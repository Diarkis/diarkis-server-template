// © 2019-2024 Diarkis Inc. All rights reserved.
package httpcmds

// rootpath is defined in cmds/main.go

import (
	"encoding/json"
	"strings"

	"github.com/Diarkis/diarkis/server/http"

	"github.com/Diarkis/diarkis-server-template/lib/onlinestatus"
)

func exposeOnlineStatus() {
	// $(uids) is a comma separated user ID list
	// e.i. /onlinestatus/uids/1,2,3,4
	http.Get("/onlinestatus/uids/:uids", getOnlineStatusList)
}

func getOnlineStatusList(resp *http.Response, req *http.Request, params *http.Params, next func(error)) {
	src, err := params.GetAsString("uids")

	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}

	uids := strings.Split(src, ",")

	list, err := onlinestatus.GetUserStatusList(uids)

	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}

	// key is UID
	ret := make(map[string]interface{}, 0)

	for _, userStatusData := range list {
		uid := userStatusData.UID
		inRoom := userStatusData.InRoom
		// key value map of type as key and ID as value per session
		sessionData := make(map[uint8]string)
		for _, sd := range userStatusData.SessionData {
			sessionData[sd.Type] = sd.ID
		}
		// populate the map
		v := make(map[string]interface{})
		v["InRoom"] = inRoom
		v["SessionData"] = sessionData
		ret[uid] = v
	}

	encoded, err := json.Marshal(ret)

	if err != nil {
		resp.Respond(err.Error(), http.Bad)
		next(err)
		return
	}

	resp.Respond(string(encoded), http.Ok)
	next(nil)
}
