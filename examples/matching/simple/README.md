This project is composed of three servers and one client.
MARS and http server are standard process with only the HTTP server
defining a simple matching profile.
The third server is the udp one that handle the matching command of the client.



# How to build

```sh
./run-mage.sh build:local
```

or

```batch
.\run-mage.bat build:local
```

This will build the three servers and the client binary.
See remote_bin folder.

# How to run

You must first start the MARS server.
```sh
./run-mage.sh server mars
```

Next you can start the HTTP server.
```sh
./run-mage.sh server http
```

Then you can start the UDP server.
```sh
./run-mage.sh server udp
```

Once all the servers are running, you can start two clients to test the matching code.

## Test matching


The test client has the following parameters.

```
Usage of ./remote_bin/cli:
  -clientKey string
        the client key to authenticate with the server
  -host string
        the address of the HTTP server (default "127.0.0.1:7000")
  -profile string
        The matching profile to use. [LevelMatch, LevelMatchExact] (default "LevelMatch")
  -tag value
        The matching tag to use. Can be set multiple times.
  -uid string
        the unique identifier of the client like user ID
  -userLevel int
        The user level (greater or equal to zero) (matching property)  (default 1)
  -userRank int
        The user rank (greater or equal to zero) (default 1)
```


### level bucket

The HTTP server defines a matching profile named `LevelMatch` associated to the property `level`.

```go
	levelMatchProfile := make(map[string]int)
	levelMatchProfile["level"] = 10
	matching.Define("LevelMatch", levelMatchProfile)
```

With this profile, each level bucket will pool users with level 1 to 10, 11 to 20, 21 to 30 and so forth...

| user 1<br>level | user 1<br>rank | user 2<br>level | user 2<br>rank | result                |
| --------------- | -------------- | --------------- | -------------- | --------------------- |
| 1               | 1              | 2               | 3              | can match together    |
| 4               | 1              | 10              | 3              | can match together    |
| 4               | 1              | 11              | 3              | cannot match together |
| 16              | 1              | 10              | 3              | can match together    |


The server defines also a second profile named `LevelMatchExact` where only exact matching level
will be matched together.

#### Matching complete

The current matching uses the user level as the unique criteria to match users together.

```sh
./remote_bin/cli -uid user-1 -userLevel 1 -userRank 1
```

```sh
./remote_bin/cli -uid user-2 -userLevel 4 -userRank 1
```

Example of output
```
Connecting to HTTP server first: http://127.0.0.1:7000/endpoint/type/UDP/user/user-1 - clientKey = 
UDP address = 127.0.0.1:7100
UDP sid         = b459ca40ddf64a908f04ff24f8251279
UDP key         = 642dcb9ad1b341a093cdfe679cb05537
UDP iv          = f9ff7e15319246859dbf1827efe7273b
UDP mac         = d633cdf78a50446caff175fa83e7a2e0
[2024/10/31 06:35:39.698]<UDPCL>        INFO Local UDP Client started on [::]:53810
[2024/10/31 06:35:39.699]<NET>          INFO Local IP Addresses. [REDACTED]
[2024/10/31 06:35:39.699]<UDPCL>        INFO [user-1] UDP connection started 127.0.0.1:7100
[2024/10/31 06:35:39.699]<UDPCL>        INFO sendLoop started 127.0.0.1:7100
Connected UDP
[2024/10/31 06:35:40.301]<CLI>         DEBUG UDP onResponse ver=2 cmd=1 status=1 payload=
[2024/10/31 06:35:42.906]<CLI>         DEBUG UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":[],"ticketType":1}
[2024/10/31 06:35:42.906]<CLI>          INFO UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":[],"ticketType":1}
[2024/10/31 06:35:42.906]<CLI>          INFO matching complete: {OwnerID:user-1 CandidateIDs:[] TicketType:1}
[2024/10/31 06:35:43.307]<CLI>         DEBUG UDP onResponse ver=1 cmd=224 status=1 payload=OK
[2024/10/31 06:35:43.508]<CLI>         DEBUG UDP onPush ver=1 cmd=224 payload=hello world
[2024/10/31 06:35:43.508]<CLI>          INFO received ticket broadcast ver=1 cmd=224 payload=hello world
[2024/10/31 06:35:43.508]<CLI>          INFO received matching ticket broadcast hello world
[2024/10/31 06:35:43.508]<CLI>          INFO test is finished, disconnect
[2024/10/31 06:35:43.508]<UDPCL>      SYSTEM [user-1] Failed to receive a packet from <nil>: read udp [::]:53810: use of closed network connection
[2024/10/31 06:35:43.508]<UDPCL>        INFO [user-1] Client disconnected from 127.0.0.1:7100
```

