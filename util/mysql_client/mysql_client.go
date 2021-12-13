package mysql_client

import (
	"fmt"
	"go-backend/config"
	"go-backend/util/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gLogger "gorm.io/gorm/logger"
	"time"
)

type Client struct {
	*gorm.DB
}

type ClientParam struct {
	Username     string
	Password     string
	Address      string
	DatabaseName string
	MaxOpenConn  int
	MaxIdleConn  int
}

func NewMysql(param *ClientParam) (*Client, error) {
	logLevel := gLogger.Silent
	if config.GetMysqlCfg().Log {
		logLevel = gLogger.Info
	}
	newLogger := gLogger.New(
		log.GetLogLogger("[MYSQL] "),
		gLogger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logLevel,
			Colorful:      false,
		},
	)
	cfg := mysql.Config{
		DSN: fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			param.Username, param.Password, param.Address, param.DatabaseName), // DSN data source name
		DefaultStringSize:         256,
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}
	db, err := gorm.Open(mysql.New(cfg), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(param.MaxIdleConn)
	sqlDB.SetMaxOpenConns(param.MaxOpenConn)
	sqlDB.SetConnMaxLifetime(time.Hour * 4)
	return &Client{DB: db}, nil
}

func (c *Client) Close() {
	d, _ := c.DB.DB()
	_ = d.Close()
}
