package config

import "time"

type Configuration struct {
	Help       bool   `config:"h" usage:"help, only show config params"`
	ConfigFile string `config:"f" usage:"config file"`

	Server ServerConfig `config:"server" usage:"" ini:"server"`
	Redis  RedisConfig  `config:"redis" usage:"" ini:"redis"`
	Mysql  MysqlConfig  `config:"mysql" usage:"" ini:"mysql"`
	MQTT   MqttConfig   `config:"mqtt" usage:"" ini:"mqtt"`
}

type ServerConfig struct {
	LogMode    string `config:"log_mode" usage:"log mode" ini:"log_mode"`
	Port       string `config:"port" usage:"server listen port" ini:"port"`
	DeployMode string `config:"deploy_mode" usage:"deploy mode" ini:"deploy_mode"`
}

type RedisConfig struct {
	Hostname      string        `config:"hostname" usage:"redis hostname" ini:"hostname"`
	Port          string        `config:"port" usage:"redis port" ini:"port"`
	Password      string        `config:"password" usage:"redis password" ini:"password"`
	MaxConnection int           `config:"max_connection" usage:"redis max_connection" ini:"max_connection"`
	Timeout       time.Duration `config:"timeout" usage:"redis timeout" ini:"timeout"`
	Log           bool          `config:"log" usage:"redis log" ini:"log"`
}
type MysqlConfig struct {
	Username string `config:"username" usage:"mysql username" ini:"username"`
	Password string `config:"password" usage:"mysql password" ini:"password"`
	Hostname string `config:"hostname" usage:"mysql hostname" ini:"hostname"`
	Database string `config:"database" usage:"mysql database" ini:"database"`
	MaxConn  int    `config:"max_connect" usage:"mysql max_connect" ini:"max_connect"`
	Log      bool   `config:"log" usage:"mysql log" ini:"log"`
}

type MqttConfig struct {
	BrokerUrl            string `config:"broker_url" usage:"mqtt broker_url" ini:"broker_url"`
	ClientIDV3           string `config:"client_id_v3" usage:"mqtt client_id_v3" ini:"client_id_v3"`
	ClientIDV5           string `config:"client_id_v5" usage:"mqtt client_id_v5" ini:"client_id_v5"`
	Username             string `config:"username" usage:"mqtt username" ini:"username"`
	Password             string `config:"password" usage:"mqtt password" ini:"password" json:",omitempty"`
	Log                  bool   `config:"log" usage:"mqtt log" ini:"log"`
	MessageExpirySeconds uint32 `config:"message_expiry_seconds" usage:"set mqtt message expiry time" ini:"message_expiry_seconds"`
	ReconnectBrokerCount int    `config:"reconnect_broker_count" usage:"set mqtt reconnect broker count" ini:"reconnect_broker_count"`
}
