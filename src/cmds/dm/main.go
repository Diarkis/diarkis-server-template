// © 2019-2024 Diarkis Inc. All rights reserved.

package dmcmds

import (
	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis/dm"
	"github.com/Diarkis/diarkis/metrics"
	"github.com/Diarkis/diarkis/user"
	"github.com/Diarkis/diarkis/util"
)

var deliveredCnt *metrics.Counter
var deliveredBytes *metrics.Counter

func Setup() {
	// Create the metrics after mesh.Setup. If you create them before mesh.Setup, Diarkis does not register them.
	diarkis.OnReady(func(next func(error)) {
		deliveredCnt = metrics.NewCounter(metrics.CounterOpts{
			Name:      "dm_delivered_total",
			Help:      "Represents the number of direct messages delivered to the users on a server (node). Disconnect messages are not included.",
			NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
		})
		deliveredBytes = metrics.NewCounter(metrics.CounterOpts{
			Name:      "dm_delivered_bytes_total",
			Help:      "Represents the total size of direct messages delivered to the users on a server (node). Disconnect messages are not included.",
			NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
		})
		next(nil)
	})
	dm.SetOnReceiveMessage(func(ver uint8, cmd uint16, message []byte, _ string, _ *user.User) bool {
		// Do not count the disconnect messages.
		if ver == util.CmdBuiltInVer && cmd == util.CmdDMDisconn {
			return true
		}
		deliveredCnt.Inc()
		deliveredBytes.Add(uint64(len(message)))
		return true
	})

	// If you do not need to send a disconnect message automatically, comment out the following line
	dm.SetOnUserDisconnect(func(disconnectedUser *user.User, peerUserID string) []byte {
		return []byte(disconnectedUser.ID)
	})
}
