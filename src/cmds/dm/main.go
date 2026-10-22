// © 2019-2024 Diarkis Inc. All rights reserved.

package dmcmds

import (
	"github.com/Diarkis/diarkis"
	"github.com/Diarkis/diarkis/dm"
	"github.com/Diarkis/diarkis/metrics"
	"github.com/Diarkis/diarkis/user"
)

var deliveredCnt *metrics.Counter
var deliveredBytes *metrics.Counter

func Setup() {
	// Create the metrics after mesh.Setup. If you create them before mesh.Setup, Diarkis does not register them.
	diarkis.OnReady(func(next func(error)) {
		deliveredCnt = metrics.NewCounter(metrics.CounterOpts{
			Name:      "dm_delivered_total",
			Help:      "Represents the number of direct messages delivered to the users on a server (node).",
			NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
		})
		deliveredBytes = metrics.NewCounter(metrics.CounterOpts{
			Name:      "dm_delivered_bytes_total",
			Help:      "Represents the total size of direct messages delivered to the users on a server (node).",
			NodeRoles: []diarkis.NodeRole{diarkis.UDP, diarkis.TCP},
		})
		next(nil)
	})
	dm.SetOnReceiveMessage(func(_ uint8, _ uint16, message []byte, _ string, _ *user.User) bool {
		deliveredCnt.Inc()
		deliveredBytes.Add(uint64(len(message)))
		return true
	})

	// If you do not need to send a disconnect message automatically, comment out the following line
	dm.SetOnUserDisconnect(func(disconnectedUser *user.User, peerUserID string) []byte {
		return []byte(disconnectedUser.ID)
	})
}
