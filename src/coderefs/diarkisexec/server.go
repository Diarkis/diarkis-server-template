package diarkisexec

import (
	"github.com/Diarkis/diarkis/user"
)

/*
Notification represents the notification data to be propagated to all connected Diarkis clients.
*/
type Notification struct {
	ID      string
	Name    string
	Ver     uint8
	Cmd     uint16
	Message []byte
	TTL     int64
}

/*
SetupNotificationService starts a notification service loop on the server.

The notification loop will execute the given callback at the given interval in seconds.

If the callback returns, *diarkisexec.Notification,
the server will send the notification message along with ver and cmd to all the clients connected to the server.

	[IMPORTANT] The server with HTTP role is not allowed to use Notifier.
	[IMPORTANT] In order to deliver notifications, all server process must start the same notification service.
	[IMPORTANT] This function does NOT guarantee the delivery of notification message to all user clients.
	            The possible failure of notification message delivery may include the following:
	              - Duplicate message ID will cause the message to be sent or ignored.
	              - User clients that are not connected at the time of this function's execution.
	              - User clients that are connected to non-responsive server process.
	              - TTL expiration.

	[NOTE]      The returned value of the callback will be sent as a notification.

Parameters

	name     - A name of the notification service.
	cb       - A callback to be invoked on each notification service loop tick. The returned values will be the notification.
	           If the callback returns a nil, notification will not be sent to the clients.
	interval - Notification delivery cycle interval in milliseconds

Example:

The code example blow demonstrates diarkisexec.SetupNotificationService fetches a notification data from a database every 60s and
sets it up to be delivered to all connected users.

	// Notification service will poll at every 60 seconds
	interval := int64(60) // 60 seconds

	diarkisexec.SetupNotificationService("Test", interval, func() (notification *diarkisexec.Notification, err error) {

		// Retrieve notification data from a database by the current time
		notificationData := someDatabase.GetNotificationDataByCurrentTime(year, month, date)

		if notificationData == nil {
			// No notification data to send out
			return nil, nil
		}

		n := &diarkisexec.Notification{}
		n.ID = notificationData.ID
		n.Name = notificationData.Name

		// Ver is used by the client to identify the message when received.
		n.Ver = notificationData.Ver

		// Cmd is used by the client to identify the message when received.
		n.Cmd = notificationData.Cmd

		n.Message = []byte("Notification message says 'Hello from the server'")

		// TTL is in seconds to indicate the availability of the notification data.
		// The notification will be available for the not-connected-clients for the duration of TTL and
		// will be sent to the clients when they connect before TTL expires.
		n.TTL = int64(60 * 60) // one hour

		return n, nil

	})
*/
func SetupNotificationService(name string, interval int64, callback func() (*Notification, error)) {
}

/*
ResponseStatusOK returns UDP or TCP server response status code for successfully handled command.
*/
func ResponseStatusOK() uint8 {
	return 0
}

/*
ResponseStatusBad returns UDP or TCP server response status code for invalid command invocation.
*/
func ResponseStatusBad() uint8 {
	return 0
}

/*
ResponseStatusErr returns UDP or TCP server response status code for server error while handling command.
*/
func ResponseStatusErr() uint8 {
	return 0
}

/*
OnKeepAlive assigns a callback to be invoked on every echo (UDP) or heartbeat (TCP).

		[IMPORTANT] Every callback must call next func(error) at the end of the operation
	             to allow Diarkis to move on to the next on keep alive operations.
*/
func OnKeepAlive(handler func(userData *user.User, next func(error))) {
}

/*
MarkServerAsOnline flags the server as "ONLINE".

	[NOTE] "ONLINE" is the default state of the server.
*/
func MarkServerAsOnline() {
}

/*
MarkServerAsOffline flags the server as "OFFLINE".

	[NOTE] "OFFLINE" server does not allow the following operations on the server:
	       - Accept new user creation and connection.
	       - Diarkis Room's new room creation.
	       - Diarkis Group's new group creation.
	       - Diarkis Session's new session creation.
	[NOTE] "OFFLINE" server still functions as normal except for the prohibited operations listed above.
*/
func MarkServerAsOffline() {
}

/*
MarkServerAsTaken flags the server as "TAKEN".

	[NOTE] "TAKEN" server does not allow the following operation on the server:
	       - Accept new user creation and connection.
	[NOTE] "TAKEN" server still functions as normal except for the prohibited operations listed above.
*/
func MarkServerAsTaken() {
}

/*
MarkServerAsTakenIf flags the server as "TAKEN" if the callback returns true.
The callback is invoked every 2 second and if the callback returns true,
the server will be marked as "ONLINE".

	[IMPORTANT] If the server is marked as "OFFLINE", this function will be ignored.

	[NOTE]      This function is executed every 2 seconds so there will be race condition and it is not precise.

TAKEN - When the server is marked as "TAKEN", the server will NOT accept new user connections.

Example: The example code below uses CCU of the node to control TAKEN <==> ONLINE

	server.MarkServerAsTakenIf(func() bool {
		// user package can tell you the CCU of the node
		ccu := user.GetCCU()
		if ccu >= maxAllowedCCU {
			// this will mark the server as TAKEN
			return true
		}
		// this will mark the server as ONLINE
		return false
	})
*/
func MarkServerAsTakenIf(callback func() bool) {
}
