package room

/*
StartP2PSync Starts peer-to-peer comminucations with the remote clients in the same room.
This method raises OnStartP2PSync event.
Pramms:
- clients: Number of clients to be linked directly via peer-to-peer.Other remote clients in the room will be linked via relay.
*/
func (room *Room) StartP2PSync(clients byte) {
}
