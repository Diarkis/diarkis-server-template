# Overview

This directory contains bots for load tests.

- room bot ( random join rooms and broadcast messages )
- resonance bot ( connect diarkis and send custom commands packet and receive packet )
- matchmaker bot ( keep matchmaking bot )

# room bot

```
# example
cd room
go build
./room localhost:7000 5 10 100 # this means 5 bot clients send 10 byte packet to room per 100ms.
```

# group bot

```
# example
cd group
go build
./group localhost:7000 5 10 100 # this means 5 bot clients send 10 byte packet to group per 100ms.
```

# resonance bot

```
# example
cd resonance
go build
./resonance localhost:7000 5 100 800 # this meands 5 bot clients send 800byte resonance command every 100ms
```

# matchmaker bot

```
# example
cd matchmaker
go build
./matchmaker localhost:7000 10 30 1000 # This means that 30% of the 10 bots will be hosts and the other 50% will be guests, searching every 800ms
```

# dm bot

```
# example
cd dm
go build
./dm host=$(HTTP host) protocol=$(UDP or TCP) bots=$(how many bots) size=$(message size) interval=$(message send interval in milliseconds)
```

# field bot

```
# example
cd field
./field %(http host) $(how many bots) $(protocol: UDP or TCP) $(update interval in milliseconds) $(map size) $(movement range)
```
