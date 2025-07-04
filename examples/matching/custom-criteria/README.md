# Overview

This project is comprised of **(3)** servers and **(1)** client.

- **MARS** server is a standard (out-of-the-box) Diarkis template. It serves to orchestrates the
  node mesh.

- **HTTP** server is a standard (out-of-the-box) Diarkis template; excepting a simple Matchmaker
  profile definition for our custom matchmaking criteria. Our Matchmaker candidate information is
  stored here.

- **UDP** server hosts the client connection and handles all incoming Matchmaker commands.
  It queries the Matchmaker storage server (**HTTP**) for valid candidates. If a valid
  matching is found via the provided pooling constraints _(**see**: `matching.AddProperty`,
  `matching.SearchProperty`)_, it attempts match the selected candidates by the (optional) custom
  matchmaking criteria.

The goal of this sample is to demonstrate how to add custom matchmaking criteria to Diarkis
Matchmaker. Custom matchmaking criteria are constraints that may not be represented as an `int`
within the map of defined `matching.AddProperty` or `matching.SearchProperty` utilized for candidate
pooling.

In this example, we use geolocation-based candidate matchmaking to demonstrate one such
implementation.

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
behavior using the custom criteria constraint.

## Testing the Custom Matchmaking Criteria

**Scenario**: The **UDP** server will prevent matchmaking if the distance between any (2) candidates
is calculated to be **>5000 km**.

The test client has the following usage:

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

**NOTE**: This sample does not implment a geolocalization database nor service. If desired, it is
up to you to plug-in such an external service into the code example.

## Examples

The custom matchmaking criteria outlined in this example uses client-provided decimal degree (DD)
coordinates as the matchmaking conditions.

### Valid Matching

Our first user, `user-a`, is connecting from **New York City, USA** (40.7128, -74.0060).

```sh
 ./remote_bin/cli -uid user-a -latitude 40.7128 -longitude -74.0060
```

Our second user, `user-b`, is connecting from **Los Angeles, USA** (34.052235, -118.243683).

```sh
 ./remote_bin/cli -uid user-b -latitude 34.052235 -longitude -118.243683
```

The distance between NYC and LA is **~3987 km**, and therefore is less than our **<5000 km**
distance matchmaking constraint. This means that our two clients are able to match together
successfully.

**Output:**

```output
[2025/07/03 06:39:11.996]<SERVER>      DEBUG    UDP|127.0.0.1:8101      Candidates matched successfully Owner=user-a Candidate=user-b Distance=3986.9854666330843 MaxAllowedDistance=5000
```

### Invalid Matching

Our first user, `user-a`, is connecting from **New York City, USA** (40.7128, -74.0060).

```sh
 ./remote_bin/cli -uid user-a -latitude 40.7128 -longitude -74.0060
```

Our second user, `user-b`, is connecting from **Tokyo, JP** (35.652832, 139.839478).

```sh
 ./remote_bin/cli -uid user-b -latitude 35.652832 -longitude 139.839478
```

The distance between NYC and Tokyo is **~18744 km**, and therefore is less than our **<5000 km**
distance matchmaking constraint. This means that our two clients will fail to match together.

**Output:**

```output
[2025/07/03 06:39:11.996]<SERVER>      DEBUG    UDP|127.0.0.1:8101      Cannot match candidates Owner=user-a Candidate=user-b Distance=18744.26784351792 MaxAllowedDistance=5000
```

---

_First created on 2025-04-03. Updated on 2025-04-04._
