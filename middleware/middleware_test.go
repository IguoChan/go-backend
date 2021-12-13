package middleware

import (
	"fmt"
	"go-backend/middleware/mq"
	"testing"
)

func TestPublish(t *testing.T) {
	c, e := GetMQ(mq.MqttV5)
	if e != nil {
		t.Error(e)
	}
	msg := &mq.Message{
		Topic:    "test_topic",
		Tags:     "hello",
		Body:     []byte("nihao,#########################"),
		Property: map[string]string{"TraceId": "trace id"},
	}
	c.Publish("", msg)
}

func TestSubscribe(t *testing.T) {
	c, e := GetMQ(mq.MqttV5)
	if e != nil {
		t.Error(e)
	}
	h := func(msg *mq.MessageExt) error {
		fmt.Println(msg.String())
		return nil
	}
	c.Subscribe("", "test_topic", "hello", h)
	select {}
}
