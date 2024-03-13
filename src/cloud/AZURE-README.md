# Diarkis in Azure
## Overview
Azure で Diarkis を動作させるためのインフラの構築手順及び、Kubernetes manifest を格納しております。

## インフラ構築手順
下記の手順でインフラの構築手順を述べていきます。

1. AKE クラスタの作成
2. AKE への接続
3. Diarkis CLIを使ってイメージをビルド
4. Container イメージを作成
5. AKS Network Security Group の設定
6. GKE クラスタへのマニフェスト反映
7. Diarkis 起動確認

## インフラ構築手順1 - AKE クラスタの作成
Diarkis を動作させる AKE クラスタを構築します。
Diarkis を動作させるのにはそれぞれの Node が PublicIP を持つように構築する必要があります。
az deployment group create --resource-group ${RESOURCE_GROUP} --template-file k8s-demo.json

## インフラ構築手順2 - AKE クラスタへの接続
kubectl に認証を通します。
az account set --subscription ${SUBSCRIPTION_ID}
az aks get-credentials --resource-group ${RESOURCE_GROUP} --name diarkis-sample

## インフラ構築手順3 - Diarkis CLIを使ってイメージをビルド
まずlinux用にクロスコンパイルを行う準備をする。
```
cp build.yml build.linux.yml
```

とし。build.linux.yml内でGOOSとGOARCHを下記のように設定する。
```
  env:
    GOOS: linux
    GOARCH: amd64
    CGO_ENABLED: 0
```

server-templateから生成した project の root から下記を実行します。
```
./diarkis-cli/bin/diarkis-cli build --host builder.diarkis.io -c build.linux.yml
```
remote_bin にサーバーの実行ファイル郡が生成されます。

## インフラ構築手順4 - Container イメージを作成
az acr create -n ${REPOSITORY_NAME} -g ${RESOURCE_GROUP} --sku Standard

az acr login --name ${REPOSITORY_NAME}

生成したプロジェクトで、dockerイメージを作成します。
```
docker build --platform=linux/amd64 -t ${REPOSITORY_NAME}.azurecr.io/udp ./remote_bin -f docker/udp/Dockerfile
docker build --platform=linux/amd64 -t ${REPOSITORY_NAME}.azurecr.io/tcp ./remote_bin -f docker/tcp/Dockerfile
docker build --platform=linux/amd64 -t ${REPOSITORY_NAME}.azurecr.io/http ./remote_bin -f docker/http/Dockerfile
docker build --platform=linux/amd64 -t ${REPOSITORY_NAME}.azurecr.io/mars ./remote_bin -f docker/mars/Dockerfile
```
image を作成したgcrにpushします。
```
docker push ${REPOSITORY_NAME}.azurecr.io/udp
docker push ${REPOSITORY_NAME}.azurecr.io/tcp
docker push ${REPOSITORY_NAME}.azurecr.io/http
docker push ${REPOSITORY_NAME}.azurecr.io/mars
```

## インフラ構築手順5 - AKS Network Security Group の設定
対象の Network Security Group に下記のルールを追加します。
```
Rule: Ingress
Protocol: TCP & UDP
Port: 7000 - 8000
Source: *
```

## インフラ構築手順6 - GKE クラスタへのマニフェスト反映
kustomize を使用し、GKE クラスタにマニフェストの反映をします
```
kustomize build overlays/dev0/ | sed -e "s/__REPOSITORY_NAME__/${REPOSITORY_NAME}/g" | kubectl apply -f -
```

## インフラ構築手順7 - Diarkis 起動確認
作成したクラスタの稼働を確認するために、Diarkis の認証エンドポイントに HTTP リクエストを送信し、動作の確認をします。
下記の結果がレスポンスされれば正常に起動が完了しております。
(下記の例では xxx や yyy 等で伏せ字をしております)
```
# http loadbalancer のエンドポイント取得
$ kubectl get service http -n dev0

$ curl {http service の EXTERNAL-IP}:7000/auth/1 # 本来は、"GET /auth/{user id}"でリクエストいただく想定です。今回は動作確認のため"/auth/1"で挙動確認しております。
{"TCP":"xxxx.bc.googleusercontent.com:7200","UDP":"yyyy.bc.googleusercontent.com:7100","sid":"xxxxx","encryptionKey":"xxxxx","encryptionIV":"xxxxx","encryptionMacKey":"xxxxx"}
```

