package middleware

import (
	"errors"
	"fmt"
	"go-backend/middleware/mq"
	mqttV5 "go-backend/util/mqtt_v5_client"
)

func GetMQ(t mq.Type) (mq.MQ, error) {
	switch t {
	//case RocketMQ:
	case mq.MqttV5:
		return mqttV5.NewMqttV5Client()
	default:
		return nil, errors.New(fmt.Sprintf("mq type[%+v] not supported", t))
	}
}
