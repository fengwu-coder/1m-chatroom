package main

import (
	"crypto/tls"
	"fmt"
	"my_1m_chatroom/api"
	"my_1m_chatroom/config"
	"my_1m_chatroom/internal/controller"
	"my_1m_chatroom/internal/server"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

var (
	configFile = "/home/feng/go/src/my_1m_chatroom/config/config.yaml"
)

func init() {
	config.InitConfig(configFile)
}

func main() {

	r := api.NewRouter(config.CRConfig)

	httpServer := &http.Server{
		Addr:      fmt.Sprintf("%s:%d", config.CRConfig.Http.Address, config.CRConfig.Http.Port),
		Handler:   r,
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println(err)
		}
	}()

	go controller.Broadcaster.Start()
	go server.Start()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-quit:
		os.Exit(1)
	}
}
