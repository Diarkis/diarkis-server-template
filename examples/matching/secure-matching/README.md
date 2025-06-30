This project is composed of three servers and one client.
MARS and http server are standard process with only the HTTP server
defining a simple matching profile.
The third server is the udp one that handle the matching command of the client.


# How to build

First, you need to

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
UDP sid         = 49287b84766f4f53ae3abc4b1ef17835
UDP key         = 6feb9e13b9fa442c9413c223fd746638
UDP iv          = 3591d278db47463aa3204a416e01fef6
UDP mac         = 721499f8bf06482ea727301eb5687927
[2025/05/02 04:11:42.831]<UDPCL>        INFO    Local UDP Client started on [::]:56950
[2025/05/02 04:11:42.833]<NET>          INFO    Local IP Addresses. [REDACTED]
[2025/05/02 04:11:42.833]<UDPCL>        INFO    [user-1] UDP connection started 127.0.0.1:7100
[2025/05/02 04:11:42.833]<UDPCL>        INFO    sendLoop started 127.0.0.1:7100
[2025/05/02 04:11:43.034]<CLI>          INFO    Connected UDP
[2025/05/02 04:11:43.034]<CLI>          INFO    start matching with profile:LevelMatch, level:1, rank:1
[2025/05/02 04:11:43.435]<CLI>         DEBUG    UDP onResponse ver=2 cmd=1 status=1 payload=
[2025/05/02 04:11:43.435]<CLI>          INFO    Matching successfully started
[2025/05/02 04:11:43.836]<CLI>         DEBUG    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":["user-2"],"ticketType":1}
[2025/05/02 04:11:43.836]<CLI>          INFO    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":["user-2"],"ticketType":1}
[2025/05/02 04:11:43.836]<CLI>          INFO    matching complete: {OwnerID:user-1 CandidateIDs:[user-2] TicketType:1}
[2025/05/02 04:11:44.037]<CLI>         DEBUG    UDP onResponse ver=1 cmd=224 status=1 payload=OK
[2025/05/02 04:11:44.237]<CLI>         DEBUG    UDP onPush ver=1 cmd=224 payload=hello world
[2025/05/02 04:11:44.237]<CLI>          INFO    received ticket broadcast ver=1 cmd=224 payload=hello world
[2025/05/02 04:11:44.237]<CLI>          INFO    received matching ticket broadcast hello world
[2025/05/02 04:11:44.237]<CLI>          INFO    test is finished, disconnect

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
UDP sid         = 58cb098792f24f8092e1c1437d78ad2d
UDP key         = a2b4ad6aa4c64b638d46f4409902ab3e
UDP iv          = 67c9fe69398447d28f7bbda13a6c6396
UDP mac         = b998b8b1e4a343279d6dabfb3330536c
[2025/05/02 04:12:46.431]<UDPCL>        INFO    Local UDP Client started on [::]:59955
[2025/05/02 04:12:46.433]<NET>          INFO    Local IP Addresses. [REDACTED]
[2025/05/02 04:12:46.434]<UDPCL>        INFO    [user-1] UDP connection started 127.0.0.1:7100
[2025/05/02 04:12:46.434]<UDPCL>        INFO    sendLoop started 127.0.0.1:7100
[2025/05/02 04:12:46.635]<CLI>          INFO    Connected UDP
[2025/05/02 04:12:46.636]<CLI>          INFO    start matching with profile:LevelMatchExact, level:5, rank:3
[2025/05/02 04:12:47.034]<CLI>         DEBUG    UDP onResponse ver=2 cmd=1 status=1 payload=
[2025/05/02 04:12:47.034]<CLI>          INFO    Matching successfully started
[2025/05/02 04:12:47.235]<CLI>         DEBUG    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-2","candidateIDs":["user-1"],"ticketType":1}
[2025/05/02 04:12:47.235]<CLI>          INFO    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-2","candidateIDs":["user-1"],"ticketType":1}
[2025/05/02 04:12:47.235]<CLI>          INFO    matching complete: {OwnerID:user-2 CandidateIDs:[user-1] TicketType:1}
[2025/05/02 04:12:47.636]<CLI>         DEBUG    UDP onPush ver=1 cmd=224 payload=hello world
[2025/05/02 04:12:47.636]<CLI>          INFO    received ticket broadcast ver=1 cmd=224 payload=hello world
[2025/05/02 04:12:47.636]<CLI>          INFO    received matching ticket broadcast hello world
[2025/05/02 04:12:47.636]<CLI>          INFO    test is finished, disconnect
```

