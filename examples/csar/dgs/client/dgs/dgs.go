// © 2019-2025 Diarkis Inc. All rights reserved.

package dgs

import (
	"fmt"

	pufferDgs "github.com/Diarkis/diarkis-server-template/examples/csar/dgs/puffer/go/dgs"
	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
)

const statusOk = uint8(1)

type DGS struct {
	tcp *tcp.Client
	udp *udp.Client
}

func SetupAsTCP(c *tcp.Client) *DGS {
	r := &DGS{tcp: c}
	r.setup()
	return r
}

func SetupAsUDP(c *udp.Client) *DGS {
	r := &DGS{udp: c}
	r.setup()
	return r
}

func (r *DGS) setup() {
	if r.tcp != nil {
		r.tcp.OnResponse(r.onResponse)
		r.tcp.OnPush(r.onPush)
		return
	}
	if r.udp != nil {
		r.udp.OnResponse(r.onResponse)
		r.udp.OnPush(r.onPush)
	}
}

func (r *DGS) onResponse(ver uint8, cmd uint16, status uint8, payload []byte) {
	fmt.Printf("onResponse ver=%d cmd=%d status=%v, payload=%x\n", ver, cmd, status, payload)

	if ver == pufferDgs.InstanceCreateRequestVer && cmd == pufferDgs.InstanceCreateRequestCmd {
		r.onInstanceCreateResponse(status, payload)
		return
	}
	if ver == pufferDgs.InstanceBackfillRequestVer && cmd == pufferDgs.InstanceBackfillRequestCmd {
		r.onInstanceBackfillResponse(status, payload)
		return
	}
}

func (r *DGS) onInstanceCreateResponse(status uint8, payload []byte) {
	if status != statusOk {
		fmt.Printf("DGS create failed. status %d, response %v\n", status, payload)
		return
	}

	proto := pufferDgs.InstanceCreateResponse{}
	proto.Unpack(payload)
	fmt.Printf("DGS create response: %s\n", proto.String())
}

func (r *DGS) onInstanceBackfillResponse(status uint8, payload []byte) {
	if status != statusOk {
		fmt.Printf("DGS backfill failed. status %d, response %v\n", status, payload)
		return
	}

	proto := pufferDgs.InstanceBackfillResponse{}
	proto.Unpack(payload)
	fmt.Printf("DGS backfill response: %s\n", proto.String())
}

func (r *DGS) onPush(ver uint8, cmd uint16, payload []byte) {
	fmt.Printf("onPush ver=%d cmd=%d payload=%x\n", ver, cmd, payload)

	if ver == pufferDgs.InstanceCreatePushVer && cmd == pufferDgs.InstanceCreatePushCmd {
		r.onInstanceCreatePush(payload)
		return
	}
	if ver == pufferDgs.InstanceBackfillRequestVer && cmd == pufferDgs.InstanceBackfillRequestCmd {
		r.onInstanceBackfillPush(payload)
		return
	}
	if ver == pufferDgs.InstanceAllocatedPushVer && cmd == pufferDgs.InstanceAllocatedPushCmd {
		r.onInstanceAllocatedPush(payload)
		return
	}
}

func (r *DGS) onInstanceCreatePush(payload []byte) {
	proto := pufferDgs.NewInstanceCreatePush()
	if err := proto.Unpack(payload); err != nil {
		fmt.Printf("onInstanceCreatePush: failed to parse push data. %v\n", err)
		return
	}
	if !proto.Success {
		fmt.Println("onInstanceCreatePush: Failed to allocate DGS...")
	}

	fmt.Printf("onInstanceCreatePush: %s\n", proto.String())
}

func (r *DGS) onInstanceBackfillPush(payload []byte) {
	proto := pufferDgs.NewInstanceCreatePush()
	if err := proto.Unpack(payload); err != nil {
		fmt.Printf("onInstanceBackfillPush: failed to parse push data. %v\n", err)
		return
	}
	if !proto.Success {
		fmt.Println("onInstanceBackfillPush: Failed to backfill DGS...")
	}

	fmt.Printf("onInstanceBackfillPush: %s\n", proto.String())
}

func (r *DGS) onInstanceAllocatedPush(payload []byte) {
	proto := pufferDgs.NewInstanceAllocatedPush()
	if err := proto.Unpack(payload); err != nil {
		fmt.Printf("onInstanceAllocatedPush: failed to parse push data. %v\n", err)
		return
	}
	if !proto.Success {
		fmt.Println("onInstanceAllocatedPush: DGS creation failed...")
	}

	fmt.Printf("onInstanceAllocatedPush: %s\n", proto.String())
}

func (r *DGS) Send(ver uint8, cmd uint16, payload []byte) {
	if r.tcp != nil {
		r.tcp.Send(ver, cmd, payload)
		return
	}
	if r.udp != nil {
		r.udp.Send(ver, cmd, payload)
	}
}

func (r *DGS) CreateSessionFromRoom() {
	r.Send(pufferDgs.InstanceCreateRequestVer, pufferDgs.InstanceCreateRequestCmd, pufferDgs.NewInstanceCreateRequest().Pack())
}

func (r *DGS) BackfillSessionFromRoom() {
	r.Send(pufferDgs.InstanceBackfillRequestVer, pufferDgs.InstanceBackfillRequestCmd, pufferDgs.NewInstanceBackfillRequest().Pack())
}
