# Overview

This project is comprised of **(3)** servers and **(1)** client.

- **MARS** server is a standard (out-of-the-box) Diarkis template. It serves to orchestrate the
  node mesh.

- **HTTP** server is a standard (out-of-the-box) Diarkis template; excepting a simple Matchmaker
  profile definition for our matchmaking. Our Matchmaker candidate information is stored here.
  It also hosts Dive and DM, which are used to deliver Session invitations between servers.

- **UDP** server hosts the client connection, handles the Session (party) commands, and handles all
  incoming Matchmaker commands. It queries the Matchmaker storage server (**HTTP**) for valid
  candidates.

The goal of this sample is to demonstrate how to implement team-based (3 vs 3) matchmaking by party
using Diarkis Matchmaker. Users form a party with Diarkis Session, and the party owner starts
matchmaking for the whole party with `matching.StartTicketWithOptions` and `TicketMembers`.

## How it works

Matchmaking is done in **(2)** steps with **(2)** ticket types.

| Ticket type               | Max members      | Purpose                                      |
| ------------------------- | ---------------- | -------------------------------------------- |
| `common.TeamTicketType`   | 3 (1 team)       | Gathers users into a team of 3               |
| `common.BattleTicketType` | 6 (2 teams of 3) | Matches a team against another team (3 vs 3) |

#### 1. Form a party with Session

- The party owner creates a Session and invites the other members.
- Invited users join the Session. When all of them have joined, the party owner starts matchmaking.
- A user who is not in a party starts matchmaking alone (solo).

#### 2. Team matchmaking via `common.TeamTicketType`

- The server gets the party members and the mesh address of the server each member is connected
  to with `session.GetAllSessionMemberData`.
- The party owner (or the solo user) starts a team ticket with `matching.StartTicketWithOptions`.
  The party members are given as `TicketMembers`, so they join the ticket from the start.
  Only the remaining seats are filled by other users.
- If the party already fills a team (3 members), the team ticket is skipped and the party goes
  straight to battle matchmaking.

#### 3. Battle matchmaking via `common.BattleTicketType`

- Once the team is complete, the team owner starts a battle ticket. All the team members are given
  as `TicketMembers`, so the battle ticket contains all the users of both teams.
- The mesh address of a team member who joined the team ticket alone (not a party member) is read
  from the ticket: each user stores the mesh address of its server as a user property
  (`meshAddr`) when it connects, and the property is copied to the ticket.

#### 4. Notify the result with `matching.TicketBroadcast`

- When the battle ticket is complete, the server broadcasts the result, including both teams, to all
  the battle ticket members with `matching.TicketBroadcast`.
- Since all the users of both teams are in the battle ticket, all **(6)** users receive the result.

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

### 2. Next, start the HTTP server, which holds the matchmaking data

```sh
./run-mage.sh server http
```

### 3. Then, start the UDP server, to manage incoming client connections

Start the UDP server twice in separate terminals to run **(2)** UDP servers.
The second one listens on port `7101`.

```sh
./run-mage.sh server udp
```

```sh
./run-mage.sh server udp
```

> All the scenarios also work with **(1)** UDP server.
> With **(2)** UDP servers, the members of a party may be connected to different servers.

Once all servers are running, you may start the client instances to test and observe the
matchmaking behavior.

## Test matching

The test client has the following parameters:

```output
Usage of ./remote_bin/cli:
  -clientKey string
        the client key to authenticate with the server
  -host string
        the address of the HTTP server (default "127.0.0.1:7000")
  -invite string
        create a party and invite the given user IDs (comma separated). Matching starts when all of them join
  -solo
        start matching alone
  -uid string
        the unique identifier of the client like user ID
```

The client runs automatically. When it connects, the user is in the "lobby", and the command line
options decide what the user does:

| Option                  | Action                                                                                  |
| ----------------------- | --------------------------------------------------------------------------------------- |
| `-invite user-2,user-3` | Creates a party, invites the users, and starts matching when all of them have joined    |
| `-solo`                 | Starts matching alone                                                                   |
| (none)                  | Waits for an invitation and joins the party. Exits if no invitation comes in 30 seconds |

### Team-Based Matchmaking Scenario

This scenario uses **(6)** clients to cover the following **(3)** patterns:

