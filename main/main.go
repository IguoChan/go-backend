package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go-backend/config"
	"go-backend/dc"
	"go-backend/dc/client"
	"go-backend/dc/dao"
	"go-backend/service/http"
	"go-backend/util/lock"
	_ "go-backend/util/mqtt_v5_client"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// 如果是debug模式，需设置logrus为debuglevel
	// 如果是release模式，需设置Gin为release模式（Gin默认debug模式）
	if config.GetServerCfg().LogMode == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	logrus.Infof("lalala")

	go func() {
		for true {
			_ = dc.STUDao.Add(&dao.Student{
				StuId:     "10001",
				FirstName: "Iguo",
				LastName:  "Chan",
				Email:     "iguochan@foxmail.com",
				PN:        "13161192916",
			})
			logrus.Infof("mysql test")
			time.Sleep(10 * time.Second)
		}
	}()

	//return

	// 开启http服务
	ch := make(chan struct{})
	go func() {
		listenAddr := ":" + config.GetServerCfg().Port
		err := http.StartHTTP(listenAddr)
		fmt.Println(err.Error())
		ch <- struct{}{}
	}()

	// 设置退出
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	lock := lock.GetBLocker("redis_lock")
	clientID := uuid.NewV4().String()
	c, ok := lock.BLock(clientID)
	if ok {
		fmt.Println("get lock success!")
		time.Sleep(10 * time.Second)
		fmt.Println("unlock!")
		lock.BUnLock(clientID, c)
	} else {
		fmt.Println("get lock failed!")
	}

	select {
	case <-interrupt:
		fmt.Printf("ctrl c, finish\n")
	case <-ch:
		fmt.Printf("happen exception\n")
	}

	http.StopHTTP()
	client.CloseRedisDB()
}
