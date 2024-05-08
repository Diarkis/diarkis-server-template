package server

import (
	"github.com/Diarkis/diarkis/user"
)

/*
NotificationHandler is invoked on every Notify. If it returns true, the user client receives the notification,
but if it returns false, the client user does not receive the notification.
*/
type NotificationHandler func(userData *user.User, notification *Notification) bool

/*
NotificationLoader is a callback to used by NotificationService to generate a notification message.
*/
type NotificationLoader func() (notification *Notification, err error)

/*
Notification is a data structure that represents a notification message.

	ID      - Unique Identifier for the server notification. This ID must not overlap.
	Ver     - Command version of the notification. If ver = 0 and cmd = 400, the client raises OnNotification event.
	Cmd     - Command ID of the notification. If ver = 0 and cmd = 400, the client raises OnNotification event.
	Message - Notification message
	TTL     - Notification TTL in milliseconds.
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
EnableNotifications initializes notifier. It must be called before diarkis.Start is called
to setup notifier properly.
In order for notifications to be sent out to the clients, you must use Notify and/or NotificationService.
*/
func EnableNotifications() {
}

/*
DeleteNotification deletes an existing Notification message struct data.
*/
func DeleteNotification(id string) {
}

/*
SetOnNotification assigns a callback function to notifications with the matching name.
The callback is invoked on every notification with the same name before it is sent to a client and decides who receives the notification
by returning true/false. If the callback returns true, the target user will receive the notification.

	[CRITICALLY IMPORTANT] Using pointer variables that are defined outside of the callback closure
	                       in the callback closure will cause those pointers to be not garbage collected leading to memory leak.
*/
func SetOnNotification(name string, cb NotificationHandler) {
}

/*
NotificationService starts a goroutine loop that invokes the given callback at the interval that is given.
Every iteration of the goroutine loop triggers a notification message being sent.

	[NOTE] The returned value of the callback will be sent as a notification.
	[IMPORTANT] In order to deliver notifications, all server process must start the same notification service.
	[IMPORTANT] This function does NOT guarantee the delivery of notification message to all user clients.
	            The possible failure of notification message delivery may include the following:
	              - Duplicate message ID will cause the message to be sent or ignored.
	              - User clients that are not connected at the time of this function's execution.
	              - User clients that are connected to non-responsive server process.
	              - TTL expiration.

Parameters

	name     - A name of the notification service.
	cb       - A callback to be invoked on each notification service loop tick. The returned values will be the notification.
	interval - Notification delivery cycle interval in milliseconds

Example:

The code example blow demonstrates server.NotificationService fetches a notification data from a database every 10,000ms and
sets it up to be delivered to all connected users.

	interval := int64(5000)

	server.EnableNotifications()

	server.NotificationService("Test", func() (notification *server.Notification, err error) {

		// Retrieve notification data from a database by the current time
		notificationData := someDatabase.GetNotificationDataByCurrentTime(year, month, date)

		if notificationData == nil {
			// No notification data to send out
			return nil, nil
		}

		n := new(server.Notification)
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
		n.TTL = int64(60)

		return n, nil
	}, interval)
*/
func NotificationService(name string, cb NotificationLoader, interval int64) {
}

/*
Notify sends out a notification message for all user clients.

	[IMPORTANT] This function does NOT guarantee the delivery of notification message to all user clients.
	            The possible failure of notification message delivery may include the following:
	              - Duplicate message ID will cause the message to be sent or ignored.
	              - User clients that are not connected at the time of this function's execution.
	              - User clients that are connected to non-responsive server process.
	              - TTL expiration.

Parameters

	name    - A notification name used to assign NotificationHandler.
	m       - Message struct that represents a notification message.
	local   - If true, the notification message will be sent to the user clients on the same server ONLY.

Returns an error if it fails to send out the notification.
*/
func Notify(m *Notification, local bool) error {
	return nil
}
