# Overview

This project is comprised of **(4)** servers and **(1)** client.

- **MARS** server is a standard (out-of-the-box) Diarkis template. It serves to orchestrates the
  node mesh.

- **HTTP** server is a standard (out-of-the-box) Diarkis template; excepting a simple Matchmaker
  profile definition for our custom matchmaking criteria. Our Matchmaker candidate information is
  stored here. It additionally handles Diarkis user authentication.

- **UDP** server hosts the client connection and handles all incoming Matchmaker commands.
  It queries the Matchmaker storage server (**HTTP**) for valid candidates. If a valid
  matching is found via the provided pooling constraints it attempts match the selected candidates.

The goal of this sample is to demonstrate how to implement secure authentication and matchmaking for
Diarkis. We will demonstrate the following security features:

- Authentication performed via an external API server.
- Safe retrieval of user rank data from the API server during ticket issuance.
- Prevention of client-side tampering of authentication data and user rank data.

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

This will create the **MARS**, **HTTP**, **UDP** server binaries, and client binary,
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
behavior using the team-based matchmaking example.

### Use Example

To test the team-based matchmaking scenario, execute the following on **(4)** separate instances of
the provided client:

```sh
./remote_bin/cli -uid user-1
```

```sh
./remote_bin/cli -uid user-2
```

```sh
./remote_bin/cli -uid user-3
```

```sh
./remote_bin/cli -uid user-4
```