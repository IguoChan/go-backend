package mq

import "fmt"

type Message struct {
	Topic    string
	Tags     string
	Keys     []string
	Body     []byte
	Property map[string]string
	Ext      map[string]interface{}
}

type MessageExt struct {
	Message
	MsgId                     string
	OffsetMsgId               string
	StoreSize                 int32
	QueueOffset               int64
	SysFlag                   int32
	BornTimestamp             int64
	BornHost                  string
	StoreTimestamp            int64
	StoreHost                 string
	CommitLogOffset           int64
	BodyCRC                   int32
	ReconsumeTimes            int32
	PreparedTransactionOffset int64
}

type MessageExtHandler func(*MessageExt) error

type Type int32

const (
	RocketMQ Type = iota + 1
	MqttV3
	MqttV5
	Kafka
)

type MQ interface {
	Publish(groupID string, msg *Message) error
	Subscribe(groupID, topic, tag string, handler MessageExtHandler) error
}

func (msg *Message) String() string {
	return fmt.Sprintf("[Topic: %s, Tags: %s, Keys: %s, Body: %s, Property: %v]",
		msg.Topic, msg.Tags, msg.Keys, string(msg.Body), msg.Property)
}

func (msgExt *MessageExt) String() string {
	return fmt.Sprintf("[Message=%s, MsgId=%s, OffsetMsgId=%s, StoreSize=%d, QueueOffset=%d, SysFlag=%d, "+
		"BornTimestamp=%d, BornHost=%s, StoreTimestamp=%d, StoreHost=%s, CommitLogOffset=%d, BodyCRC=%d, "+
		"ReconsumeTimes=%d, PreparedTransactionOffset=%d]", msgExt.Message.String(), msgExt.MsgId, msgExt.OffsetMsgId,
		msgExt.StoreSize, msgExt.QueueOffset, msgExt.SysFlag, msgExt.BornTimestamp, msgExt.BornHost,
		msgExt.StoreTimestamp, msgExt.StoreHost, msgExt.CommitLogOffset, msgExt.BodyCRC, msgExt.ReconsumeTimes,
		msgExt.PreparedTransactionOffset)
}
