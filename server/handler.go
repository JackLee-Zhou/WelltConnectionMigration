package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"myConnect/connect"
	"myConnect/tlog"
	"myConnect/types"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // TODO: remove when using at prod...
		return true
	},
}

// RegisterWebHookHandler 注册 webhook
// 注册 webhook 的话就是在使用 http 来进行通信
func RegisterWebHookHandler(c *gin.Context) {
	var req RegisterWebHookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, map[string]string{"message": "missing or invalid request notification"})
		return
	}
	connect.AddWebHookList(req.ClientID, req.WebHook)
	c.JSON(200, map[string]bool{"message": true})

}

func Sub(c *gin.Context) {

}

func ConnectHandler(c *gin.Context) {
	if !c.IsWebsocket() {
		c.String(400, "Bad Request")
		return
	}
	// 协议升级 为 DApp 和 Wallet 建立链接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		tlog.Errorf("websocket upgrade error: %s\n", err.Error())
		return
	}
	tlog.Infof("success to peer connection: %s\n", conn.RemoteAddr())
	go func(con *websocket.Conn) {
		//	读取 conn 中的数据
		for {
			_, msgBt, err := con.ReadMessage()
			if err != nil {
				tlog.Errorf("websocket read client data error: %s\n", err.Error())
				return
			}
			// 这里做版本判断
			var msg types.PubSub
			if err = json.Unmarshal(msgBt, &msg); err != nil {
				tlog.Errorf("websocket unmarshal msg error: %s\n", err.Error())
				break
			}
			// TODO 不是这样判断的
			switch msg.JsonRpc {
			case "1.0":
				connect.V1Handler(&msg)
			case "2.0":
				connect.V2Handler(&msg, conn)
			default:
				tlog.Errorf("websocket msg version is not exist! %s\n", msg.JsonRpc)
			}
		}
	}(conn)
}

func Info(c *gin.Context) {

}
