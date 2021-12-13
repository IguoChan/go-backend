package mqtt_v3_client

import (
	"fmt"
	mqttV3 "github.com/eclipse/paho.mqtt.golang"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go-backend/config"
	"time"
)

var (
	defaultClient mqttV3.Client
)

func onConnectHandler(client mqttV3.Client) {
	// TODO: called when the client is connected. Both at initial connection time and upon automatic reconnect.
	logrus.Debugf("mqtt %s on connect")
}

func connectionLostHandler(client mqttV3.Client, err error) {
	// TODO: will be executed in the case where the client unexpectedly loses connection with the mqttV3 broker.
	//logrus.Debugf("mqtt %s connect lost %s", clientID, err)
	//go opentracer.ReportError("MQTTConnectionLost", err.Error())
}

func connectMQTTV3Broker() error {
	cfg := config.GetMQTTCfg()
	clientID := cfg.ClientIDV3 + uuid.NewV4().String()

	opts := mqttV3.NewClientOptions()
	opts.AddBroker(cfg.BrokerUrl).SetClientID(clientID).SetUsername(cfg.Username).SetPassword(cfg.Password)
	opts.SetMaxReconnectInterval(10 * time.Second).SetCleanSession(false).SetResumeSubs(true).SetKeepAlive(10 * time.Second)
	opts.SetConnectionLostHandler(connectionLostHandler).SetOnConnectHandler(onConnectHandler)
	defaultClient = mqttV3.NewClient(opts)
	if token := defaultClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("connect to mqtt broker %s using clientid %s error: %s", cfg.BrokerUrl, clientID, token.Error())
	}
	logrus.Infof("connect mqtt broker %s, clientid %s", cfg.BrokerUrl, clientID)
	return nil
}

func init() {
	logger := log.
}
