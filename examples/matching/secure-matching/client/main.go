// © 2019-2025 Diarkis Inc. All rights reserved.

package main

import (
	"github.com/Diarkis/diarkis/client/go/test/cli"
)

func main() {
	cli.RegisterCommands("ticket", cli.TicketCommands)
	cli.Connect()
	cli.Run()
}
