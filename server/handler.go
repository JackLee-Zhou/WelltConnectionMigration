package server

import (
	"github.com/gin-gonic/gin"
	"myConnect/connect"
)

func RegisterWebHookHandler(c *gin.Context) {

}

func Sub(c *gin.Context) {

}

func ConnectHandler(c *gin.Context) {
	connect.WebSocketHandler(c)
}

func Info(c *gin.Context) {

}
