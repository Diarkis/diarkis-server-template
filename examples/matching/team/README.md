This project is composed of three servers and one client.
MARS and http server are standard process with only the HTTP server
defining a simple matching profile.
The third server is the udp one that handle the matching command of the client.
The goal of this sample is to explain how implement team matching using two types
of matching ticket.

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
  -uid string
        the unique identifier of the client like user ID
```


./remote_bin/cli -uid user-1
```
Connecting to HTTP server first: http://127.0.0.1:7000/endpoint/type/UDP/user/user-1 - clientKey = 
UDP address = 127.0.0.1:7100
UDP sid         = 2c55fad0047c4979b418a0b36279a964
UDP key         = 3caf704cab1843c98b94b2bdc1667890
UDP iv          = 4dcd68b2ad494af5bcced5fa9ad9cb77
UDP mac         = 027e54e102e547809eba6169228f49c9
[2025/05/09 04:45:18.498]<UDPCL>        INFO    Local UDP Client started on [::]:47288
[2025/05/09 04:45:18.498]<NET>          INFO    Local IP Addresses. [REDACTED]
[2025/05/09 04:45:18.498]<UDPCL>        INFO    [user-1] UDP connection started 127.0.0.1:7100
[2025/05/09 04:45:18.498]<UDPCL>        INFO    sendLoop started 127.0.0.1:7100
[2025/05/09 04:45:18.699]<CLI>          INFO    Connected UDP
[2025/05/09 04:45:18.699]<CLI>          INFO    start matching
[2025/05/09 04:45:19.100]<CLI>         DEBUG    UDP onResponse ver=2 cmd=1 status=1 payload=
[2025/05/09 04:45:19.100]<CLI>          INFO    Matching successfully started
[2025/05/09 04:45:20.103]<CLI>         DEBUG    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":["user-1","user-3"],"ticketType":1,"teams":null}
[2025/05/09 04:45:20.103]<CLI>          INFO    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":["user-1","user-3"],"ticketType":1,"teams":null}
[2025/05/09 04:45:20.103]<CLI>          INFO    matching complete: {OwnerID:user-1 CandidateIDs:[user-1 user-3] TicketType:1 Teams:[]}
[2025/05/09 04:45:20.103]<CLI>          INFO    Team created. Members ["user-1" "user-3"]
[2025/05/09 04:45:20.304]<CLI>         DEBUG    UDP onPush ver=2 cmd=3 payload={"ownerID":"user-1","membersIDs":["user-1","user-3"],"ticketType":1}
[2025/05/09 04:45:22.111]<CLI>         DEBUG    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:22.111]<CLI>          INFO    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:22.111]<CLI>          INFO    matching complete: {OwnerID:user-1 CandidateIDs:[] TicketType:2 Teams:[[user-1 user-3] [user-2 user-4]]}
[2025/05/09 04:45:22.312]<CLI>         DEBUG    UDP onPush ver=2 cmd=3 payload={"ownerID":"user-1","membersIDs":["user-2","user-1"],"ticketType":2}
[2025/05/09 04:45:22.512]<CLI>         DEBUG    UDP onPush ver=2 cmd=4 payload={"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:22.512]<CLI>          INFO    Team matched: {"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:22.512]<CLI>          INFO    test is finished, disconnect
[2025/05/09 04:45:22.513]<UDPCL>      SYSTEM    [user-1] Failed to receive a packet from <nil>: read udp [::]:47288: use of closed network connection
[2025/05/09 04:45:22.513]<UDPCL>        INFO    [user-1] Client disconnected from 127.0.0.1:7100
```

./remote_bin/cli -uid user-2
```
Connecting to HTTP server first: http://127.0.0.1:7000/endpoint/type/UDP/user/user-2 - clientKey = 
UDP address = 127.0.0.1:7100
UDP sid         = 84348f1fe06b4132aa5d8b3c114dc68a
UDP key         = df39ed6292384c44b55dd5c8df4ed0cb
UDP iv          = e1c6dba63f9c4523b8a7dbe7fdd9c687
UDP mac         = e8f6968e7e5f4aa8bfa7908b393640a5
[2025/05/09 04:45:20.459]<UDPCL>        INFO    Local UDP Client started on [::]:43287
[2025/05/09 04:45:20.459]<NET>          INFO    Local IP Addresses. [REDACTED]
[2025/05/09 04:45:20.459]<UDPCL>        INFO    [user-2] UDP connection started 127.0.0.1:7100
[2025/05/09 04:45:20.459]<UDPCL>        INFO    sendLoop started 127.0.0.1:7100
[2025/05/09 04:45:20.660]<CLI>          INFO    Connected UDP
[2025/05/09 04:45:20.661]<CLI>          INFO    start matching
[2025/05/09 04:45:21.061]<CLI>         DEBUG    UDP onResponse ver=2 cmd=1 status=1 payload=
[2025/05/09 04:45:21.061]<CLI>          INFO    Matching successfully started
[2025/05/09 04:45:22.064]<CLI>         DEBUG    UDP onPush ver=2 cmd=3 payload={"ownerID":"user-2","membersIDs":["user-2","user-4"],"ticketType":1}
[2025/05/09 04:45:22.265]<CLI>         DEBUG    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-2","candidateIDs":["user-2","user-4"],"ticketType":1,"teams":null}
[2025/05/09 04:45:22.265]<CLI>          INFO    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-2","candidateIDs":["user-2","user-4"],"ticketType":1,"teams":null}
[2025/05/09 04:45:22.265]<CLI>          INFO    matching complete: {OwnerID:user-2 CandidateIDs:[user-2 user-4] TicketType:1 Teams:[]}
[2025/05/09 04:45:22.265]<CLI>          INFO    Team created. Members ["user-2" "user-4"]
[2025/05/09 04:45:22.466]<CLI>         DEBUG    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:22.466]<CLI>          INFO    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:22.466]<CLI>          INFO    matching complete: {OwnerID:user-1 CandidateIDs:[] TicketType:2 Teams:[[user-1 user-3] [user-2 user-4]]}
[2025/05/09 04:45:22.667]<CLI>         DEBUG    UDP onPush ver=2 cmd=3 payload={"ownerID":"user-1","membersIDs":["user-2","user-1"],"ticketType":2}
[2025/05/09 04:45:22.867]<CLI>         DEBUG    UDP onPush ver=2 cmd=4 payload={"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:22.867]<CLI>          INFO    Team matched: {"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:22.867]<CLI>          INFO    test is finished, disconnect
[2025/05/09 04:45:22.868]<UDPCL>      SYSTEM    [user-2] Failed to receive a packet from <nil>: read udp [::]:43287: use of closed network connection
[2025/05/09 04:45:22.868]<UDPCL>        INFO    [user-2] Client disconnected from 127.0.0.1:7100
```

./remote_bin/cli -uid user-3
```
Connecting to HTTP server first: http://127.0.0.1:7000/endpoint/type/UDP/user/user-3 - clientKey = 
UDP address = 127.0.0.1:7100
UDP sid         = 085da0e072534313a54ab5a5ace94b38
UDP key         = e2286442678a455db251edcbf8daa98d
UDP iv          = 57828ed4d2564ccca3443addc4e9ffb3
UDP mac         = bce7b50580d7439d8917db3edd335111
[2025/05/09 04:45:19.545]<UDPCL>        INFO    Local UDP Client started on [::]:60819
[2025/05/09 04:45:19.546]<NET>          INFO    Local IP Addresses. [REDACTED]
[2025/05/09 04:45:19.546]<UDPCL>        INFO    [user-3] UDP connection started 127.0.0.1:7100
[2025/05/09 04:45:19.546]<UDPCL>        INFO    sendLoop started 127.0.0.1:7100
[2025/05/09 04:45:19.747]<CLI>          INFO    Connected UDP
[2025/05/09 04:45:19.747]<CLI>          INFO    start matching
[2025/05/09 04:45:20.147]<CLI>         DEBUG    UDP onResponse ver=2 cmd=1 status=1 payload=
[2025/05/09 04:45:20.147]<CLI>          INFO    Matching successfully started
[2025/05/09 04:45:20.348]<CLI>         DEBUG    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":["user-1","user-3"],"ticketType":1,"teams":null}
[2025/05/09 04:45:20.348]<CLI>          INFO    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-1","candidateIDs":["user-1","user-3"],"ticketType":1,"teams":null}
[2025/05/09 04:45:20.348]<CLI>          INFO    matching complete: {OwnerID:user-1 CandidateIDs:[user-1 user-3] TicketType:1 Teams:[]}
[2025/05/09 04:45:20.348]<CLI>          INFO    Team created. Members ["user-1" "user-3"]
[2025/05/09 04:45:20.549]<CLI>         DEBUG    UDP onPush ver=2 cmd=3 payload={"ownerID":"user-1","membersIDs":["user-1","user-3"],"ticketType":1}
[2025/05/09 04:45:21.952]<CLI>         DEBUG    UDP onPush ver=2 cmd=4 payload={"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:21.952]<CLI>          INFO    Team matched: {"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:21.952]<CLI>          INFO    test is finished, disconnect
[2025/05/09 04:45:21.952]<UDPCL>      SYSTEM    [user-3] Failed to receive a packet from <nil>: read udp [::]:60819: use of closed network connection
[2025/05/09 04:45:21.952]<UDPCL>        INFO    [user-3] Client disconnected from 127.0.0.1:7100
```

./remote_bin/cli -uid user-4
```
Connecting to HTTP server first: http://127.0.0.1:7000/endpoint/type/UDP/user/user-4 - clientKey = 
UDP address = 127.0.0.1:7100
UDP sid         = 2dba29dccab749c484c294b018abcf9f
UDP key         = 2cfc78eb9b4b4d2488d6af8696b2465f
UDP iv          = eaa4a83372934179a7fc2dd9041cfdb5
UDP mac         = e11952ea64934b5395b42d8f3455937a
[2025/05/09 04:45:21.517]<UDPCL>        INFO    Local UDP Client started on [::]:55037
[2025/05/09 04:45:21.517]<NET>          INFO    Local IP Addresses. [REDACTED]
[2025/05/09 04:45:21.517]<UDPCL>        INFO    [user-4] UDP connection started 127.0.0.1:7100
[2025/05/09 04:45:21.517]<UDPCL>        INFO    sendLoop started 127.0.0.1:7100
[2025/05/09 04:45:21.718]<CLI>          INFO    Connected UDP
[2025/05/09 04:45:21.718]<CLI>          INFO    start matching
[2025/05/09 04:45:22.119]<CLI>         DEBUG    UDP onResponse ver=2 cmd=1 status=1 payload=
[2025/05/09 04:45:22.119]<CLI>          INFO    Matching successfully started
[2025/05/09 04:45:22.320]<CLI>         DEBUG    UDP onPush ver=2 cmd=3 payload={"ownerID":"user-2","membersIDs":["user-2","user-4"],"ticketType":1}
[2025/05/09 04:45:22.520]<CLI>         DEBUG    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-2","candidateIDs":["user-2","user-4"],"ticketType":1,"teams":null}
[2025/05/09 04:45:22.520]<CLI>          INFO    UDP onPush ver=1 cmd=220 payload={"ownerID":"user-2","candidateIDs":["user-2","user-4"],"ticketType":1,"teams":null}
[2025/05/09 04:45:22.521]<CLI>          INFO    matching complete: {OwnerID:user-2 CandidateIDs:[user-2 user-4] TicketType:1 Teams:[]}
[2025/05/09 04:45:22.521]<CLI>          INFO    Team created. Members ["user-2" "user-4"]
[2025/05/09 04:45:22.721]<CLI>         DEBUG    UDP onPush ver=2 cmd=4 payload={"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:22.721]<CLI>          INFO    Team matched: {"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-3"],["user-2","user-4"]]}
[2025/05/09 04:45:22.721]<CLI>          INFO    test is finished, disconnect
[2025/05/09 04:45:22.721]<UDPCL>      SYSTEM    [user-4] Failed to receive a packet from <nil>: read udp [::]:55037: use of closed network connection
[2025/05/09 04:45:22.721]<UDPCL>        INFO    [user-4] Client disconnected from 127.0.0.1:7100
```