| Pattern         | Users                  | Matchmaking                                        |
| --------------- | ---------------------- | -------------------------------------------------- |
| 3-member party  | user-1, user-2, user-3 | Skips team matchmaking and goes straight to battle |
| 2-member party  | user-4, user-5         | Team matchmaking, waits for 1 more member          |
| Solo (no party) | user-6                 | Team matchmaking, joins the team of user-4         |

Start the clients in separate terminals in the following order.
**Start the invited users first**, because an invitation is not delivered to a user who is not
connected yet.

```sh
./remote_bin/cli -uid user-2
```

```sh
./remote_bin/cli -uid user-3
```

```sh
./remote_bin/cli -uid user-5
```

```sh
./remote_bin/cli -uid user-1 -invite user-2,user-3
```

```sh
./remote_bin/cli -uid user-4 -invite user-5
```

```sh
./remote_bin/cli -uid user-6 -solo
```

All **(6)** users receive the battle result with the two teams:
`[["user-1","user-2","user-3"],["user-4","user-5","user-6"]]`.

**Output (user-1, party owner):**

```output
Connecting to HTTP server first: http://127.0.0.1:7000/endpoint/type/UDP/user/user-1 - clientKey =
UDP address = 127.0.0.1:7100
UDP sid         = d64692cca96a4517ad641b990890f57f
UDP key         = af2dfbe273424928a15f1896d88db999
UDP iv          = 0aedc288437849a0916851e039357daa
UDP mac         = c04440ccc5114a89bed2e0d48bc2919d
[2026/11/25 06:37:34.038]<UDPCL>         INFO    Local UDP Client started on [::]:64311
[2026/11/25 06:37:34.038]<UDPCL>         INFO    Local IP Addresses. [REDACTED]
[2026/11/25 06:37:34.038]<UDPCL>         INFO    [user-1] UDP connection started 127.0.0.1:7100
[2026/11/25 06:37:34.038]<UDPCL>         INFO    sendLoop started 127.0.0.1:7100
[2026/11/25 06:37:34.240]<CLI>           INFO    Connected UDP
[2026/11/25 06:37:34.240]<CLI>           INFO    create a party
[2026/11/25 06:37:34.642]<CLI>          DEBUG    UDP onResponse {"ver":1,"cmd":801,"status":1,"payload":"\u0000a26c4a27af0846758e150d1274e3c6a47f0000011fa5d2587350"}
[2026/11/25 06:37:34.642]<CLI>           INFO    party created {"sessionID":"a26c4a27af0846758e150d1274e3c6a47f0000011fa5d2587350","invitees":["user-2","user-3"]}
[2026/11/25 06:37:35.044]<CLI>          DEBUG    UDP onResponse {"ver":1,"cmd":802,"status":1,"payload":"OK"}
[2026/11/25 06:37:35.245]<CLI>          DEBUG    UDP onPush {"ver":1,"cmd":804,"payload":"\u0000user-2"}
[2026/11/25 06:37:35.245]<CLI>           INFO    party member joined {"memberID":"user-2"}
[2026/11/25 06:37:35.446]<CLI>          DEBUG    UDP onPush {"ver":1,"cmd":804,"payload":"\u0000user-3"}
[2026/11/25 06:37:35.446]<CLI>           INFO    party member joined {"memberID":"user-3"}
[2026/11/25 06:37:35.446]<CLI>           INFO    all party members joined
[2026/11/25 06:37:35.446]<CLI>           INFO    start matching
[2026/11/25 06:37:35.847]<CLI>          DEBUG    UDP onResponse {"ver":2,"cmd":1,"status":1,"payload":""}
[2026/11/25 06:37:35.847]<CLI>           INFO    Matching successfully started
[2026/11/25 06:37:41.073]<CLI>          DEBUG    UDP onPush {"ver":2,"cmd":3,"payload":"{\"ownerID\":\"user-1\",\"membersIDs\":[\"user-3\",\"user-6\",\"user-4\",\"user-5\",\"user-1\",\"user-2\"],\"ticketType\":2}"}
[2026/11/25 06:37:42.078]<CLI>          DEBUG    UDP onPush {"ver":2,"cmd":5,"payload":"{\"ownerID\":\"user-1\",\"candidateIDs\":null,\"ticketType\":2,\"teams\":[[\"user-1\",\"user-2\",\"user-3\"],[\"user-4\",\"user-5\",\"user-6\"]]}"}
[2026/11/25 06:37:42.078]<CLI>           INFO    Team matched: {"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-2","user-3"],["user-4","user-5","user-6"]]}
[2026/11/25 06:37:42.078]<CLI>           INFO    test is finished, disconnect
[2026/11/25 06:37:42.079]<UDPCL>       SYSTEM    [user-1] Failed to receive a packet from <nil>: read udp [::]:64311: use of closed network connection
[2026/11/25 06:37:42.279]<UDPCL>         INFO    [user-1] Client disconnected from 127.0.0.1:7100
```

