// © 2019-2024 Diarkis Inc. All rights reserved.

package customcmds

import (
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

func resonanceCmd(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Debug("Resonance command has received %#v from the client SID:%s - UID:%s", payload, userData.SID, userData.ID)
	reliable := false
	userData.ServerRespond(payload, ver, cmd, server.Ok, reliable)
	logger.Debug("Resonance command has respond %#v", payload)
	next(nil)
}
