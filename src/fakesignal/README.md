# fakesignal

Unix signal (SIGUSR1, SIGUSR2, SIGHUP) emulation utility for Windows.

## Overview

Since Windows does not have Unix signals like SIGUSR1 or SIGUSR2, this utility emulates these signals using Windows Event API. It is used to send signals to Diarkis server processes.

## Features

- **Supported signals**: SIGUSR1, SIGUSR2, SIGHUP
- Signal emulation using Windows Event API
- Send signals by specifying process ID (PID)

## Usage

```bash
fakesignal <PID> <SIGNAL>
```

### Examples

```bash
fakesignal 1234 SIGUSR1
fakesignal 5678 SIGHUP
```

## Implementation Details

- `main.go` - CLI tool for sending signals
- `signal.go` - Signal handling mechanism (used in-process)
- Windows Event name format: `DIARKIS#SIGNAL_<SIGNAME>#PID_<PID>`

## Build Requirements

`//go:build windows` - Windows platform only
