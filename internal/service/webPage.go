package service

import (
	"fmt"
	"text/template"

	"github.com/gin-gonic/gin"
)

func WebPageHandler(c *gin.Context) {
	tpl, err := template.ParseFiles("/home/feng/go/src/my_1m_chatroom/web/home.html")
	if err != nil {
		fmt.Fprint(c.Writer, "模板解析错误！")
		panic(err)
	}

	err = tpl.Execute(c.Writer, nil)
	if err != nil {
		fmt.Fprint(c.Writer, "模板执行错误！")
		return
	}
}
