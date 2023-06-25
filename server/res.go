package server

import "encoding/json"

type (
	Res struct {
		ID      string `json:"id"`
		JsonRpc string `json:"jsonrpc"`
		Result  string `json:"result"`
	}
)

func (r Res) Marshal() ([]byte, error) {
	return json.Marshal(r)
}
