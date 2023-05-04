package server

import (
	"errors"
	"io"
	"log"
	"my_1m_chatroom/internal/controller"
	"reflect"
	"sync"
	"syscall"
	"time"

	"github.com/spf13/cast"
	"golang.org/x/sys/unix"

	// "nhooyr.io/websocket"
	// "nhooyr.io/websocket/wsjson"
	"github.com/gorilla/websocket"
)

type epoll struct {
	fd          int
	connections map[int]websocket.Conn
	users       map[int]*controller.User
	lock        *sync.RWMutex
}

var (
	Epoller  *epoll
	batchNum = 100
)

func init() {
	fd, err := unix.EpollCreate1(0)
	if err != nil {
		log.Printf("failed to init epoller, error = %v\n", err)
		return
	}

	Epoller = &epoll{
		fd:          fd,
		connections: make(map[int]websocket.Conn),
		lock:        &sync.RWMutex{},
		users:       make(map[int]*controller.User),
	}
}

func websocketFd1(conn websocket.Conn) int {
	v := reflect.Indirect(reflect.ValueOf(conn))
	rwc := reflect.Indirect(v.FieldByName("rwc")).Elem()
	tcpConn := reflect.Indirect(rwc).FieldByName("conn")
	fdVal := reflect.Indirect(tcpConn).FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}

// gobwas/ws
func websocketFd(conn websocket.Conn) int {
	tcpConn := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")

	return int(pfdVal.FieldByName("Sysfd").Int())
}

// func GetFdFromConn(c websocket.Conn) int {
// 	v := reflect.Indirect(reflect.ValueOf(c))
// 	conn := v.FieldByName("conn")
// 	netFD := reflect.Indirect(conn.FieldByName("fd"))
// 	pfd := netFD.FieldByName("pfd")
// 	fd := int(pfd.FieldByName("Sysfd").Int())
// 	return fd
// }

// origin websocket
func GetFdFromConn(conn websocket.Conn) int {
	connVal := reflect.Indirect(reflect.ValueOf(conn)).FieldByName("conn").Elem()
	tcpConn := reflect.Indirect(connVal).FieldByName("conn")
	fdVal := tcpConn.FieldByName("fd")
	pfdVal := reflect.Indirect(fdVal).FieldByName("pfd")
	return int(pfdVal.FieldByName("Sysfd").Int())
}

func (e *epoll) Add(conn websocket.Conn, user *controller.User) error {
	fd := GetFdFromConn(conn)
	err := syscall.EpollCtl(e.fd, syscall.EPOLL_CTL_ADD, fd, &syscall.EpollEvent{Events: unix.POLLIN | unix.POLLHUP, Fd: int32(fd)})
	if err != nil {
		return err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	e.users[fd] = user
	e.connections[fd] = conn

	return nil
}

func (e *epoll) Remove(conn websocket.Conn) error {
	fd := GetFdFromConn(conn)
	err := syscall.EpollCtl(e.fd, syscall.EPOLL_CTL_DEL, fd, nil)
	if err != nil {
		return err
	}

	e.lock.Lock()
	defer e.lock.Unlock()
	delete(e.users, fd)
	delete(e.connections, fd)
	return nil

}

func (e *epoll) wait() ([]websocket.Conn, []int, error) {
	events := make([]unix.EpollEvent, batchNum)
	n, err := unix.EpollWait(e.fd, events, 100)

	if err != nil {
		return nil, nil, err
	}
	e.lock.Lock()
	defer e.lock.Unlock()
	connections := make([]websocket.Conn, n)
	fds := make([]int, n)
	for i := 0; i < n; i++ {
		connections[i] = e.connections[int(events[i].Fd)]
		fds[i] = int(events[i].Fd)
	}
	return connections, fds, nil

}

func Start() {
	for {
		connections, fds, err := Epoller.wait()
		if err != nil {
			log.Printf("Failed to epoll wait %v", err)
			continue
		}

		for i, connection := range connections {
			receivedMsg := make(map[string]string)
			var err error

			// err = wsjson.Read(context.Background(), &connection, &receivedMsg)
			err = connection.ReadJSON(&receivedMsg)
			// data, _, err := wsutil.ReadClientData(connection)
			// json.Unmarshal(data, &receivedMsg)
			if err != nil {
				// var closeErr websocket.CloseError
				if errors.As(err, &websocket.CloseError{}) {
					log.Printf("Connection error: %v\n", err)
					continue
				} else if errors.Is(err, io.EOF) {
					continue
				}
			}

			deliverMsg := controller.NewMessage(controller.MsgTypeNormal, Epoller.users[fds[i]], receivedMsg["content"], time.Unix(0, cast.ToInt64(receivedMsg["send_time"])))
			controller.Broadcaster.Broadcast(deliverMsg)
		}
	}
}
