package api

import (
	"my_1m_chatroom/internal/service"

	"github.com/gin-gonic/gin"
)

func ChatHandlerRoutes(group *gin.RouterGroup) {
	group.GET("/ws", service.ConnectionHandler)
	group.GET("/user_list", service.UserListHandleFunc)
	group.GET("/", service.WebPageHandler)

}
