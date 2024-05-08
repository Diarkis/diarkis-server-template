package room

import (
	"github.com/Diarkis/diarkis/user"
)

/*
JoinRoomCallback is callback function for room.Join.

	[IMPORTANT] This callback function is invoked on the server where the user is "originally" connected.
	            room.Join may move the user to a different server
	            if the room that the user is joining is not on the same server.
	            When the user is moved to another server, this callback function is invoked on the server where
	            the user was originally connected.
*/
type JoinRoomCallback func(err error, memberIDs []string, ownerID string, createdTime int64)

/*
JoinReturnData represents internally used data
*/
type JoinReturnData struct {
	MemberUIDs  []string `json:"memberUIDs"`
	OwnerID     string   `json:"ownerID"`
	CreatedTime int64    `json:"createdTime"`
}

/*
Join joins a room and notify the other members of the room on joining the room.

	[NOTE] Uses mutex lock internally.
	[NOTE] If message is empty, Broadcast will not be sent.
	[NOTE] This function is asynchronous. What this means is that it involves server-to-server communication internally.

	[IMPORTANT] The room to join may not be on the same server process and if that is the case, the client will move to where the room is.
	            That means the client will re-connect to a different Diarkis server in the Diarkis cluster.
	            This is handled by Diarkis client SDK. The join will be completed when the user successfully changes the server.

	[IMPORTANT] The callback is called on the original server that the user client is connected to
	            and the room to join may not be on the same server.
	            Please do not read and/or write room properties in the callback as the room may not be accessible
	            and creates a bug in your application code.
	            Do not use functions that require the room to be on the same server to work.

	[IMPORTANT] This function works even if the room is on a different server.

Error Cases

	┌─────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────┐
	│ Error                                   │ Reason                                                                         │
	╞═════════════════════════════════════════╪════════════════════════════════════════════════════════════════════════════════╡
	│ User is already in another room         │ The user is not allowed to join more than one room at a time.                  │
	├─────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Node of room not found                  │ Server of the room is not available in the Diarkis cluster.                    │
	├─────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Failed to transfer user                 │ Transferring the user to the server of the room failed.                        │
	├─────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Room not found                          │ Room to join not found.                                                        │
	├─────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Room join rejected                      │ SetOnJoinByID and/or SetJoinCondition callback(s) rejected the user to join.   │
	├─────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────┤
	│ Not allowed to join without reservation │ The room is room including reservation and the user does not have reservation. │
	└─────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────┘

# Error Evaluations

The possible error evaluations:

Besides the following errors, there are Mesh module errors as well.

	room.IsJoinError(err error) // Failed to join a room

Parameters

	roomID   - Target room ID to join.
	userData - User to join the room.
	ver      - Command version to be used for successful join message sent to other room members.
	cmd      - Command ID to be used for successful join message sent to other room members.
	message  - Message byte array to be sent as successful message to other room members.
	           If message is either nil or empty, the message will not be sent.
	callback - Callback to be invoked when join operation completes (both success and failure).

# Packets sent from the server for join

Join operation on the server may send multiple packets to the client.

Here are the details on each packet sent from the server and their order.

▶︎ Re-connect instruction push

This packet is sent from the server when the server requires the client to re-connect to another server.

▶︎ On join response

This packet is always sent from the server when join operation is completed.

	ver = 1
	cmd = 101

Command ver and ID for on join packet is as shown above:

# Order of packets from the server and join callback

▶︎ With re-connect instruction:

Join request → Re-connect instruction push → On join response → On member join push

	    ┌──────────────────────────┐                        ┌──────────────────┐
	    │     UDP/TCP server A     │                        │ UDP/TCP server B │──────────────┐
	    └──────────────────────────┘                        └──────────────────┘              │
	                 ▲        │                                   │                           │
	<1> Join request │        │                                   │                           │
	                 │        │ <2> Re-connect instruction push   │                           │
	                 │        │                                   │ <3> On join response      │
	            ╭────────╮◀︎───┘                                   │                           │ <4> On member join push
	            │ Client │◀︎───────────────────────────────────────┘                           │
	            ╰────────╯◀︎───────────────────────────────────────────────────────────────────┤
	                  ╭────────╮╮╮                                                            │
	Other room members│ Client │││◀︎───────────────────────────────────────────────────────────┘
	                  ╰────────╯╯╯

▶︎ Without re-connect instruction:

Join request → On join response → On member join push

	    ┌────────────────────────────────────┐
	┌─▶︎ │           UDP/TCP server A         │────┐
	│   └────────────────────────────────────┘    │
	│                      │                      │
	│ <1> Join request     │                      │
	│                      │ <2> On join response │
	│                      ▼                      │ <3> On member join push
	│                  ╭────────╮                 │
	└──────────────────│ Client │◀︎────────────────┤
	                   ╰────────╯                 │
	                   ╭────────╮╮╮               │
	 Other room members│ Client │││◀︎──────────────┘
	                   ╰────────╯╯╯

The order of packet and callback change when re-connect or not.

# Callback and Server Location

When the user joins a room that exists on a different server, the client is required to re-connect to the server where the room is.
Diarkis handles this internally, but it is useful to understand the mechanism behind it.

The diagram below shows the sequence of operations and callbacks as well as where they are executed.

The number in <> indicates the order of execution of callbacks and operations.

	┌──────────────────┐                      ┌────────────────────────┐ <7> On member join push waits for on join response
	│ UDP/TCP server A │                      │    UDP/TCP server B    │─────────────────────────────────────┐
	│                  │  <2> Join operation  │                        │                                     │
	│                  │─────────────────────▶︎│                        │ <6> On join response waits for      │
	│                  │                      │                        │     SetOnJoinCompleteByID execution │
	│                  │◀︎─────────────────────│ [ Room A exists here ] │─────────────────────────────┐       │
	└──────────────────┘ <3> Join operation   └────────────────────────┘                             │       │
	        ▲                callback is              ▲ <5> Callback of SetOnJoinCompleteByID is     │       │
	        │                executed on server A     │     executed on server B                     │       │
	        │                                         │                                              │       │
	        │ <1> Join request                        │ <4> Re-connect                               │       │
	        │                                         │     Client re-connects to where the room is  │       │
	    ╭────────╮                                    │                                              │       │
	    │ Client │────────────────────────────────────┘                                              │       │
	    ╰────────╯◀︎──────────────────────────────────────────────────────────────────────────────────┘       │
	        ▲                                                                                                │
	        │                                                                                                │
	        └────────────────────────────────────────────────────────────────────────────────────────────────┤
	                   ╭────────╮╮╮                                                                          │
	 Other room members│ Client │││◀︎─────────────────────────────────────────────────────────────────────────┘
	                   ╰────────╯╯╯

# Join Sequences

Join operation and its sequences are different when the room to join is on the same server or on a different server.

▶︎ Same Server Join

	The sequence below is join operation sequence where the room is on the same server as the user.

	                                                [ Room is here ]
	          ╭────────╮                              ┌──────────┐
	          │ Client │                              │ Server A │
	          ╰────────╯                              └──────────┘
	               │                                        │
	               ┝━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━▶︎[ Join Request ]
	               │                                        │
	               │                            [ room.SetJoinCondition ]
	               │                                        │
	               │                               ┌────────┴───┐
	               │                               │            │
	       [ Join Failure Response ]◀︎━━━━━[ Room Join Failed ]  │
	               │                                            │
	               │                                            │
	               │                                            │
	               │                                [ room.SetOnJoinByID ]
	               │                                            │
	               │                                   ┌────────┴───┐
	               │                                   │            │
	[ Join Failure Response ]◀︎━━━━━━━━━━━━━━━━[ Room Join Failed ]  │
	               │                                                │
	               │                                                │
	               │                                     [ Reservation Check ]
	               │                                                │
	               │                                       ┌────────┴───┐
	               │                                       │            │
	[ Join Failure Response ]◀︎━━━━━━━━━━━━━━━━━━━━[ Room Join Failed ]  │
	               │                                                    │
	               │                                       [ room.SetOnJoinByID ]
	               │                                                    │
	               │                                           ┌────────┴───┐
	               │                                           │            │
	[ Join Failure Response ]◀︎━━━━━━━━━━━━━━━━━━━━━━━━[ Room Join Failed ]  │
	               │                                                        │
	               │                                             [ room.SetOnRoomChange ]
	               │                                                        │
	               │                                          [ room.SetOnJoinCompleteByID ]
	               │                                                        │
	               │                                                 [ Join Success ]
	               │                                                        │
	     [ On Join Response ]◀︎━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┥
	               │                                                        │
	      [ On Member Join ]◀︎━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┥ On Member Join will be sent to all members.
	               │                                                        │
	               ▽                                                        ▽

▶︎ Remote Server Join

	The sequence below is the operation sequence where the room is on a different server from the user.

	                                         [ Room is here ]
	       ╭────────╮       ┌──────────┐       ┌──────────┐
	       │ Client │       │ Server A │       │ Server B │
	       ╰────────╯       └──────────┘       └──────────┘
	            │                 │                 │
	            ┝━━━━━━━━▶︎[ Join Request ]          │
	            │                 │                 │
	            │                 ├────────▶︎[ Join Request ]
	            │                 │                 │
	            │                 │     [ room.SetJoinCondition ]
	            │                 │                 │
	            │                 │             ┌───┴───┐
	            │                 │             │       │
	            │         [ Join Failure ]◀︎─────┘       │
	            │                 │                     │
	[ Join Failure Response ]◀︎━━━━┥                     │
	            │                 │                     │
	            │                 │          [ room.SetOnJoinByID ]
	            │                 │                     │
	            │                 │                 ┌───┴───┐
	            │                 │                 │       │
	            │         [ Join Failure ]◀︎─────────┘       │
	            │                 │                         │
	[ Join Failure Response ]◀︎━━━━┥                         │
	            │                 │                         │
	            │                 │               [ Reservation Check ]
	            │                 │                         │
	            │                 │                     ┌───┴───┐
	            │                 │                     │       │
	            │         [ Join Failure ]◀︎─────────────┘       │
	            │                 │                             │
	[ Join Failure Response ]◀︎━━━━┥                             │
	            │                 │                             │
	            │                 │                  [ room.SetOnJoinByID ]
	            │                 │                             │
	            │                 │                         ┌───┴───┐
	            │                 │                         │       │
	            │         [ Join Failure ]◀︎─────────────────┘       │
	            │                 │                                 │
	[ Join Failure Response ]◀︎━━━━┥                                 │
	            │                 │                                 │
	            │                 │                    [ room.SetOnRoomChange ]
	            │                 │                                 │
	            │          [ Join Success ]◀︎────────────────────────┤ Callback of room.Join on Server A is invoked.
	            │                 │                                 │
	[ Re-connect Instruction ]◀︎━━━┥ Disconnect from Server A.       │
	            │                 │                                 │
	            ┝━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━▶︎[ Re-connect ]
	            │                 │                                 │
	    [ On Connect ]◀︎━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┥
	            │                 │                                 │
	            │                 │                    [ room.SetOnJoinCompleteByID ]
	            │                 │                                 │
	            │                 │                          [ Join Success ]
	            │                 │                                 │
	  [ On Join Response ]◀︎━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┥
	            │                 │                                 │
	   [ On Member Join ]◀︎━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┥ On Member Join will be sent to all members.
	            │                 │                                 │
	            ▽                 ▽                                 ▽
*/
func Join(roomID string, userData *user.User, ver uint8, cmd uint16, message []byte, callback JoinRoomCallback) {
}
