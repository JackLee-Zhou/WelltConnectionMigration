package connect

import "sync"

type WebHookList struct {
	ClientID string
	WebHook  string
}

var WebHookListMap = make(map[string][]*WebHookList)
var hookListMu = new(sync.Mutex)

func AddWebHookList(clientID, webHook string) {
	hookListMu.Lock()
	defer hookListMu.Unlock()
	WebHookListMap[clientID] = append(WebHookListMap[clientID], &WebHookList{
		ClientID: clientID,
		WebHook:  webHook,
	})
}
