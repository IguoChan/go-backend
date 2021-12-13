package http

import (
	"context"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var httpServerHandle *http.Server

func setRoutes(e *gin.Engine) {
	e.GET("/health", stdLogger(false, false), func(c *gin.Context) {
		c.String(http.StatusOK, "OK")
	})
}

func StartHTTP(listenAddr string) error {
	e := gin.New()
	pprof.Register(e)
	setRoutes(e)
	httpServerHandle = &http.Server{
		Addr:    listenAddr,
		Handler: e,
	}
	return httpServerHandle.ListenAndServe()
}

func StopHTTP() {
	if httpServerHandle != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := httpServerHandle.Shutdown(ctx); err != nil {
			logrus.Errorf("shutdown http server error: %s\n", err.Error())
		}
	}
}
