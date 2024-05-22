package customcmds

import (
	custom "github.com/Diarkis/diarkis-server-template/puffer/go/custom"
	userstatus "github.com/Diarkis/diarkis-server-template/puffer/go/userstatus"
	sessiondata "github.com/Diarkis/diarkis-server-template/puffer/go/sessiondata"
	"github.com/Diarkis/diarkis-server-template/lib/onlinestatus"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

func getUserStatusList(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	req := custom.NewOnlineStatusRequest()
	req.Unpack(payload)

	logger.Sys("Get user online status list: UIDs:%v", req.UIDs)

	list, err := onlinestatus.GetUserStatusList(req.UIDs)

	if err != nil {
		userData.ServerRespond(util.ErrData(err.Error(), util.ErrCodeInvalidParams), ver, cmd, server.Bad, true)
		next(err)
		return
	}

	resp := custom.NewOnlineStatusResponse()
	resp.UserStatusList = make([]*userstatus.UserStatus, len(list))

	for i := 0; i < len(list); i++ {
		item := list[i]
		us := userstatus.NewUserStatus()
		us.UID = item.UID
		us.InRoom = item.InRoom
		us.SessionData = make([]*sessiondata.UserSessionData, len(item.SessionData))
		for j := 0; j < len(item.SessionData); j++ {
			sd := item.SessionData[j]
			us.SessionData[j] = sessiondata.NewUserSessionData()
			us.SessionData[j].Type = sd.Type
			us.SessionData[j].ID = sd.ID
		}
		resp.UserStatusList[i] = us
	}

	packed := resp.Pack()

	userData.ServerRespond(packed, ver, cmd, server.Ok, true)
	next(nil)
}
