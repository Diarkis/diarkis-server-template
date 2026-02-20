# linode での構築

## クラスタ作成

https://cloud.linode.com/kubernetes/clusters
diarkis-load-test とかで作っとく

## ローカルの kubectl で接続

config をダウンロードして ~/.kube/config にマージ

```bash
vim ~/.kube/config
kubectl-ctx
```

## KUBECONFIG でファイル変更する方法もあります

```
export KUBECONFIG=/tmp/diarkis-load-test-kubeconfig.yaml
```

## diarkis-cli でビルド

./diarkis-cli/bin/diarkis-cli build --host builder.diarkis.io -c build.linux.yml

## ローカルビルド

GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o remote_bin/mars ./mars/main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o remote_bin/udp ./servers/udp/main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o remote_bin/tcp ./servers/tcp/main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o remote_bin/http ./servers/http/main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o remote_bin/health-check ./healthcheck/main.go
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o remote_bin/ms ./mars-stats/main.go

# docker build

DOCKER_REPOSITORY=kojinakahatadiarkis/diarkis
docker login

docker build --platform linux/amd64 -t ${DOCKER_REPOSITORY}:udp ./remote_bin -f docker/udp/Dockerfile
docker build --platform linux/amd64 -t ${DOCKER_REPOSITORY}:tcp ./remote_bin -f docker/tcp/Dockerfile
docker build --platform linux/amd64 -t ${DOCKER_REPOSITORY}:http ./remote_bin -f docker/http/Dockerfile
docker build --platform linux/amd64 -t ${DOCKER_REPOSITORY}:mars ./remote_bin -f docker/mars/Dockerfile

# docker push

docker push ${DOCKER_REPOSITORY}:udp
docker push ${DOCKER_REPOSITORY}:tcp
docker push ${DOCKER_REPOSITORY}:http
docker push ${DOCKER_REPOSITORY}:mars

# 以下で上記の一連の作業が一発でできる

make build-linode

kustomize build k8s/linode/overlays/dev0/ | kubectl apply -f -

## metrics API を有効にする

kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

kubectl edit deployments.apps -n kube-system metrics-server

args に以下追加

```
      - args:
        - --kubelet-insecure-tls
```

## prometheus / grafana

### 7. Prometheus / Grafana を立てる

https://www.notion.so/kube-prometheus-c6192836218a43c085a9df8dca896c41

## bot build

```bash
cd bot/field
go build -o field-bot
# linux 向け
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o field-bot-amd64

SERVER_IP=172.233.68.22
# 100
./field-bot ${SERVER_IP} 100 udp 2000 4500 500
# 1 は old, 2 は new(puffer)
./field-bot ${SERVER_IP} 100 udp 1 2000 4500 500
./field-bot ${SERVER_IP} 100 udp 2 2000 4500 500
setsid ./field-bot ${SERVER_IP} 100 udp 2000 4500 500 > /tmp/field-bot.log
# 300
./field-bot ${SERVER_IP} 300 udp 1 2000 4500 500
./field-bot ${SERVER_IP} 300 udp 2 2000 4500 500
# 500
./field-bot ${SERVER_IP} 500 udp 1 2000 4500 500
./field-bot ${SERVER_IP} 500 udp 2 2000 4500 500
# 1000
./field-bot ${SERVER_IP} 1000 udp 1 2000 4500 500
./field-bot ${SERVER_IP} 1000 udp 2 2000 4500 500
setsid ./field-bot ${SERVER_IP} 1000 udp 2000 4500 500 > /tmp/field-bot.log
# 2000
./field-bot ${SERVER_IP} 2000 udp 2000 4500 500
# 5000
./field-bot ${SERVER_IP} 5000 udp 2000 4500 500
# client 1000 field size 16000
./field-bot ${SERVER_IP} 1000 udp 2000 16000 500
# client 1000 field size 80000
./field-bot ${SERVER_IP} 1000 udp 2000 80000 500
```

これをどっかの VM で叩いて計測する

## bot room

```bash
cd bot/room
go build -o room-bot
# linux 向け
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o room-bot-amd64

SERVER_IP=172.233.68.47
HOST=${SERVER_IP} BOTS=100 ROOM_SIZE=64 PACKET_INTERVAL=100 PACKET_SIZE=100 AUTH_INTERVAL=300 ./room-bot
HOST=${SERVER_IP} BOTS=500 ROOM_SIZE=64 PACKET_INTERVAL=100 PACKET_SIZE=100 AUTH_INTERVAL=300 ./room-bot
HOST=${SERVER_IP} BOTS=500 ROOM_SIZE=100 PACKET_INTERVAL=200 PACKET_SIZE=300 AUTH_INTERVAL=300 ./room-bot
HOST=${SERVER_IP} BOTS=1000 ROOM_SIZE=64 PACKET_INTERVAL=100 PACKET_SIZE=100 AUTH_INTERVAL=300 ./room-bot

HOST=${SERVER_IP} BOTS=100 ROOM_SIZE=100 PACKET_INTERVAL=200 PACKET_SIZE=300 AUTH_INTERVAL=300 ./room-bot

./room-bot ${SERVER_IP} 5 10 100 # this means 5 bot clients send 10 byte packet to room per 100ms.


# profiling

DIARKIS_SIGUSR1 ファイルを `Debug` という内容で保存
SIGUSR1 シグナルをプロセスに送ると、プロファイリングできるようになる

```

echo "Debug" > /tmp/DIARKIS_SIGUSR1
kill -s SIGUSR1 1
wget -O profile.cpu "http://127.0.0.1:6060/debug/pprof/profile?seconds=600"
wget -O profile.heap "http://127.0.0.1:6060/debug/pprof/heap?seconds=10"
wget -O profile.goroutine "http://127.0.0.1:6060/debug/pprof/goroutine?debug=1"

```

```