**Output (user-2, invited party member):**

```output
Connecting to HTTP server first: http://127.0.0.1:7000/endpoint/type/UDP/user/user-2 - clientKey =
UDP address = 127.0.0.1:7100
UDP sid         = 017f2b0800174e36997aa188c4c63bc4
UDP key         = b43ce9a630a24be3829b6d915a51f527
UDP iv          = 508968ae293d44859c470f6a9b75659f
UDP mac         = ae653391f3944479a539137d7b87b59c
[2026/11/25 06:37:23.516]<UDPCL>         INFO    Local UDP Client started on [::]:65116
[2026/11/25 06:37:23.516]<UDPCL>         INFO    Local IP Addresses. [REDACTED]
[2026/11/25 06:37:23.516]<UDPCL>         INFO    [user-2] UDP connection started 127.0.0.1:7100
[2026/11/25 06:37:23.516]<UDPCL>         INFO    sendLoop started 127.0.0.1:7100
[2026/11/25 06:37:23.718]<CLI>           INFO    Connected UDP
[2026/11/25 06:37:23.718]<CLI>           INFO    waiting for an invitation {"timeout":"30s"}
[2026/11/25 06:37:34.979]<CLI>          DEBUG    UDP onPush {"ver":1,"cmd":802,"payload":"\u0000a26c4a27af0846758e150d1274e3c6a47f0000011fa5d2587350party"}
[2026/11/25 06:37:34.979]<CLI>           INFO    invited to a party, join {"sessionID":"a26c4a27af0846758e150d1274e3c6a47f0000011fa5d2587350"}
[2026/11/25 06:37:35.381]<CLI>          DEBUG    UDP onResponse {"ver":1,"cmd":804,"status":1,"payload":"OK"}
[2026/11/25 06:37:35.581]<CLI>          DEBUG    UDP onPush {"ver":1,"cmd":804,"payload":"\u0000user-2"}
[2026/11/25 06:37:35.581]<CLI>           INFO    party member joined {"memberID":"user-2"}
[2026/11/25 06:37:35.782]<CLI>          DEBUG    UDP onPush {"ver":1,"cmd":804,"payload":"\u0000user-3"}
[2026/11/25 06:37:35.782]<CLI>           INFO    party member joined {"memberID":"user-3"}
[2026/11/25 06:37:41.007]<CLI>          DEBUG    UDP onPush {"ver":2,"cmd":3,"payload":"{\"ownerID\":\"user-1\",\"membersIDs\":[\"user-3\",\"user-6\",\"user-4\",\"user-5\",\"user-1\",\"user-2\"],\"ticketType\":2}"}
[2026/11/25 06:37:42.012]<CLI>          DEBUG    UDP onPush {"ver":2,"cmd":5,"payload":"{\"ownerID\":\"user-1\",\"candidateIDs\":null,\"ticketType\":2,\"teams\":[[\"user-1\",\"user-2\",\"user-3\"],[\"user-4\",\"user-5\",\"user-6\"]]}"}
[2026/11/25 06:37:42.012]<CLI>           INFO    Team matched: {"ownerID":"user-1","candidateIDs":null,"ticketType":2,"teams":[["user-1","user-2","user-3"],["user-4","user-5","user-6"]]}
[2026/11/25 06:37:42.012]<CLI>           INFO    test is finished, disconnect
[2026/11/25 06:37:42.012]<UDPCL>       SYSTEM    [user-2] Failed to receive a packet from <nil>: read udp [::]:65116: use of closed network connection
[2026/11/25 06:37:42.213]<UDPCL>         INFO    [user-2] Client disconnected from 127.0.0.1:7100
```

---

_First created on 2025-07-03. Updated on 2026-11-25._
