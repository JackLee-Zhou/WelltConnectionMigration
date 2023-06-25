package server

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"myConnect/tlog"
	"time"
)

func Init() {
	r := gin.Default()
	r.Use(ginzap.Ginzap(tlog.Logger, time.RFC3339, true))

	r.Use(ginzap.RecoveryWithZap(tlog.Logger, true))

	r.POST("/register-webhook", RegisterWebHookHandler)
	r.GET("/", ConnectHandler)
	r.POST("/subscribe", Sub)

	// TODO 要启动 Https
	r.Run(":8000")
}
