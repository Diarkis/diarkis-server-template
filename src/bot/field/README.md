# Field Sample Bot

Program in Go which creates a discretionary number of clients connected to a Diarkis field server and make them move randomly on the grid.

# Usage

## Dependencies

The source code of the Diarkis Client is needed in order to build the bot.
You can set the path to your own Diarkis folder by editing the last line of `go.mod` with:

```
replace github.com/Diarkis/diarkis => /path/to/my/diarkis
```

## Build

In this folder simply run:

```
go build
```

## Run

The application take 4 arguments:

```
./diarkis-field-bot $HOST_ENDPOINT $BOTS_NUM udp $PACKET_SEND_INTERVAL_MILLISECOND
```

Example, for starting 100 bots updating every 2 seconds:

```
./diarkis-field-bot 127.0.0.1:7000 100 udp 2000
```

# About

## Random walk

The movements behavior is defined in the `randomSync` function in `main.go`.

Depending on the number of UDP nodes (server side) the Grid position may change, and subsequently the bot position too.\
The number of server and the map size (or field/grid size) is hardcoded and set to 1 by default.\
You're free to change those values according to your server environment with:

```
const SERVER_COUNT = 1
const MAP_SIZE = 4500
```

The bot may move or not after every intervals, this probability is defined between 0 and 100.\
It be modified with:

```
var moveRatio = 25
```

The covered distance at each moves is always the same, and can be modified with:

```
const BOT_MOVE_RANGE = 500
```
