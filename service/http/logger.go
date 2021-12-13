package http

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go-backend/util"
	"go-backend/util/log"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type colorModeValue int

const (
	autoColor colorModeValue = iota
	disableColor
	forceColor
)

var (
	httpLogger      *logrus.Logger
	loggerColorMode = autoColor
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

// Write 返回 response body 时并存储它
func (w *bodyLogWriter) Write(b []byte) (int, error) {
	//memory copy here!
	w.bodyBuf.Write(b)
	return w.ResponseWriter.Write(b)
}

type logFormatterParams struct {
	Request *http.Request

	// TimeStamp shows the time after the server returns a response.
	TimeStamp time.Time
	// StatusCode is HTTP response code.
	StatusCode int
	// Latency is how much time the server cost to process a certain request.
	Latency time.Duration
	// ClientIP equals Context's ClientIP method.
	ClientIP string
	// Method is the HTTP method given to the request.
	Method string
	// Path is a path the client requests.
	Path string
	// ErrorMessage is set if error has occurred in processing the request.
	ErrorMessage string
	// isTerm shows whether does gin's output descriptor refers to a terminal.
	isTerm bool
	// BodySize is the size of the Response Body
	BodySize int
	// Keys are the keys set on the request's context.
	Keys map[string]interface{}
	// ReqBody is the Request Body
	ReqBody []byte
	// ResBody is the Response Body
	ResBody *bodyLogWriter
	// logReqBody should log request body
	logReqBody bool
	// logResBody should log response body
	logResBody bool
	// traceID is OAuth-Nonce or TraceId in header
	traceID string
}

func (p *logFormatterParams) IsOutputColor() bool {
	return loggerColorMode == forceColor || (loggerColorMode == autoColor && p.isTerm)
}

// StatusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func (p *logFormatterParams) StatusCodeColor() string {
	code := p.StatusCode

	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// MethodColor is the ANSI color for appropriately logging http method to a terminal.
func (p *logFormatterParams) MethodColor() string {
	method := p.Method

	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}

// ResetColor resets all escape attributes.
func (p *logFormatterParams) ResetColor() string {
	return reset
}

// logFormatter is the default log format function Logger middleware uses.
var logFormatter = func(param logFormatterParams) string {
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}

	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}

	if param.logReqBody {
		if param.logResBody {
			return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v [TraceId: %s]\n\trequest %s\n\tresponse %s\n%s\n",
				param.TimeStamp.Format("2006/01/02 - 15:04:05.000"),
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
				param.traceID,
				string(param.ReqBody),
				param.ResBody.bodyBuf.String(),
				param.ErrorMessage,
			)
		}
		return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v [TraceId: %s]\n\trequest %s\n%s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05.000"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.traceID,
			string(param.ReqBody),
			param.ErrorMessage,
		)
	} else if param.logResBody {
		return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v [TraceId: %s]\n\tresponse %s\n%s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05.000"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.traceID,
			param.ResBody.bodyBuf.String(),
			param.ErrorMessage,
		)
	} else {
		return fmt.Sprintf("[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %#v [TraceId: %s]\n%s\n",
			param.TimeStamp.Format("2006/01/02 - 15:04:05.000"),
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.traceID,
			param.ErrorMessage,
		)
	}
}

func entryLog(start time.Time, method, path, query, traceID string) {
	if query != "" {
		path = path + "?" + query
	}
	httpLogger.Infof("[GIN] %v %-7s %#v [TraceID: %s] start\n",
		start.Format("2006/01/02 - 15:04:05:000"),
		method,
		path,
		traceID)
}

func stdLogger(isReqBody, isResBody bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		traceID := ""
		switch {
		case c.GetHeader("OAuth-Nonce") != "":
			traceID = c.GetHeader("OAuth-Nonce")
		case c.GetHeader("TraceID") != "":
			traceID = c.GetHeader("TraceId")
		default:
			traceID = util.GenerateTraceID()
		}
		c.Set("TraceID", traceID)

		entryLog(start, c.Request.Method, path, query, traceID)

		var reqBody []byte
		if isReqBody {
			reqBody, _ = ioutil.ReadAll(c.Request.Body)
			c.Request.Body.Close()
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(reqBody))
		}

		var blw *bodyLogWriter
		if isResBody {
			blw = &bodyLogWriter{bodyBuf: bytes.NewBufferString(""), ResponseWriter: c.Writer}
			c.Writer = blw
		}

		// Process request
		c.Next()

		param := logFormatterParams{
			Request:    c.Request,
			Keys:       c.Keys,
			logReqBody: isReqBody,
			logResBody: isResBody,
			isTerm:     util.CheckIfTerminal(httpLogger.Out),
		}

		// Stop timer
		param.TimeStamp = time.Now()
		param.Latency = param.TimeStamp.Sub(start)

		param.ClientIP = c.ClientIP()
		param.Method = c.Request.Method
		param.StatusCode = c.Writer.Status()
		param.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()

		param.BodySize = c.Writer.Size()
		param.ReqBody = reqBody
		param.ResBody = blw

		param.traceID = traceID

		if query != "" {
			path = path + "?" + query
		}

		param.Path = path
		httpLogger.Info(logFormatter(param))
	}
}

func init() {
	//formatter := new(logrus.TextFormatter)
	httpLogger = &logrus.Logger{
		Out:          gin.DefaultWriter,
		Formatter:    new(logrus.TextFormatter),
		Hooks:        make(logrus.LevelHooks),
		Level:        logrus.InfoLevel,
		ExitFunc:     os.Exit,
		ReportCaller: false,
	}
	log.SetLogrusRotateHook(httpLogger)
}
