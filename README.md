# Overview

Diarkis server cluster is made up with HTTP, TCP, UDP and WebSocket servrers.

Each protocol servers run independently within the cluster, but you do not have to have all protocols.

Only HTTP server is required in the cluster and the rest of the servers should be chosen according to your application's requirements.

# How To Initialize Server Project

You may initialize your Diarkis server project from this repository:

```
make init project_id={project ID} builder_token={build token} output={absolute path to install the server project}
```

**NOTE**: To get your porject ID and builder token, please contact us at https://diarkis.io/en/contact

# Structure

```
───┬─ servers/ ─┬──── http/main.go      [HTTP server main]
   │            │
   │            ├──── udp/main.go       [UDP server main]
   │            │
   │            ├──── tcp/main.go       [TCP server main]
   │            │
   │            ├──── connector/main.go [Connector server main]
   │            │
   │            └──── ws/main.go        [WebSocket server main]
   │
   ├─ mars/ ───────── main.go
   │
   │
   ├─ healthcheck/ ── main.go
   │
   │
   ├─ configs/ ─┬──── shared/ [Shared configuration directory] ────────────────────┬─ field.json
   │            │                                                                  ├─ group.json
   │            ├──── http/     [HTTP configuration directory] ──────── main.json  ├─ log.json
   │            │                                                                  ├─ matching.json
   │            ├──── udp/      [UDP configuration directory]  ──────── main.json  └─ mesh.json
   │            │
   │            ├──── tcp/      [TCP configuration directory]  ──────── main.json
   │            │
   │            ├──── connector [Connector configuration directory] ─── main.json
   │            │
   │            └──── ws/       [WebSocket configuration directory] ─── main.json
   │
   ├─ cmds/  [Custom client command directory] ─┬── main.go [Entry point for all cmds]
   │                                            ├── http   ──────────────────────────────────────┬─── main.go
   ├─ lib/   [Shared library directory]         ├── room   ──────────────────────────── main.go  └─── matching.go
   │                                            ├── group  ──────────────────────────── main.go
   ├─ bin/   [Built server binary directory]    ├── field  ──────────────────────────── main.go
   │                                            │
   └─ go.mod [Go module file for the project]   └── custom ──────────────────────────── main.go
```

# Server Entry Points

## HTTP server

```
servers/http/main.go
```

## UDP server

```
servers/udp/main.go
```

## TCP server

```
servers/tcp/main.go
```

## WebSocket server

```
servers/ws/main.go
```

## Connector server

```
servers/connector/main.go
```

# MARS

Diarkis server cluster requires it's unique server called MARS server.

You need to simply build MARS server and deploy it along with other Diarkis servers.

## Starting MARS Server With Configuration File

By default, there is MARS configuration JSON file: `configs/mars/main.json`

```
./mars <path to the config JSON file>
```

# Commands

This is where you add your custom commands.

## UDP and TCP

```
cmds/
```

## WebSocket

```
ws_cmds/
```

# Configurations

This is where you place your configuration JSON files.

```
configs/
```

# MatchMaker

This is where you define your own MatchMaker rules.

The template provides HTTP API endpoints, but you may implement UDP, TCP, WebSocket commands for MatchMaker as well.

```
cmds/http/matching.go
```

## Default HTTP API Endpoints

We have pre-define rules of `rank` and `score`.

`rank` and `score` are also the `mmID` of each rule.

### Adding data to MatchMaker pool to be searched

```
POST /mm/add/:mmID/:uniqueID/:ttl
```

- `mmID` is the ID of MatchMaker rule.

- `uniqueID` is the unique ID of the data that is to be added to MatchMaker pool.

- `ttl` is the TTL of the data that is to be added to MatcMaker pool.

#### Request Body

- `props` is the JSON data to represent conditions of the data. This is used by search.

- `metadata` is the JSON data to be returned when search returns with results.

### Searching candidates from MatchMaker pool

```
POST /mm/search/:mmIDs/:limit
```

- `mmIDs` is a comma separated MatchMaker rule IDs.

- `limit` defines how many matching results you expect.

#### Request Body

- `props` is the JSON data representing search conditions.

# Custom Commands

This is where you impleement your own custom commands for TCP, UDP/RUDP.

```
/cmds/custom/main.go
```

This is where you implement your own custom commands for WebSocket.

```
ws_cmds/custom/main.go
```
