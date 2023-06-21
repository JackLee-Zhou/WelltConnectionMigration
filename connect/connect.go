package connect

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"myConnect/tlog"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // TODO: remove when using at prod...
		return true
	},
}

type (
	Msg struct {
		Version string `json:"version"`
	}
	MsgV2 struct {
		ID      string `json:"id"`
		Topic   string `json:"topic"`
		Message string `json:"message"`
		TTL     uint32 `json:"ttl"`
		Tag     uint32 `json:"tag"`
	}
	MsgV1 struct {
		Topic   string `json:"topic"`
		Type    string `json:"type"`
		Payload string `json:"payload"`
		//Silent  bool   `json:"silent"`
	}
	Event struct {
	}
)

type WPool struct {
	sync.Locker
	nodes  map[string]*websocket.Conn // key : Topic value : conn
	events map[*websocket.Conn][]*Event
}

func WebSocketHandler(c *gin.Context) {
	if !c.IsWebsocket() {
		c.String(400, "Bad Request")
		return
	}
	// 协议升级 为 DApp 和 Wallet 建立链接
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	//conn, _, _, err := ws.UpgradeHTTP(c.Request, c.Writer)
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
			var msg Msg
			if err = json.Unmarshal(msgBt, &msg); err != nil {
				tlog.Errorf("websocket unmarshal msg error: %s\n", err.Error())
				break
			}

			/**
			tlog.Infof("incoming: %s\n", msg.String())
			switch msg.Type {
			case "subscribe":
				subscribeController(con, msg.Topic)
			case "publish":
				publishController(con, msg)
			}
			*/

			switch msg.Version {
			case "v1":
				v1Handler(&MsgV1{})
			case "v2":
				v2Handler(&MsgV2{})
			default:
				tlog.Errorf("websocket msg version is not exist! %s\n", msg.Version)
			}

		}
	}(conn)
}

func v1Handler(msg *MsgV1) {

}

func v2Handler(msg *MsgV2) {

}

// Pub 推送
func (w *WPool) Pub() {
	w.Locker.Lock()
	defer w.Locker.Unlock()
}

// Sub 注册监听
func (w *WPool) Sub() {
	w.Locker.Lock()
	defer w.Locker.Unlock()
}
