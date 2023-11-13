package fieldlib

import (
	"encoding/binary"

	"github.com/Diarkis/diarkis/client/go/tcp"
	"github.com/Diarkis/diarkis/client/go/udp"
	"github.com/Diarkis/diarkis/packet"
	"github.com/Diarkis/diarkis/util"
)

const builtInVer = util.CmdBuiltInVer
const syncInitCmd = util.CmdFieldSyncInit
const syncCmd = util.CmdFieldSync
const reconnCmd = util.CmdFieldReconn
const disappearCmd = util.CmdFieldDisappear
const serverSyncCmd = util.CmdFieldServerSync

type Field struct {
	tcp                *tcp.Client
	udp                *udp.Client
	onResponseSyncInit []func(messageList [][]byte)
	onSync             []func(message []byte)
	onReconnect        []func()
	onDisappear        []func(uid string)
	onServerSync       []func(inSight bool, message []byte)
}

func NewFieldAsTCP(tcpClient *tcp.Client) *Field {
	f := new(Field)
	f.SetupAsTCP(tcpClient)
	return f
}

func NewFieldAsUDP(udpClient *udp.Client) *Field {
	f := new(Field)
	f.SetupAsUDP(udpClient)
	return f
}

func (f *Field) SetupAsTCP(tcpClient *tcp.Client) bool {
	if f.tcp == nil && f.udp == nil {
		f.tcp = tcpClient
		f.setup()
		return true
	}
	return false
}

func (f *Field) SetupAsUDP(udpClient *udp.Client) bool {
	if f.tcp == nil && f.udp == nil {
		f.udp = udpClient
		f.setup()
		return true
	}
	return false
}

func (f *Field) setup() {
	if f.tcp != nil {
		f.tcp.OnReconnect(f.dispatchOnReconnectCallbacks)
		f.tcp.OnResponse(f.dispatchOnResponseCallbacks)
		f.tcp.OnPush(f.dispatchOnPushCallbacks)
		return
	}
	if f.udp != nil {
		f.udp.OnReconnect(f.dispatchOnReconnectCallbacks)
		f.udp.OnResponse(f.dispatchOnResponseCallbacks)
		f.udp.OnPush(f.dispatchOnPushCallbacks)
		return
	}
}

func (f *Field) SyncInit(x, y, z int64, syncLimit uint16, customFilterID uint8, msg []byte, reliable bool, uid string) {
	payload := f.createSyncPayload(x, y, z, syncLimit, customFilterID, msg, reliable, uid)
	f.send(builtInVer, syncInitCmd, payload, reliable)
}

func (f *Field) Sync(x, y, z int64, syncLimit uint16, customFilterID uint8, msg []byte, reliable bool, uid string) {
	payload := f.createSyncPayload(x, y, z, syncLimit, customFilterID, msg, reliable, uid)
	f.send(builtInVer, syncCmd, payload, reliable)
}

func CreateHeader(oid string) []byte {
	header := []byte{0, 0, uint8(len(oid))}
	oidBytes := []byte(oid)
	header = append(header[:], oidBytes[:]...)
	return header
}

func (f *Field) createSyncPayload(
	x, y, z int64,
	syncLimit uint16,
	customFilterID uint8,
	msg []byte,
	reliable bool,
	uid string) []byte {

	payload := make([]byte, (8*3)+2+1+1)
	binary.BigEndian.PutUint64(payload[0:8], uint64(x))
	binary.BigEndian.PutUint64(payload[8:16], uint64(y))
	binary.BigEndian.PutUint64(payload[16:24], uint64(z))
	binary.BigEndian.PutUint16(payload[24:26], syncLimit)
	payload[26] = byte(customFilterID)
	reliableByte := byte(0x01)
	if !reliable {
		reliableByte = byte(0x00)
	}
	header := CreateHeader(uid)
	payload[27] = reliableByte
	msg = append(header[:], msg[:]...)
	return append(payload[:], msg[:]...)
}

func (f *Field) Disappear() {
	f.send(builtInVer, disappearCmd, make([]byte, 0), true)
}

func (f *Field) send(ver uint8, cmd uint16, payload []byte, reliable bool) {
	if f.tcp != nil {
		f.tcp.Send(ver, cmd, payload)
		return
	}
	if f.udp != nil && reliable {
		f.udp.RSend(ver, cmd, payload)
		return
	}
	if f.udp != nil && !reliable {
		f.udp.Send(ver, cmd, payload)
		return
	}
}

func (f *Field) OnResponseSyncInit(cb func([][]byte)) {
	f.onResponseSyncInit = append(f.onResponseSyncInit, cb)
}

func (f *Field) OnSync(cb func([]byte)) {
	f.onSync = append(f.onSync, cb)
}

func (f *Field) OnReconnect(cb func()) {
	f.onReconnect = append(f.onReconnect, cb)
}

func (f *Field) OnDisappear(cb func(string)) {
	f.onDisappear = append(f.onDisappear, cb)
}

func (f *Field) OnServerSync(cb func(bool, []byte)) {
	f.onServerSync = append(f.onServerSync, cb)
}

func (f *Field) dispatchOnReconnectCallbacks() {
	for _, cb := range f.onReconnect {
		cb()
	}
}

func (f *Field) dispatchOnResponseCallbacks(ver uint8, cmd uint16, status uint8, payload []byte) {
	if ver != builtInVer {
		return
	}
	switch cmd {
	case syncInitCmd:
		f.dispatchOnResponseSyncInit(payload)
	}
}

func (f *Field) dispatchOnPushCallbacks(ver uint8, cmd uint16, payload []byte) {
	if ver != builtInVer {
		return
	}
	switch cmd {
	case syncCmd:
		f.dispatchOnSync(payload)
	case disappearCmd:
		f.dispatchOnDisappear(payload)
	case serverSyncCmd:
		f.dispatchOnServerSync(payload)
	}
}

func (f *Field) dispatchOnResponseSyncInit(payload []byte) {
	messageList := packet.BytesToBytesList(payload)
	for _, cb := range f.onResponseSyncInit {
		cb(messageList)
	}
}

func (f *Field) dispatchOnSync(payload []byte) {
	for _, cb := range f.onSync {
		cb(payload)
	}
}

func (f *Field) dispatchOnDisappear(payload []byte) {
	uid := string(payload)
	for _, cb := range f.onDisappear {
		cb(uid)
	}
}

func (f *Field) dispatchOnServerSync(payload []byte) {
	inSight := true
	if payload[0] == 0x00 {
		inSight = false
	}
	message := payload[1:]
	for _, cb := range f.onServerSync {
		cb(inSight, message)
	}
}
