package controller

import (
	"regexp"
	"time"
)

const (
	MsgTypeNormal = iota
	MsgTypeWelcome
	MsgTypeUserEnter
	MsgTypeUserLeave
	MsgTypeError
)

// 给用户发送的消息
type Message struct {
	// 哪个用户发送的消息
	User    *User     `json:"user"`
	Type    int       `json:"type"`
	Content string    `json:"content"`
	MsgTime time.Time `json:"msg_time"`

	ClientSendTime time.Time `json:"client_send_time"`

	//  @
	Ats []string `json:"ats"`

	// 用户列表不通过 WebSocket 下发
	// Users []*User `json:"users"`
}

func NewMessage(msgType int, user *User, content string, sendTime time.Time) *Message {
	newMsg := &Message{
		Type:           msgType,
		User:           user,
		Content:        content,
		ClientSendTime: sendTime,
		MsgTime:        time.Now(),
	}

	reg := regexp.MustCompile(`@[^\s@]{2,20}`)
	newMsg.Ats = reg.FindAllString(newMsg.Content, -1)

	return newMsg
}
