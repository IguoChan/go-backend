package log

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

var (
	Hook   *RotateHook
	logger = logrus.StandardLogger()
)

const (
	LOGGER_DIR       = "/home/iguochan/log/go-backend"
	LOGGER_FILE_NALE = "go-backend"
)

func GetLogrusLogger() *logrus.Logger {
	return logger
}

func GetLogLogger(prefix string) *log.Logger {
	logger := log.New(logger.Out, prefix, log.LstdFlags)
	SetLogRotateHook(logger)
	return logger
}

func SetLogrusRotateHook(logger *logrus.Logger) {
	logger.AddHook(Hook)
	Hook.RegisterLogrusLogger(logger)
}

func SetLogRotateHook(logger *log.Logger) {
	Hook.RegisterLogLogger(logger)
}

func init() {
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "[2006-01-02T15:04:05.000]"
	logrus.SetFormatter(customFormatter)

	os.MkdirAll(LOGGER_DIR, 0755)
	Hook = NewLogRotateHook(LOGGER_DIR, LOGGER_FILE_NALE, HOUR)
	SetLogrusRotateHook(logger)
}
