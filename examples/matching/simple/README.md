# Overview

This project is comprised of **(3)** servers and **(1)** client.

- **MARS** server is a standard (out-of-the-box) Diarkis template. It serves to orchestrates the
  node mesh.

- **HTTP** server is a standard (out-of-the-box) Diarkis template; excepting (2) simple Matchmaker
  profile definitions for our custom matchmaking criteria. Our Matchmaker candidate information is
  stored here.

- **UDP** server hosts the client connection and handles all incoming Matchmaker commands.
  It queries the Matchmaker storage server (**HTTP**) for valid candidates. If a valid
  matching is found via the provided pooling constraints it attempts match the selected candidates.

The goal of this sample is to demostrate how to set-up a simple matchmaking scenario using Diarkis
Matchmaker. We will demonstrate how to use custom pooling constraints to match candidates together.

**NOTE**: For illustration purposes we introduce **(2)** properties `level` and `rank` to illustrate
how to perform basic matchmaking via a pooling constratint. However, for the sake of simplicity,
pooling for this example is only performed via `level`. We encourage you, as an exercise, to
implement a separate pooling-constraint upon on our templated `rank` constraint to master this
topic.

## How to Build

You can build all **(3)** servers and the client binary using the provided Mage build scripts.

### To build on Linux or macOS

```sh
./run-mage.sh build:local
```

### To build on Windows

```sh
.\run-mage.bat build:local
```

This will create the **MARS**, **HTTP**, and **UDP** server binaries, and client binary,
placing them inside the `remote_bin` directory.

## How to Run

This project requires all **(3)** servers to be running before clients can test matchmaking
behavior.

### 1. First, start the MARS server to orchestrate the node mesh

```sh
./run-mage.sh server mars
```

### 2. Next, start the HTTP server, which holds the custom matchmaking criteria

```sh
./run-mage.sh server http
```

### 3. Then, start the UDP server, to manage incoming client connections

```sh
./run-mage.sh server udp
```

Once all servers are running, you may start two client instances to test and observe the matchmaking
behavior using the custom pooling constraints.

## Testing Matchmaker

The test client has the following parameters:

```output
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

## Buckets and Pooling

### LevelMatch

On the **HTTP** server, we define a matching profile called `LevelMatch` associated to the property `level`.

```go
	levelMatchProfile := make(map[string]int)
	levelMatchProfile["level"] = 10
	matching.Define("LevelMatch", levelMatchProfile)
```

With this profile, each level bucket will pool users by the value of their **level** property in static intervals of `10`. E.g:

```example
[1–10], [11–20], [21–30], ..., [n–(n+9)] // and so on...
```

| User 1 (`level`) | User 2 (`level`) | Match Outcome |
|:-----------------|:-----------------|:--------------|
| `1`              | `2`              | OK            |
| `4`              | `10`             | OK            |
| `4`              | `11`             | FAIL          |
| `16`             | `11`             | OK            |

### Use Example

To test our `LevelMatch` pooling constraint, execute the following on **(2)** separate instances of
the provided client:

```sh
./remote_bin/cli -uid user-1 -userLevel 1 -userRank 1 -profile LevelMatch
```

```sh
./remote_bin/cli -uid user-2 -userLevel 4 -userRank 1 -profile LevelMatch
```

As our candidates are within the level range for our pooling constraint `LevelMatch` they will
successfully match them together.

**Output:**

```output
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



### LevelMatchExact

On the **HTTP** server, we also define a second profile called `LevelMatchExact` which only requires
that only candidates of an exactly equivalent `level` property may be matched together.

```go
	levelMatchProfile := make(map[string]int)
	levelMatchProfile["level"] = 1
	matching.Define("LevelMatch", levelMatchProfile)
```
| User 1 (`level`) | User 2 (`level`) | Match Outcome |
|:-----------------|:-----------------|:--------------|
| `1`              | `2`              | FAIL          |
| `4`              | `10`             | FAIL          |
| `4`              | `11`             | FAIL          |
| `16`             | `11`             | FAIL          |
| `5`              | `5`              | OK            |
| `30`             | `30`             | OK            |

### Use Example

To test our `LevelMatchExact` pooling constraint, execute the following on **(2)** separate instances of
the provided client:


```sh
./remote_bin/cli -uid user-1 -userLevel 5 -userRank 3 -profile LevelMatchExact
```

```sh
./remote_bin/cli -uid user-2 -userLevel 5 -userRank 7 -profile LevelMatchExact
```

As our candidates are both the same level, our pooling constraint `LevelMatchExact` will allow them
to successfully match them together.

**Output:**

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

---

*First created on 2025-04-03. Updated on 2025-04-04.*
