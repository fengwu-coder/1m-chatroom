package api

import (
	"log"
	"my_1m_chatroom/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

func NewRouter(config config.ChatRoomConfig) *gin.Engine {
	r := gin.New()

	r.Use(cors.Default())

	r.Use(gin.Recovery())
	// r.Use(gindump.DumpWithOptions(true, true, false, false, false, func(dumpStr string) {
	// 	loggers.LogInstance().Debugln(dumpStr)
	// }))
	// r.Use(i18n.LangMiddleware())
	r.Use(gindump.DumpWithOptions(true, true, false, false, false, func(dumpStr string) {
		log.Println(dumpStr)
	}))
	chatGroup := r.Group("/api/v1")

	ChatHandlerRoutes(chatGroup)

	return r
}
