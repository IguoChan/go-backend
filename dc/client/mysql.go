package client

import (
	"go-backend/config"
	"go-backend/util/mysql_client"
)

var DefaultMysqlClient *mysql_client.Client

func CloseMysqlDB() {
	if DefaultMysqlClient != nil {
		DefaultMysqlClient.Close()
	}
}

func init() {
	mysqlCfg := config.GetMysqlCfg()
	param := &mysql_client.ClientParam{
		Username:     mysqlCfg.Username,
		Password:     mysqlCfg.Password,
		Address:      mysqlCfg.Hostname,
		DatabaseName: mysqlCfg.Database,
		MaxOpenConn:  mysqlCfg.MaxConn,
		MaxIdleConn:  mysqlCfg.MaxConn,
	}
	var err error
	DefaultMysqlClient, err = mysql_client.NewMysql(param)
	if err != nil {
		panic(err)
	}
}
