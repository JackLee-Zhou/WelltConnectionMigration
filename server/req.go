package server

type (
	RegisterWebHookReq struct {
		ClientID string `json:"clientId"`
		WebHook  string `json:"webhook"`
	}
)
