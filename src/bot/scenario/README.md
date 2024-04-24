# Scenario bot Documentation

This is a template to develop a scenario test where you can test Diarkis by imitating an actual operation after integrating to an application.
It is also used for load test.

## How to Build

```
go build -o remote_bin/bot ./bot/scenario
```

## How to Run

### Run Locally

#### Usage

```
./remote_bin/bot --type "Ticket" --run "OneTicket"
```

#### Options

- `type`: determines which scenario to run. It is set in `ScenarioFactoryList` in `bot/scenario/scenarios/main.go`
- run: selects with which parameter it runs the scenario. It is the top level key in a json file in `bot/scenario/config/`
- `howmany` (optional): sets with how many clients you want to run the scenario.
- `duration` (optional): sets how long you want to run the scenario.
- `interval` (optional): sets how long the bot should wait for after in between each spawn of bot clients.

> [!NOTE]  
> All optional fields can be also set in config file.  
> If you specify as bot parameter, it will overwrite the value set in the config file.

> [!NOTE]  
> All Json files under `bot/scenario/config/` will be read by running bot.  
> It is a best practice to have root level keys in `global.json` whereas all root level keys set in different files are also assumed as global config.

### Run on Server

#### Usage

```
DIARKIS_BOT_SERVER_MODE=true ./remote_bin/bot
```

Once you run the bot, it will listen the HTTP request to execute a scenario.  
The default port is 9500 but you can change it with an environment variable.

#### Get Prometheus metrics

```
curl http://x.x.x.x:9500/metrics/
```

#### Environment variables

- `DIARKIS_BOT_SERVER_MODE`: Set "true" to run as server mode. (default: false)
- `DIARKIS_BOT_CONFIG`: Set path to the directory for config file. (default: ./bot/scenario/config)
- `DIARKIS_BOT_ADDRESS`: Set an address for the bot to listen. (default: localhost)
- `DIARKIS_BOT_PORT`: Set a port for the bot to listen. (default: 9500)

### Run on Kubernetes

#### Create and Push an image

```
REPOSITORY_NAME=__YOUR_REPOSITORY__
docker build --platform=linux/amd64 -f docker/bot/Dockerfile remote_bin -t ${REPOSITORY_NAME}/bot:dev0
docker push ${REPOSITORY_NAME}/bot:dev0
```

#### Deploy on a Cluster

We use kustomize to generate a config map from json files.

```
kustomize build ./bot/scenario/ | kubectl apply -f -
```

soon it will create a service with external IP Address.

```
NAME                TYPE           CLUSTER-IP   EXTERNAL-IP     PORT(S)        AGE
service/bot         LoadBalancer   y.y.y.y      x.x.x.x         80:3xxxx/TCP   57s
```

Then you can issue the same query as [Run on Server](#run-on-server)

```
curl -X POST http://x.x.x.x:9500/run/ -d '{"type":"Ticket","run":"OneTicket","howmany":10,"duration":10}'
```

## How to get Result

### CSV File

It will generate a result csv file under `/tmp/` after the bot run.
Here's an example output.

```
ticketType,ver3-cmd100-Push,ISSUE_TICKET-TYPE0,MATCHING_DURATION-average,active-users
0,21,13,1.3549708437500003,4
```

### Prometheus metrics

If you run the bot in server mode, you can issue the query below to get the Prometheus metrics.

```
curl -X POST http://x.x.x.x:9500/run/ -d '{"type":"Ticket","run":"OneTicket","howmany":10,"duration":10}'
```

## Development Guide

1. Implement scenario.
   1. Copy `bot/scenario/scenarios/ticket.go` and rename it to describe your own scenario.
   2. It should implement an interface named `Scenario` in `scenarios` package.
   3. Create a function to new a implemented struct and set it in `ScenarioFactoryList` in `bot/scenario/scenarios/main.go`.
2. Create config file.
   1. Copy `bot/scenario/config/ticket.json` and rename it.
   2. Change the root key (ex. "OneTicket") into some descriptive thing for your config set.

### Tips

- UDP client in `bot_client` package wraps `Send`, `RSend`, `OnResponse` and `OnPush`. It contains metrics collector so it's better to use it rather than use the one in Diarkis core directly.
  - If you want to use the one in the core, it's recommended to call `report.IncrementResponseMetrics` at the same time.
- Use `(*CustomMetrics).Add()` and `(*CustomMetrics).Increment()` to collect values for a statistics. It will be available as a metrics [How to Get Result](#how-to-get-result).
