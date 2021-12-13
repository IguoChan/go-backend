package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"go-backend/util"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

type RotatePeriod int64

const (
	MINUTE RotatePeriod = 60
	HOUR                = MINUTE * 60
	DAY                 = HOUR * 24

	MinuteFormat = "2006-01-02-15-04"
	HourFormat   = "2006-01-02-15"
	DayFormat    = "2006-01-02"
)

var (
	FormatMap = map[RotatePeriod]string{
		MINUTE: MinuteFormat,
		HOUR:   HourFormat,
		DAY:    DayFormat,
	}
	logRegex *regexp.Regexp
)

type RotateHook struct {
	lock                 *sync.RWMutex
	logDir               string
	fileName             string
	rotatePeriod         RotatePeriod
	currentSegmentEnd    time.Time
	currentSegmentOutput *os.File
	logLoggers           []*log.Logger
	logrusLoggers        []*logrus.Logger
}

func (h *RotateHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *RotateHook) Fire(entry *logrus.Entry) error {
	if entry.Time.Unix() >= h.currentSegmentEnd.Unix() {
		h.ChangeOutput(entry.Time)
	}
	return nil
}

func NewLogRotateHook(logDir, filename string, rotatePeriod RotatePeriod) *RotateHook {
	hook := &RotateHook{
		lock:         new(sync.RWMutex),
		logDir:       logDir,
		fileName:     filename,
		rotatePeriod: rotatePeriod,
	}
	hook.ChangeOutput(time.Now())
	return hook
}

func (h *RotateHook) ChangeOutput(t time.Time) {
	h.lock.RLock()
	defer h.lock.RUnlock()

	if t.Unix() < h.currentSegmentEnd.Unix() {
		return
	}

	cso := h.currentSegmentOutput
	if cso != nil {
		defer func() {
			go func() {
				defer cso.Close()

				// 压缩
				backPath := cso.Name() + ".gz"
				backFile, err := os.Create(backPath)
				if err != nil {
					fmt.Printf("logrotate create backfile err: %+v\n", err)
					return
				}
				defer backFile.Close()

				_, err = cso.Seek(0, 0)
				if err != nil {
					fmt.Printf("logrotate seek %s err: %+v\n", cso.Name(), err)
					return
				}
				err = util.Zip(backFile, cso)
				if err != nil {
					fmt.Printf("logrotate zip %s err: %+v\n", cso.Name(), err)
					return
				}
				err = os.Remove(cso.Name())
				if err != nil {
					fmt.Printf("logrotate remove %s err: %+v\n", cso.Name(), err)
					return
				}
			}()
		}()
	}

	sec := t.Unix()
	endUnix := sec - (sec % int64(h.rotatePeriod)) + int64(h.rotatePeriod)

	suffix := t.Format(FormatMap[h.rotatePeriod])
	cr, err := os.OpenFile(path.Join(h.logDir, strings.Join([]string{h.fileName, suffix}, ".")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Printf("logrotate change output err: %+v\n", err)
		return
	}

	h.currentSegmentEnd = time.Unix(endUnix, 0)

	for i := range h.logLoggers {
		h.logLoggers[i].SetOutput(cr)
	}
	for i := range h.logrusLoggers {
		h.logrusLoggers[i].SetOutput(cr)
	}
	h.currentSegmentOutput = cr
}

func (h *RotateHook) RegisterLogrusLogger(logger *logrus.Logger) {
	logger.SetOutput(h.currentSegmentOutput)
	h.logrusLoggers = append(h.logrusLoggers, logger)
}

func (h *RotateHook) RegisterLogLogger(logger *log.Logger) {
	logger.SetOutput(h.currentSegmentOutput)
	h.logLoggers = append(h.logLoggers, logger)
}
