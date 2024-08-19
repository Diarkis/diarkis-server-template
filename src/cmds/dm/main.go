// © 2019-2024 Diarkis Inc. All rights reserved.

package dmcmds

import (
	"github.com/Diarkis/diarkis/dm"
	"github.com/Diarkis/diarkis/user"
)

func Setup() {
	// If you do not need to send a disconnect message automatically, comment out the following line
	dm.SetOnUserDisconnect(func(disconnectedUser *user.User, peerUserID string) []byte {
		return []byte(disconnectedUser.ID)
	})
}
