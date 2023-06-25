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
	TTL     uint32 `json:"ttl"`
	Tag     uint32 `json:"tag"`
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
