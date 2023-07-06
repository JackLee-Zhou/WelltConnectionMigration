package types

import "encoding/json"

type (
	Res struct {
		ID      string `json:"id"`
		JsonRpc string `json:"jsonrpc"`
		Result  string `json:"result"`
	}
)

type Params struct {
	ID      string `json:"id"`
	Topic   string `json:"topic"`
	Message string `json:"message"`

	//TODO 如何处理 TTL 和 Tag 数据是什么
	TTL               uint32 `json:"ttl"` // 有效时间 用于 JWT 中存储有效时间
	Tag               uint32 `json:"tag"`
	PayloadParamsData `json:"data"`
}

type PayloadParamsData struct {
	ID   string `json:"id"`
	Data struct {
		Topic       string `json:"topic"`
		Message     string `json:"message"`
		PublishedAt uint32 `json:"publishedAt"`
		//TODO 如何处理 TTL 和 Tag 数据是什么
		TTL uint32 `json:"ttl"`
		Tag uint32 `json:"tag"`
	} `json:"data"`
}

type (
	RegisterWebHookReq struct {
		ClientID string `json:"clientId"`
		WebHook  string `json:"webhook"`
	}
	// PubSub 通用的 pub 或者 Sub 请求
	PubSub struct {
		ID      string `json:"id"`
		JsonRpc string `json:"jsonrpc"`
		Method  string `json:"method"`
		Params  `json:"params"`
	}
)

func (p Params) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (p PubSub) Marshal() ([]byte, error) {
	return json.Marshal(p)
}

func (r Res) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
