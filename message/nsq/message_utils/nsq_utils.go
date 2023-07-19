package message_utils

import (
	"auction-website/utils"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"time"
)

type NSQMessage struct {
	Topic     string `json:"topic"`      // 消息主题
	Channel   string `json:"channel"`    // 消息队列
	Timestamp int64  `json:"timestamp"`  // 消息时间戳
	MessageID string `json:"message_id"` // 消息ID
	Body      []byte `json:"body"`       // 消息体
}

// 创建nsq信息
func CreateMessage(topic, channel string, msg interface{}) ([]byte, error) {
	body, err := utils.ToJSONBytes(msg)
	if err != nil {
		return nil, err
	}
	// 创建一个NSQMessage消息
	message := NSQMessage{
		Topic:     topic,
		Channel:   channel,
		Timestamp: time.Now().Unix(),
		MessageID: uuid.New().String(),
		Body:      body,
	}
	// 将NSQMessage编码为JSON格式，并发布到NSQ中
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		return nil, errors.New("json marshal failed" + err.Error())
	}
	return jsonMessage, nil
}

// 解析nsq信息
func ParseMessage(jsonMessage []byte) (NSQMessage, error) {
	var message NSQMessage
	err := json.Unmarshal(jsonMessage, &message)
	if err != nil {
		return NSQMessage{}, errors.New("json unmarshal failed: " + err.Error())
	}
	//fmt.Println(message)
	return message, nil
}
