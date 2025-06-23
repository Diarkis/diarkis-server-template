// © 2019-2024 Diarkis Inc. All rights reserved.

package cmds

import (
	httpcmds "github.com/Diarkis/diarkis-server-template/examples/dive/user-online-status/cmds/http"
	"github.com/Diarkis/diarkis-server-template/examples/dive/user-online-status/lib/onlinestatus"
)

func SetupUDP() {
	onlinestatus.Setup()
}

func SetupHTTP() {
	httpcmds.Expose()
	onlinestatus.Setup()
}
