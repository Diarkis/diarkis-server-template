package room

import (
	"github.com/Diarkis/diarkis/user"
)

/*
SetOnMigrated registers a callback on migrated.
The callbacks will be invoked on the server node where the room is migrated to

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.

	callback - A callback to be invoked when migration is successful.
*/
func SetOnMigrated(callback func(previousRoomID string, currentRoomID string)) {
}

/*
MigrateRoom migrates a room to another node and moves its members to the new room.

This is intended to be used when the current node is about to shutdown (in offline state),
but you need to keep the current room.

	[IMPORTANT] This function does not work if the room is not on the same server.
	[IMPORTANT] Migration will change the room ID.
	[IMPORTANT] The callback may be invoked on the server where the room is not held.
	            Do not use functions that requires the room to be on the same server.

	[NOTE] Uses mutex lock internally.
	[NOTE] This function is asynchronous.

Error Cases

	┌─────────────────────────────────────────┬────────────────────────────────────────────────────────────────────┐
	│ Error                                   │ Reason                                                             │
	╞═════════════════════════════════════════╪════════════════════════════════════════════════════════════════════╡
	│ Node (server) to migrate room not found │ The server to migrate the room to is not found.                    │
	│ Must be owner                           │ Only the room owner is allowed to perform room migration.          │
	│ User transfer failed                    │ The owner user failed to transfer to another server for migration. │
	└─────────────────────────────────────────┴────────────────────────────────────────────────────────────────────┘

Parameters

	ver      - Command ver to send the new migrated room ID to all members as migration instruction.
	cmd      - Command ID to send the new migrated room ID to all members as migration instruction.
	userData - Owner user to perform the migration.
	callback - Invoked when migration operation is completed.

When a room is migrated to another server, the users of the room will be re-connected to the server where the room is migrated.
This is handled by Diarkis internally, but it is useful to understand the mechanism of it.

The diagram below shows the sequence of operations and callbacks as well as where they are executed.

The number in < > indicates the order of execution of callbacks and operations.

	┌────────────────────────┐                   ┌────────────────────────┐
	│ UDP/TCP server A       │  <2> Migrate      │    UDP TCP server B    │
	│                        │──────────────────▶︎│                        │
	│ [ Room A exists here ] │                   │  [ Room A migrated ]   │ <3> SetOnMigrated is executed on server B
	│                        │                   │                        │     Re-connections will take place
	└────────────────────────┘                   └────────────────────────┘     AFTER the execution of SetOnMigrated
	        ▲                                      │                 ▲
	        │                                      │                 │
	        │ <1> Migrate                          │                 │ <6> Re-connect
	        │                                      │                 │     (All members of the room will reconnect asynchronously)
	        │ ┌────────────────────────────────────┘                 │
	        │ │ <6> Re-connect instruction                           │
	        │ │     (All members of the room will receive            │
	        │ │      instruction asynchronously)                     │
	        │ ▼                                                      │
	    ╭────────╮                                                   │
	    │ Client │───────────────────────────────────────────────────┘
	    ╰────────╯(Client re-connects to where the migrated room is)

# Owner user:

The owner of the room will NOT receive on join event response (cmd:101) from the server.
Instead, the owner will receive on migration complete event push (cmd:20) from the server
when the owner user completes the server re-connection.

	  ╭────────╮                 ┌──────────┐      ┌──────────┐
	  │ Client │                 │ Server A │      │ Server B │
	  ╰────────╯                 └──────────┘      └──────────┘
	      │                           │                 │
	      ┝━━━━━━━━━━━━━━━━━▶︎[ Migration Request ]      │
	      │                           │                 │
	      │                           ├───────────▶︎[ Migration ] Room is copied here with the new room ID.
	      │                           │                 │
	      │                           │ [ SetOnMigrated Callback Invoked ]
	      │                           │                 │
	      │           [ Migration Confirmation ]◀︎───────┤
	      │                           │                 │
	[ Migration Instruction ]◀︎━━━━━━━━┥(ver: custom)    │
	      │                           │(cmd: custom)    │
	      │                           │                 │
	      ┝━━━━━━━━━━━━▶︎[ Leave From Previous Room ]    │
	      │                           │                 │
	      │                           │                 │
	[ Leave Confirmation ]◀︎━━━━━━━━━━━┥                 │
	      │                           │                 │
	      ┝━━━━━━━▶︎[ Join Request To Migrated Room ]    │
	      │                           │                 │
	      │                           ├────────────▶︎[ Join ]
	      │                           │                 │
	      │              [ Join Confirmation ]◀︎─────────┤
	      │                           │                 │
	[ Re-connect Instruction ]◀︎━━━━━━━┥                 │
	      │                           │                 │
	      ┝━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━▶︎[ Re-connect ]
	      │                           │                 │
	      │                           │                 │
	[ Migration Confirmation ]◀︎━━━━━━━━━━━━━━━━━━━━━━━━━┥ For owner user: cmd=20 on migration complete
	      │                           │                 │
	      ▽                           ▽                 ▽

# Non-owner user:

Non-owner user WILL receive on join event response (cmd:101) from the server when completing the server re-connection.

	  ╭────────╮                 ┌──────────┐      ┌──────────┐
	  │ Client │                 │ Server A │      │ Server B │
	  ╰────────╯                 └──────────┘      └──────────┘
	      │                           │                 │
	      ┝━━━━━━━━━━━━━━━━━▶︎[ Migration Request ]      │
	      │                           │                 │
	      │                           ├───────────▶︎[ Migration ] Room is copied here with the new room ID.
	      │                           │                 │
	      │                           │ [ SetOnMigrated Callback Invoked ]
	      │                           │                 │
	      │           [ Migration Confirmation ]◀︎───────┤
	      │                           │                 │
	[ Migration Instruction ]◀︎━━━━━━━━┥(ver: custom)    │
	      │                           │(cmd: custom)    │
	      │                           │                 │
	      ┝━━━━━━━━━━━━▶︎[ Leave From Previous Room ]    │
	      │                           │                 │
	      │                           │                 │
	[ Leave Confirmation ]◀︎━━━━━━━━━━━┥                 │
	      │                           │                 │
	      ┝━━━━━━━▶︎[ Join Request To Migrated Room ]    │
	      │                           │                 │
	      │                           ├────────────▶︎[ Join ]
	      │                           │                 │
	      │              [ Join Confirmation ]◀︎─────────┤
	      │                           │                 │
	[ Re-connect Instruction ]◀︎━━━━━━━┥                 │
	      │                           │                 │
	      ┝━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━▶︎[ Re-connect ]
	      │                           │                 │
	      │                           │                 │
	[ Join Confirmation ]◀︎━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┥ Treat join confirmation as migration confirmation.
	      │                           │                 │
	      ▽                           ▽                 ▽
*/
func MigrateRoom(ver uint8, cmd uint16, userData *user.User, callback func(err error, newRoomID string, oldRoomID string)) {
}
