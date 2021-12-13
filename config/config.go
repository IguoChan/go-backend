package config

import (
	"bytes"
	"encoding/json"
	"flag"
	"github.com/sirupsen/logrus"
	"go-backend/util"
	"gopkg.in/ini.v1"
	"os"
	"regexp"
	"testing"
	"time"
)

var envRegex = "^\\$\\{.*\\}$"

var cfg = &Configuration{
	Help:       false,
	ConfigFile: "./main/config.ini",

	Server: ServerConfig{
		LogMode:    "release",
		Port:       "9410",
		DeployMode: "multi-machine",
	},

	Redis: RedisConfig{
		Hostname:      "192.168.1.211",
		Port:          "6379",
		Password:      "123456",
		MaxConnection: 64,
		Timeout:       5 * time.Second,
		Log:           true,
	},

	Mysql: MysqlConfig{
		Username: "root",
		Password: "123456",
		Hostname: "127.0.0.1:3306",
		Database: "go_backend",
		MaxConn:  64,
		Log:      true,
	},

	MQTT: MqttConfig{
		BrokerUrl:            "127.0.0.1:1883",
		ClientIDV3:           "client_id_v3",
		ClientIDV5:           "client_id_v5",
		Username:             "username",
		Password:             "password",
		Log:                  true,
		ReconnectBrokerCount: 6,
		MessageExpirySeconds: 30,
	},
}

func iniLoadFile(path string) error {
	cfgFile, err := ini.Load(path)
	if err != nil {
		return err
	}
	parseConfigFile(cfgFile)

	if err := cfgFile.MapTo(cfg); err != nil {
		return err
	}

	return nil
}

func parseConfigFile(cfg *ini.File) {
	sections := cfg.Sections()
	for _, section := range sections {
		keys := section.Keys()
		emptyKeys := make([]string, 0, len(keys))
		for _, key := range keys {
			value := key.Value()
			match, err := regexp.MatchString(envRegex, value)
			if err != nil {
				logrus.Errorf("regexp match string %s %s error: %s", value, envRegex, err)
				continue
			}
			if match {
				key.SetValue(os.Getenv(value[2 : len(value)-1]))
			}

			if key.Value() == "" {
				emptyKeys = append(emptyKeys, key.Name())
			}
		}

		for _, key := range emptyKeys {
			section.DeleteKey(key)
		}
	}
}

func PrintConfigInfo() {
	tempCfg := *cfg
	str := bytes.NewBuffer([]byte("go back-end config: "))
	encoder := json.NewEncoder(str)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	if err := encoder.Encode(tempCfg); err != nil {
		logrus.Error("encode print info failed")
		return
	}
	logrus.Info(str.String())
}

func GetServerCfg() *ServerConfig {
	return &cfg.Server
}

func GetRedisCfg() *RedisConfig {
	return &cfg.Redis
}

func GetMysqlCfg() *MysqlConfig {
	return &cfg.Mysql
}

func GetMQTTCfg() *MqttConfig {
	return &cfg.MQTT
}

func init() {
	testing.Init()

	flagSet := flag.CommandLine
	if err := registerFlagSet(cfg, flagSet); err != nil {
		panic(err)
	}

	args := os.Args[1:]
	if cfg.ConfigFile != "" && util.IsExists(cfg.ConfigFile) {
		if err := iniLoadFile(cfg.ConfigFile); err != nil {
			panic(err)
		}
	}
	if err := flagSet.Parse(args); err != nil {
		panic(err)
	}

	if cfg.Help {
		flagSet.Usage()
		os.Exit(0)
		return
	}
	PrintConfigInfo()
}
