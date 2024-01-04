package customcmds

import (
	"math"

	"github.com/Diarkis/diarkis/field"
	"github.com/Diarkis/diarkis/server"
	"github.com/Diarkis/diarkis/user"
)

// add to matching and create a room
func getFieldInfo(ver uint8, cmd uint16, payload []byte, userData *user.User, next func(error)) {
	logger.Sys("Get Field Info Received from user: {}", userData.ID)
	fieldInfo := NewGetFieldInfo()

	fieldInfo.GridCount = int64(field.GetNumberOfGrids())
	fieldInfo.FieldSize = int64(float64(field.GetGridSize()) * math.Sqrt(float64(fieldInfo.GridCount)))
	userData.ServerRespond(fieldInfo.Pack(), ver, cmd, server.Ok, true)
	next(nil)
}
