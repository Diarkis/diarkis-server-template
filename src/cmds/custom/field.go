package customcmds

import (
	dpayload "{0}/lib/payload"

	"github.com/Diarkis/diarkis/field"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

// add to matching and create a room
func getFieldInfo(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Sys("Get Field Info Received from user: {}", userData.ID)
	fieldInfo := dpayload.NewGetFieldInfo()

	fieldInfo.NodeCount = int32(field.GetNodeNum())
	fieldInfo.FieldSize = field.GetFieldSize()
	fieldInfo.FieldOfVisionSize = field.GetFieldOfVisionSize()
	userData.ServerRespond(fieldInfo.Pack(), ver, cmd, server.Ok, true)
	next(nil)
}
