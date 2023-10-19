package p2pcmds

import "github.com/Diarkis/diarkis/p2p"

// Expose exposes the P2P commands
func Expose() {
	p2p.ExposeCommands()
}
