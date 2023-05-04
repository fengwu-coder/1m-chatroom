package controller

import (
	"fmt"
	"time"
)

var (
	MaxUserNum = 100
	MaxMsgNum  = 128
)

type broadcaster struct {
	UserList map[int]*User

	MessageChannel chan *Message

	LeavingChannel  chan *User
	EnteringChannel chan *User

	// 获取用户列表
	requestUsersChannel chan struct{}
	usersChannel        chan []*User
}

var Broadcaster = &broadcaster{
	UserList:       make(map[int]*User),
	MessageChannel: make(chan *Message, MaxMsgNum),

	LeavingChannel:  make(chan *User),
	EnteringChannel: make(chan *User),

	requestUsersChannel: make(chan struct{}),
	usersChannel:        make(chan []*User),
}

func (broadcaster *broadcaster) GetUserList() []*User {
	broadcaster.requestUsersChannel <- struct{}{}
	return <-broadcaster.usersChannel
}

func (broadcaster *broadcaster) Broadcast(msg *Message) {
	broadcaster.MessageChannel <- msg
}

func (broadcaster *broadcaster) AddEnteringUser(user *User) {
	broadcaster.EnteringChannel <- user
}

func (broadcaster *broadcaster) Start() {
	for {
		select {
		case msg := <-broadcaster.MessageChannel:
			for _, user := range broadcaster.UserList {
				if user.UserID == msg.User.UserID {
					continue
				}
				user.Mc <- *msg
			}
		case leftUser := <-broadcaster.LeavingChannel:
			delete(broadcaster.UserList, int(leftUser.UserID))
			leftUser.DeleteUser()
		case EnteredUser := <-broadcaster.EnteringChannel:
			broadcaster.UserList[int(EnteredUser.UserID)] = EnteredUser
			// for _, user := range broadcaster.UserList {
			// 	user.Mc <- *msg
			// }
		case <-broadcaster.requestUsersChannel:
			userList := make([]*User, 0, len(broadcaster.UserList))
			for _, user := range broadcaster.UserList {
				userList = append(userList, user)
			}

			broadcaster.usersChannel <- userList
		}
	}
}

func WelcomeMessage(user *User) *Message {
	return &Message{
		User:    SystemAdmin,
		Type:    MsgTypeWelcome,
		Content: fmt.Sprintf("Hi, %s. Welcome to the Chatroom, now Let's chat!", user.UserName),
		MsgTime: time.Now(),
	}
}

func EnterMessage(user *User) *Message {
	return &Message{
		User:    SystemAdmin,
		Type:    MsgTypeUserEnter,
		Content: fmt.Sprintf("%s joined the Chatroom", user.UserName),
		MsgTime: time.Now(),
	}
}
