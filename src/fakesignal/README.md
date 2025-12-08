# Overview

Diarkis by default uses UNIX signals **USR1** and **USR2** to manage its debugging state.

On Windows, to emulate the signal you can use `fakesignal.exe`.

To send **USR1** signal to one of Diarkis process, run the following command.

```batch
win-amd64/fakesignal.exe <PID> SIGUSR1 # on amd64 arch
win-arm64/fakesignal.exe <PID> SIGUSR1 # on arm64 arch
```
