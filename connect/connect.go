package connect

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"myConnect/tlog"
	"myConnect/types"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { // TODO: remove when using at prod...
		return true
	},
}

const (
	PubType      = "irn_publish"
	BatchPubType = "irn_batchPublish" // 批量发布
	SubType      = "irn_subscribe"
	UnSubType    = "irn_unsubscribe"    // 取消订阅
	BatchSubType = "irn_batchSubscribe" // 批量订阅
)

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
		Silent  bool   `json:"silent"`
	}
	Event struct {
		Topic string `json:"topic"`
		ID    string `json:"id"`
	}
)

type WPool struct {
	sync.Mutex
	nodes map[*websocket.Conn]peer // key : conn value : peer 该连接 对应的订阅和发布
	//events map[*websocket.Conn][]*Event
}

type peer struct {
	pubs map[string]*types.Params
	subs map[string]struct{}
}

var wPool = &WPool{
	nodes: make(map[*websocket.Conn]peer),
	//events: make(map[*websocket.Conn][]*Event),
}

func WebSocketHandler(c *gin.Context) {
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

			switch msg.JsonRpc {
			case "1.0":
				v1Handler(&msg)
			case "2.0":
				v2Handler(&msg, conn)
			default:
				tlog.Errorf("websocket msg version is not exist! %s\n", msg.JsonRpc)
			}

		}
	}(conn)
}

func v1Handler(msg *types.PubSub) {
	switch msg.Method {
	case PubType:
	case SubType:
		//wPool.pub()
	default:
		tlog.Errorf("websocket msg type is not exist! %s\n", msg.Method)
	}
}

func v2Handler(msg *types.PubSub, con *websocket.Conn) {
	switch msg.Method {
	case PubType:
		wPool.pub(con, &msg.Params)
	case BatchPubType:
		// Message 为 PublishedMessage[]
		wPool.batchPub(con, msg.Params.Message)
	case SubType:
		subID := wPool.sub(msg.Params.Topic, msg.ID, msg.JsonRpc, con)
		res := &types.Res{
			ID:      msg.ID,
			JsonRpc: msg.JsonRpc,
			Result:  subID,
		}
		tlog.Infof("sub res is %+v \n", res)
		data, err := json.Marshal(res)
		if err != nil {
			tlog.Errorf("websocket marshal res error: %s topic is %s ", err.Error(), msg.Params.Topic)
			return
		}
		err = con.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			tlog.Errorf("websocket write message error: %s topic is %s ", err.Error(), msg.Params.Topic)
			return
		}
	case BatchSubType:
		wPool.batchSub(msg.Params.Topic, msg.ID, msg.JsonRpc, con)
	case UnSubType:
		wPool.unSub(msg.Params.Topic, msg.Params.ID)
	default:
		tlog.Errorf("websocket msg type is not exist! %s\n", msg.Method)
	}
}

// batchPub 批量推送
func (w *WPool) batchPub(con *websocket.Conn, msgs string) {
	params := []types.Params{}
	if err := json.Unmarshal([]byte(msgs), &params); err != nil {
		tlog.Errorf("websocket unmarshal msg error: %s\n", err.Error())
		return
	}
	for _, v := range params {
		temp := v
		w.pub(con, &temp)
	}
}

// pub 推送
func (w *WPool) pub(con *websocket.Conn, params *types.Params) {
	w.Lock()
	defer w.Unlock()
	// 获取那些连接监听了这个 topic
	cons := w.getSub(params.Topic)
	if cons == nil {
		tlog.Infof("websocket topic is not exist! %s\n", params.Topic)
		return
	}
	for _, v := range cons {
		temp := v
		doPub(temp, params)
	}
	w.setPub(con, params)
}

// setPub 设置推送事件
func (w *WPool) setPub(con *websocket.Conn, params *types.Params) {
	if n, ok := w.nodes[con]; ok {
		n.pubs[params.Topic] = params
	} else {
		w.nodes[con] = peer{
			pubs: map[string]*types.Params{
				params.Topic: params,
			},
			subs: map[string]struct{}{},
		}
	}
}

// getPub 获取待推送事件
func (w *WPool) getPub(topic string) *types.Params {
	for _, v := range w.nodes {
		if p, ok := v.pubs[topic]; ok {
			return p
		}
	}
	return nil
}

// getSub 获取订阅者
func (w *WPool) getSub(topic string) []*websocket.Conn {
	cons := []*websocket.Conn{}
	for con, v := range w.nodes {
		if _, ok := v.subs[topic]; ok {
			cons = append(cons, con)
			continue
		}
	}
	return cons
}

// doPub 开始推送消息
func doPub(con *websocket.Conn, params *types.Params) {
	err := con.WriteMessage(websocket.TextMessage, []byte(params.Message))
	if err != nil {
		tlog.Errorf("websocket write message error: %s topic is %s \n", err.Error(), params.Topic)
		return
	}
}

func (w *WPool) batchSub(topics string, id, rpcVersion string, conn *websocket.Conn) {
	ts := []string{}
	if err := json.Unmarshal([]byte(topics), &ts); err != nil {
		tlog.Errorf("websocket unmarshal msg error: %s\n", err.Error())
		return
	}
	subIDs := []string{}
	for _, topic := range ts {
		temp := topic
		subId := w.sub(temp, id, rpcVersion, conn)
		subIDs = append(subIDs, subId)
	}
	bID, err := json.Marshal(subIDs)
	if err != nil {
		tlog.Errorf("websocket batchSub marshal msg error: %s\n", err.Error())
		return
	}
	res := &types.Res{
		ID:      id,
		JsonRpc: rpcVersion,
		Result:  string(bID),
	}
	data, _ := json.Marshal(res)
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		tlog.Errorf("websocket batchSub write message error: %s topic is %s \n", err.Error(), topics)
		return
	}
}

// sub 注册监听
func (w *WPool) sub(topic, id, rpcVersion string, conn *websocket.Conn) (subID string) {
	w.Lock()
	defer w.Unlock()
	subID = toID(topic)
	// 表明 该 连接  订阅了这个 topic
	if _, ok := w.nodes[conn]; !ok {
		//	分配空间
		w.nodes[conn] = peer{
			pubs: make(map[string]*types.Params),
			subs: map[string]struct{}{},
		}
	}
	w.nodes[conn].subs[topic] = struct{}{}
	msg := w.getPub(topic)

	if msg == nil {
		tlog.Infof("websocket topic pub is not exist! %s\n", topic)
	} else {
		// 推送一下之前存在的信息
		doPub(conn, msg)
	}
	return
}

// unSub 取消监听
func (w *WPool) unSub(topic, id string) {

}

// toId 生成 每个 topic 计算  sha256 生产 id
func toID(topic string) string {
	nhex := sha256.New()
	defer nhex.Reset()
	timeStamp := time.Now().Unix()

	// 加入时间戳偏移
	_, err := nhex.Write([]byte(topic + strconv.FormatInt(timeStamp, 10)))
	if err != nil {
		tlog.Errorf("sha256 write error: %s topic is %s ", err.Error(), topic)
		return ""
	}
	bytes := nhex.Sum(nil)
	res := hex.EncodeToString(bytes)

	return res
}
