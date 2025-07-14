// © 2019-2024 Diarkis Inc. All rights reserved.

package cmds

import (
	customcmds "github.com/Diarkis/diarkis-server-template/examples/csar/dgs/cmds/custom"
)

func SetupUDP() {
	customcmds.Expose()
}

func SetupTCP() {
	customcmds.Expose()
}

func SetupHTTP() {
}
