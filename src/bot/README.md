# Overview
This directory contains bots.

- room bot ( random join rooms and broadcast messages )
- resonance bot ( connect diarkis and send custom commands packet and receive packet )
- matchmaker bot ( keep matchmaking bot )

# room bot
```
# example
cd room
go build 
./room localhost:7000 5 10 100 # this means 5 bot clients send 10 byte packet per 100ms.
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
