package customcmds

import (
	"github.com/Diarkis/diarkis-server-template/lib/onlinestatus"
	ponlinestatus "github.com/Diarkis/diarkis-server-template/puffer/go/onlinestatus"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

func getUserStatusList(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	req := ponlinestatus.NewOnlineStatusRequest()
	req.Unpack(payload)

	logger.Sys("Get user online status list: UIDs:%v", req.UIDs)

	list, err := onlinestatus.GetUserStatusList(req.UIDs)

	if err != nil {
		userData.ServerRespond(util.ErrData(err.Error(), util.ErrCodeInvalidParams), ver, cmd, server.Bad, true)
		next(err)
		return
	}

	resp := ponlinestatus.NewOnlineStatusResponse()
	resp.UserStatusList = make([]*ponlinestatus.UserStatus, len(list))

	for i := 0; i < len(list); i++ {
		item := list[i]
		us := ponlinestatus.NewUserStatus()
		us.UID = item.UID
		us.InRoom = item.InRoom
		us.SessionData = make([]*ponlinestatus.UserSessionData, len(item.SessionData))
		for j := 0; j < len(item.SessionData); j++ {
			sd := item.SessionData[j]
			us.SessionData[j] = ponlinestatus.NewUserSessionData()
			us.SessionData[j].Type = sd.Type
			us.SessionData[j].ID = sd.ID
		}
		resp.UserStatusList[i] = us
	}

	packed := resp.Pack()

	userData.ServerRespond(packed, ver, cmd, server.Ok, true)
	next(nil)
}
