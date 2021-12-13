package mqttV5_client

import (
	"context"
	"errors"
	"fmt"
	"github.com/eclipse/paho.golang/packets"
	mqttV5 "github.com/eclipse/paho.golang/paho"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go-backend/config"
	"go-backend/middleware/mq"
	"go-backend/util"
	"go-backend/util/log"
	"net"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type MqttV5Client struct {
	c *mqttV5.Client
}

var (
	defaultClient *mqttV5.Client
	connectStatus atomic.Value
	defaultMsgCh  = make(chan *mqttV5.Publish, 1000)
	stopReceive   = make(chan bool)
	subscribeMap  = make(map[string]string)
	mu            = &sync.Mutex{}
)

type connectStatusType int32

const (
	DISCONNECT connectStatusType = iota
	CONNECTING
	CONNECTED

	defaultQOS = byte(1)
)

func connectMQTTV5Broker() error {
	cfg := config.GetMQTTCfg()
	connBroker, err := net.Dial("tcp", cfg.BrokerUrl)
	if err != nil {
		logrus.Infof("create mqtt err: %+v", err)
		return err
	}
	clientCfg := mqttV5.ClientConfig{
		Conn: connBroker,
		OnServerDisconnect: func(disconnect *mqttV5.Disconnect) {
			logrus.Infof("OnServerDisconnect disconnect: %+v", disconnect)
			if disconnect.ReasonCode == packets.DisconnectServerShuttingDown && disconnect.Properties.ReasonString == "GONE" {
				logrus.Error(disconnect.Properties.ReasonString)
				connectStatus.Store(DISCONNECT)
				reconnect()
			}
		},
		OnClientError: func(err error) {
			logrus.Infof("OnClientError err: %+v", err)
			if err != nil {
				connectStatus.Store(DISCONNECT)
				reconnect()
			}
		},
	}
	defaultClient = mqttV5.NewClient(clientCfg)
	conn := &mqttV5.Connect{
		KeepAlive:    10,
		ClientID:     cfg.ClientIDV5 + uuid.NewV4().String(),
		CleanStart:   true,
		Username:     cfg.Username,
		Password:     []byte(cfg.Password),
		UsernameFlag: true,
		PasswordFlag: true,
	}
	ca, err := defaultClient.Connect(context.Background(), conn)
	if err != nil {
		logrus.Infof("connect mqtt broker failed, err:%+v", err)
		return err
	}

	logger := log.GetLogLogger("[MQTT V5] ")
	defaultClient.SetErrorLogger(logger)
	if cfg.Log {
		defaultClient.SetDebugLogger(logger)
		defaultClient.PingHandler.SetDebug(logger)
		defaultClient.Router.SetDebugLogger(logger)
	}

	logrus.Infof("connect mqtt broker %s, clientID: %s, reasonCode: %+v", cfg.BrokerUrl, cfg.ClientIDV5, ca)
	connectStatus.Store(CONNECTED)
	return nil
}

func reconnect() {
	if connectStatus.Load().(connectStatusType) == DISCONNECT {
		connectStatus.Store(CONNECTING)
		for i := 0; i < config.GetMQTTCfg().ReconnectBrokerCount; i++ {
			if err := connectMQTTV5Broker(); err != nil {
				logrus.Infof("reconnect to mqtt v5 broker failed: %d times", i+1)
				time.Sleep(10 * time.Second)
			} else {
				connectStatus.Store(CONNECTED)
				return
			}
		}
		connectStatus.Store(DISCONNECT)
	}
}

func Disconnect() error {
	d := &mqttV5.Disconnect{ReasonCode: 0}
	err := defaultClient.Disconnect(d)
	if err != nil {
		logrus.Errorf("disconnect mqtt v5 failed, err %+v", err)
		return err
	}
	return nil
}

func (m *MqttV5Client) Publish(groupID string, msg *mq.Message) error {
	topic := msg.Topic
	if msg.Tags != "" {
		topic = fmt.Sprintf("%s/%s", msg.Topic, msg.Tags)
	}

	qos, ok := msg.Ext["QOS"].(byte)
	if !ok {
		qos = defaultQOS
	}

	retained, ok := msg.Ext["RETAIN"].(bool)
	if !ok {
		retained = false
	}

	var uerProperties mqttV5.UserProperties = nil
	if msg.Property != nil {
		uerProperties = make(mqttV5.UserProperties, 0, len(msg.Property))
		for k, v := range msg.Property {
			uerProperties.Add(k, v)
		}
	}
	me := config.GetMQTTCfg().MessageExpirySeconds
	properties := mqttV5.PublishProperties{
		MessageExpiry: &me,
		User:          uerProperties,
	}

	sendMsg := &mqttV5.Publish{
		Topic:      topic,
		QoS:        qos,
		Payload:    msg.Body,
		Retain:     retained,
		Properties: &properties,
	}
	res, err := defaultClient.Publish(context.Background(), sendMsg)
	if err != nil {
		logrus.Errorf("mqtt v5 publish failed, err %+v , code %+v", err, util.JsonString(res))
		return err
	}
	logrus.Infof("mqtt v5 publish %+v success code %+v", util.JsonString(msg), util.JsonString(res))
	return nil
}
func convertToMessageExt(msg *mqttV5.Publish) *mq.MessageExt {
	message := &mq.MessageExt{
		Message: mq.Message{
			Body: msg.Payload,
		},
	}
	properties := make(map[string]string, 0)
	for _, v := range msg.Properties.User {
		properties[v.Key] = v.Value
	}
	message.Message.Property = properties
	topics := strings.Split(msg.Topic, "/")
	if len(topics) > 1 {
		message.Topic = topics[0]
		message.Tags = topics[1]
	}
	return message
}

func (m *MqttV5Client) Subscribe(groupID, topic, tag string, handler mq.MessageExtHandler) error {
	if tag != "" {
		topic = fmt.Sprintf("%s/%s", topic, tag)
	}

	mu.Lock()
	defer mu.Unlock()
	if _, ok := subscribeMap[topic]; ok {
		return nil
	}

	defaultClient.Router.RegisterHandler(topic, func(publish *mqttV5.Publish) {
		defaultMsgCh <- publish
	})

	s := &mqttV5.Subscribe{
		Subscriptions: map[string]mqttV5.SubscribeOptions{
			topic: {QoS: defaultQOS},
		},
	}

	res, err := defaultClient.Subscribe(context.Background(), s)
	if err != nil {
		logrus.Errorf("mqtt v5 subscribe failed, err %+v,  reason %+v", err, string(res.Reasons))
		return err
	}
	logrus.Infof("mqtt v5 subscribe res, properties %+v reason %+v", util.JsonString(res.Properties), res.Reasons)
	go func() {
		for {
			select {
			case msg := <-defaultMsgCh:
				handler(convertToMessageExt(msg))
			case <-stopReceive:
				return
			}
		}
	}()
	subscribeMap[topic] = topic
	return nil
}

func NewMqttV5Client() (*MqttV5Client, error) {
	if connectStatus.Load().(connectStatusType) != CONNECTED {
		return nil, errors.New("mqtt v5 has not been init")
	}
	return &MqttV5Client{
		defaultClient,
	}, nil
}

func init() {
	connectStatus.Store(DISCONNECT)
	if err := connectMQTTV5Broker(); err != nil {
		panic(err)
	}
}
