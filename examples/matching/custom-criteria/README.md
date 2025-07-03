# Overview

This project is comprised of **(3)** servers and **(1)** client.

- **MARS** server is standard (out-of-the-box) Diarkis template. It serves to orchestrates the
  node mesh.
- **HTTP** servers is standard (out-of-the-box) Diarkis template, with only a simple matching
  profile for our custom matching criteria. Our matching information is stored on this server type.

- **UDP** server hosts the client connection and handles all incoming Matchmaker commands.
  It queries the matchmaking storage server (**HTTP**) for valid candidates, and if a valid
  pairing exists, it completes said matching for the selected candidates.

The goal of this sample is to demonstrate how to add custom criteria to Diarkis Matchmaker.
These are criteria that may not be represented as an `int` constraint within the map of defined `AddProperty` nor `SearchProperty`.

In this example, we use geolocation-based candidate matchmaking to demonstrate one such implementation.

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

This will create the **MARS**, **HTTP**, and **UDP** server binaries, along with the client binary, placing them inside the `remote_bin` directory.

## How to Run

This project requires all **(3)** servers to be running before clients can test matchmaking behavior.

### 1. First, start the MARS server to orchestrate the node mesh

```sh
./run-mage.sh server mars
```

### 2. Next, start the HTTP server, which holds the custom matching criteria

```sh
./run-mage.sh server http
```

### 3. Then, start the UDP server, to manage incoming client connections

```sh
./run-mage.sh server udp
```

Once all servers are running, you may start two client instances to test and observe the matchmaking behavior using the custom criteria flow.

## Testing the Custom Matchmaking Criteria

**Scenario**: The **UDP** server will prevent candidate to matchmaking if the distance between them is calculated to be >5000 km.

The test client has the following usage.

```output
Usage of ./remote_bin/cli:
  -clientKey string
        the client key to authenticate with the server
  -host string
        the address of the HTTP server (default "127.0.0.1:7000")
  -latitude float
        the client's latitude (default 35.65877910898138)
  -longitude float
        the client's longitude (default 139.70128360250285)
  -uid string
        the unique identifier of the client like user ID
```

**NOTE**: This sample does not implment any geolocalization database nor service. If desired, it is up to you to plug such an external service into the example code.

## Examples

The custom criteria outlined in this example uses client-provided decimal degrees (DD) coordinates as the
matching conditions.

### Valid Matching

Our first user, `user-a`, is connecting from **New York City, USA** (40.7128, -74.0060).

```sh
 ./remote_bin/cli -uid user-a -latitude 40.7128 -longitude -74.0060
```

Our second user, `user-b`, is connecting from **Los Angeles, USA** (34.052235, -118.243683).

```sh
 ./remote_bin/cli -uid user-b -latitude 34.052235 -longitude -118.243683
 ```

The distance between NYC and LA is **~3987 km**, and therefore is less than our **<5000 km** distance
matchmaking constraint. This means that our two clients will be able to match together successfully.

### Invalid Matching

Our first user, `user-a`, is connecting from **New York City, USA** (40.7128, -74.0060).

```sh
 ./remote_bin/cli -uid user-a -latitude 40.7128 -longitude -74.0060
```

Our second user, `user-b`, is connecting from **Tokyo, JP** (35.652832, 139.839478).

```sh
 ./remote_bin/cli -uid user-b -latitude 35.652832 -longitude 139.839478
 ```

The distance between NYC and Tokyo is **~18744 km**, and therefore is less than our **<5000 km** distance
matchmaking constraint. This means that our two clients will fail to match together.
