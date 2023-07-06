package server

type (
	RegisterWebHookReq struct {
		// 这个怎么用 v1 是用 topic 来标识的
		ClientID string `json:"clientId" binding:"required"`
		WebHook  string `json:"webhook" binding:"required"`
	}
)
