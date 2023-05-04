package controller

import (
	"context"
	"errors"
	"io"
	"log"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/cast"
	// "nhooyr.io/websocket"
	// "nhooyr.io/websocket/wsjson"
)

var (
	GlobalUserID uint64
	SystemAdmin  = &User{}
)

type User struct {
	// UserInfo UserInfo `json:"userinfo"`
	UserName  string       `json:"nickname"`
	UserID    int          `json:"uid"`
	EnteredAt time.Time    `json:"enter_at"`
	Address   string       `json:"addr"`
	Token     string       `json:"token"`
	Mc        chan Message `json:"-"`
	Conn      *websocket.Conn
}

// type UserInfo struct {
// 	UserName  string    `json:"username"`
// 	UserID    uint64    `json:"userid"`
// 	EnteredAt time.Time `json:"EnteredAt"`
// 	Address   string    `json:"Address"`
// 	Token     string    `json:"token"`
// }

func NewUser(conn *websocket.Conn, Name, Address, Token string) *User {
	// UserInfo := &UserInfo{
	// 	UserName:  Name,
	// 	Address:   Address,
	// 	EnteredAt: time.Now(),
	// }
	user := &User{
		// UserInfo: *UserInfo,
		UserName:  Name,
		Address:   Address,
		EnteredAt: time.Now(),
		Mc:        make(chan Message, 64),
		Conn:      conn,
	}

	if user.Token != "" {
		// TODO JWT TOKEN
		userid, err := parseTokenAndValidate(Token, Name)
		if err == nil {
			user.UserID = userid
		}
	}

	if user.UserID == 0 {
		user.UserID = int(atomic.AddUint64(&GlobalUserID, 1))
		user.Token = genToken(user.UserID, user.UserName)
	}

	return user
}

func (user *User) DistributeMessage(ctx context.Context) {
	for msg := range user.Mc {
		// data, err := json.Marshal(msg)
		// if err != nil {
		// 	log.Panicln(err)
		// }
		// wsjson.Write(ctx, user.Conn, msg)
		user.Conn.WriteJSON(msg)
		// wsutil.WriteClientMessage(user.Conn, ws.OpBinary, data)
	}
}

func (user *User) DeleteUser() {
	close(user.Mc)
}

func (user *User) ReceiveMessage(ctx context.Context) error {
	receivedMsg := make(map[string]string)
	var err error
	for {
		err = user.Conn.ReadJSON(&receivedMsg)
		// err = wsjson.Read(ctx, user.Conn, &receivedMsg)
		// data, _, err := wsutil.ReadClientData(*user.Conn)
		// json.Unmarshal(data, &receivedMsg)

		if err != nil {
			if errors.As(err, &websocket.CloseError{}) {
				log.Printf("Connection error: %v\n", err)
				return nil
			} else if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}

		deliverMsg := NewMessage(MsgTypeNormal, user, receivedMsg["content"], time.Unix(0, cast.ToInt64(receivedMsg["send_time"])))
		Broadcaster.Broadcast(deliverMsg)
	}
}