### exact level

The HTTP server defines a matching profile named `LevelMatchExact` associated to the property `level`.

```go
	levelMatchProfile := make(map[string]int)
	levelMatchProfile["level"] = 1
	matching.Define("LevelMatch", levelMatchProfile)
```

When using this profile only user with the exact same level will be matched together.

| user 1<br>level | user 1<br>rank | user 2<br>level | user 2<br>rank | result                |
| --------------- | -------------- | --------------- | -------------- | --------------------- |
| 1               | 1              | 2               | 3              | cannot match together |
| 4               | 1              | 10              | 3              | cannot match together |
| 4               | 1              | 11              | 3              | cannot match together |
| 16              | 1              | 10              | 3              | cannot match together |
| 5               | 1              | 5               | 3              | can match together    |

#### Matching complete

The current matching uses the user level as the unique criteria to match users together.

```sh
./remote_bin/cli -uid user-1 -userLevel 5 -userRank 3 -profile LevelMatchExact
```

```sh
./remote_bin/cli -uid user-2 -userLevel 5 -userRank 7 -profile LevelMatchExact
```

Example of output

```
Connecting to HTTP server first: http://127.0.0.1:7000/endpoint/type/UDP/user/user-1 - clientKey = 
UDP address = 127.0.0.1:7100
UDP sid         = 084408fc79114e409c25fe3a58e95c3c
UDP key         = ea964564e9e04fd8a4203df7b7cbf74d
UDP iv          = 993b115118054260a21626c85faada1a
UDP mac         = cb6b82e82ab3425983fe3dd0b4eda9e0
[2024/11/01 04:46:12.383]<UDPCL>        INFO Local UDP Client started on [::]:60891
[2024/11/01 04:46:12.384]<NET>          INFO Local IP Addresses. [REDACTED]
[2024/11/01 04:46:12.384]<UDPCL>        INFO [user-1] UDP connection started 127.0.0.1:7100
[2024/11/01 04:46:12.384]<UDPCL>        INFO sendLoop started 127.0.0.1:7100
[2024/11/01 04:46:12.585]<CLI>          INFO Connected UDP
[2024/11/01 04:46:12.585]<CLI>          INFO start matching with profile:LevelMatchExact, level:5, rank:3
[2024/11/01 04:46:12.986]<CLI>         DEBUG UDP onResponse ver=2 cmd=1 status=1 payload=
[2024/11/01 04:46:12.986]<CLI>          INFO Matching successfully started
[2024/11/01 04:46:13.789]<CLI>         DEBUG UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":[],"ticketType":1}
[2024/11/01 04:46:13.789]<CLI>          INFO UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":[],"ticketType":1}
[2024/11/01 04:46:13.789]<CLI>          INFO matching complete: {OwnerID:user-1 CandidateIDs:[] TicketType:1}
[2024/11/01 04:46:14.190]<CLI>         DEBUG UDP onResponse ver=1 cmd=224 status=1 payload=OK
[2024/11/01 04:46:14.391]<CLI>         DEBUG UDP onPush ver=1 cmd=224 payload=hello world
[2024/11/01 04:46:14.391]<CLI>          INFO received ticket broadcast ver=1 cmd=224 payload=hello world
[2024/11/01 04:46:14.391]<CLI>          INFO received matching ticket broadcast hello world
[2024/11/01 04:46:14.391]<CLI>          INFO test is finished, disconnect
[2024/11/01 04:46:14.391]<UDPCL>        INFO [user-1] Client disconnected from 127.0.0.1:7100
```


## Tag ???


FIXME(Henry) fix the log to have non empty candidateIDs
