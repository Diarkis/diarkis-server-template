This project is composed of three servers and one client.
MARS and http server are standard process with only the HTTP server
defining a simple matching profile.
The third server is the udp one that handle the matching command of the client.
The goal of this sample is to explain how to add custom criteria to matching
that cannot be set as a matching property.
Here we use the client geolocation as an example.

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
  -latitude float
        the client's latitude (default 35.65877910898138)
  -longitude float
        the client's longitude (default 139.70128360250285)
  -uid string
        the unique identifier of the client like user ID
```

The UDP server will prevent client to match together if the distance between
them is more than 100.
Note that this sample does not rely on any geolocalization DB nor service.
It is up to you to plug such external service into the code.
