This directory contains a sample that relies on dive to keep track
of online users.

## Build Servers for Local Use

The following command will be building the servers for your local machine.

```sh
./run-mage.sh build:local
```

or on Windows

```bat
.\run-mage.bat build:local
```

## Start the servers

First of all you need to start diarkis mars server.

```sh
./run-mage.sh server mars
```

Then you can start the http server.

```sh
./run-mage.sh server http
```

And the udp server.

```sh
./run-mage.sh server udp
```

## Connect one or more client to the server

```
./remote_bin/cli -uid user-1
```

Let the user-2 create a room on connection.
```
./remote_bin/cli -uid user-2 -create-room
```

## Check online status

You can retrieve the online status of the users by visiting the url below.
[http://localhost:7000/onlinestatus/uids/&lt;user-id&gt;](http://localhost:7000/onlinestatus/uids/<user-id>).

In order to retrieve the online status of the users user-1 and user-2 you connected
to the server using the `remote_bin/cli` binary.
[http://localhost:7000/onlinestatus/uids/user-1,user-2](http://localhost:7000/onlinestatus/uids/user-1,user-2)

```json
{"user-1":{"InRoom":false,"SessionData":{}},"user-2":{"InRoom":true,"SessionData":{}}}
```

## How does this work

In [lib/onlinestatus](lib/onlinestatus/main.go) the Setup function
registers a callback to be called when a user establishes a connection
to the server.
See the call to `user.OnNew`.
The setup function also registers a callback on keep alive in
order to refresh the client's online status while the client is connected.
```go
	diarkis.OnReady(func(next func(error)) {
		server.OnKeepAlive(updateUserStatus)
		next(nil)
	})
```

The http server exposes an endpoint to retrieve the online status of one or more
users.
