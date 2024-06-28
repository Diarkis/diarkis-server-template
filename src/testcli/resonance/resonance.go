// © 2019-2024 Diarkis Inc. All rights reserved.

package resonance

import (
	"fmt"

	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
)

const (
	version        uint8  = 2
	resonanceCmdID uint16 = 13
)

type Resonance struct {
	tcp *tcp.Client
	udp *udp.Client
}

func SetupAsTCP(c *tcp.Client) *Resonance {
	r := &Resonance{tcp: c}
	r.setup()
	return r
}

func SetupAsUDP(c *udp.Client) *Resonance {
	r := &Resonance{udp: c}
	r.setup()
	return r
}

func (r *Resonance) setup() {
	if r.tcp != nil {
		r.tcp.OnResponse(r.onResponse)
		return
	}
	if r.udp != nil {
		r.udp.OnResponse(r.onResponse)
	}
}

func (r *Resonance) onResponse(ver uint8, cmd uint16, status uint8, payload []byte) {
	if ver != version || cmd != resonanceCmdID {
		return
	}

	fmt.Printf("Resonance command response: %v\n", string(payload))
}

func (r *Resonance) Send(cmd uint16, payload []byte) {
	if r.tcp != nil {
		r.tcp.Send(version, resonanceCmdID, payload)
		return
	}
	if r.udp != nil {
		r.udp.Send(version, resonanceCmdID, payload)
	}
}

func (r *Resonance) Resonate(message string) {
	r.Send(resonanceCmdID, []byte(message))
}
