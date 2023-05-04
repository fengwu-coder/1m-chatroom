package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"my_1m_chatroom/internal/controller"
	"my_1m_chatroom/internal/server"
	"net/http"

	"github.com/gin-gonic/gin"

	websocket1 "github.com/gorilla/websocket"
	// "nhooyr.io/websocket"
)

func ConnectionHandler(c *gin.Context) {
	req := c.Request
	w := c.Writer
	// conn2 :=
	// ws.DefaultDialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	// ws.NewMask()
	// conn, _, _, err := ws.UpgradeHTTP(req, resp)
	// header := ws.Header{}
	// header.Masked = false
	// var mask [4]byte = [4]byte{0, 0, 0, 0}
	// header.Mask = mask
	// ws.WriteHeader(conn, header)
	// ws.WriteHeader(conn, ws.Header{Masked: false, Mask: mask})
	// Upgrade connection
	upgrader := websocket1.Upgrader{}
	conn, err := upgrader.Upgrade(w, req, nil)
	// fmt.Println(conn1)
	// , &websocket.AcceptOptions{InsecureSkipVerify: true}
	// conn, err := websocket.Accept(w, c.Request, nil)

	if err != nil {
		return
	}
	// fmt.Println(conn1)
	token := req.FormValue("token")
	name := req.FormValue("nickname")
	user := controller.NewUser(conn, name, req.RemoteAddr, token)

	if err := server.Epoller.Add(*conn, user); err != nil {
		log.Printf("Failed to add connection %v", err)
		// conn.Close(websocket.StatusAbnormalClosure, err.Error())
	}
	// wsjson.Write(req.Context(), conn, )
	// wsutil.WriteServerMessage(c.Writer, ws.OpText, []byte("123"))
	go user.DistributeMessage(context.Background())

	user.Mc <- *controller.WelcomeMessage(user)

	enterMsg := controller.EnterMessage(user)

	controller.Broadcaster.AddEnteringUser(user)
	log.Println("user:", name, "joins chat")
	controller.Broadcaster.Broadcast(enterMsg)

	// err = user.ReceiveMessage(context.Background())
	// 根据读取时的错误执行不同的 Close
	// if err == nil {
	// 	log.Println("close connection")
	// 	conn.Close(ws.StatusNormalClosure, "")
	// } else {
	// 	log.Println("read from client error:", err)
	// 	conn.Close(ws.StatusInternalError, "Read from client error")
	// }

}

func UserListHandleFunc(c *gin.Context) {
	c.Writer.Header().Add("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)

	userList := controller.Broadcaster.GetUserList()
	b, err := json.Marshal(userList)

	if err != nil {
		fmt.Fprint(c.Writer, `[]`)
	} else {
		fmt.Fprint(c.Writer, string(b))
	}
}

// func ConnectionHandler(c *gin.Context) {
// 	conn, err := websocket.Accept(c.Writer, c.Request, &websocket.AcceptOptions{InsecureSkipVerify: true})
// 	if err != nil {
// 		log.Println("websocket accept error:", err)
// 		return
// 	}
// 	req := c.Request
// 	token := req.FormValue("token")
// 	name := req.FormValue("nickname")
// 	user := controller.NewUser(conn, name, req.RemoteAddr, token)

// 	go user.DistributeMessage(req.Context())

// 	user.Mc <- *controller.WelcomeMessage(user)

// 	enterMsg := controller.EnterMessage(user)

// 	controller.Broadcaster.AddEnteringUser(user)
// 	log.Println("user:", name, "joins chat")
// 	controller.Broadcaster.Broadcast(enterMsg)

// 	err = user.ReceiveMessage(req.Context())
// 	// 根据读取时的错误执行不同的 Close
// 	if err == nil {
// 		log.Println("close connection")
// 		conn.Close(websocket.StatusNormalClosure, "")
// 	} else {
// 		log.Println("read from client error:", err)
// 		conn.Close(websocket.StatusInternalError, "Read from client error")
// 	}

// }
