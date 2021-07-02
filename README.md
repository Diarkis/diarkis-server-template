# TemplateOf Server Code

Server code template aims to provide entry point to server development using Diarkis.

# How To Install Template

```
make install project={Project name of your application} output={Install destination path}
```

# Structure

```
───┬─ servers/ ─┬─ http/main.go      [HTTP server main]
   │            │
   │            ├─ udp/main.go       [UDP server main]
   │            │
   │            ├─ tcp/main.go       [TCP server main]
   │            │
   │            ├─ connector/main.go [Connector server main]
   │            │
   │            └─ ws/main.go        [WebSocket server main]
   │
   ├─ configs/ ─┬─ shared/ [Shared configuration directory] ────────────────────┬─ field.json
   │            │                                                               ├─ group.json
   │            ├─ http/     [HTTP configuration directory] ──────── main.json  ├─ log.json
   │            │                                                               ├─ matching.json
   │            ├─ udp/      [UDP configuration directory]  ──────── main.json  └─ mesh.json
   │            │
   │            ├─ tcp/      [TCP configuration directory]  ──────── main.json
   │            │
   │            ├─ connector [Connector configuration directory] ─── main.json
   │            │
   │            └─ ws/       [WebSocket configuration directory] ─── main.json
   │
   ├─ cmds/  [Custom client command directory] ─┬─ main.go [Entry point for all cmds]
   │                                            ├── http   ──────────────────────────────────────┬─── main.go
   ├─ lib/   [Shared library directory]         ├── room   ──────────────────────────── main.go  └─── matching.go
   │                                            ├── group  ──────────────────────────── main.go
   ├─ bin/   [Built server binary directory]    ├── field  ──────────────────────────── main.go
   │                                            │
   └─ go.mod [Go module file for the project]   └── custom ──────────────────────────── main.go

```
