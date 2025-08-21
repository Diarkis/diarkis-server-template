This directory contains a sample http server exposing a unique endpoint
using JSON as input/output.

## Build Servers for Local Use

The following command will be building the servers for your local machine.

```sh
./run-mage.sh build:local
```

or on Windows

```bat
.\run-mage.bat build:local
```

## Start the servers

First of all you need to start diarkis mars server.

```sh
./run-mage.sh server mars
```

Then you can start the http server.

```sh
./run-mage.sh server http
```

## Test the endpoint

In the http server log you should search for the following log

```
<HTTP>         INFO Server starting -> PublicEndpoint:127.0.0.1:7000
```

You can test the endpoint by using either curl or the provided test client.

```sh
go run ./test_client_http.go -endpoint http://localhost:7000/echo
```
