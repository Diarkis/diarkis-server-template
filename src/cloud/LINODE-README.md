# Diarkis in LINODE

## Overview

LKE で Diarkis を動作させるためのインフラの構築手順及び、Kubernetes manifest を格納しております。

## インフラ構築手順

下記の手順でインフラの構築手順を述べていきます。

1. firewall の作成
1. コンテナレジストリの連携
1. LKE クラスタの作成
1. LKE への接続
1. Container イメージを作成
1. manifest の apply
1. Diarkis 起動確認

## インフラ構築手順1 - firewall の作成

1. Linode にログインする
1. Firewall メニューから “Create Firewall” を押下
1. 作成された Firewall を選択する

Diarkis は、外部と通信する際、

```
TCP: 7000-8000
UDP: 7000-8000
```

のポートレンジを使用しますのでこの範囲のアクセスを firewall によって許可する必要があります。  
Inbound Rules に以下を追加してください。

1. TCP 7000-8000 Accept
2. UDP 7000-8000 Accept
3. UDP 6779 Accept
4. TCP 30000-32767 Accept
   1. 上記は Node Port になり得る値をすべて許可していますが、可能な限り割り当てられたポートのみ設定するようにしてください。
   2. 2023/12 現在、Node Port の Firewall 自動設定は Linode に対する Feature Request 中
5. “Default inbound policy” Drop

また、Linode は Node を跨いだ通信も Firewall の設定が必要であるため、Diarkis の內部通信に使用する以下のポートも許可します。

1. UDP 6779 Accept
1. UDP 8000-9000 Accept

## インフラ構築手順2 - コンテナレジストリの連携

Linode は Docker イメージのリポジトリサービスが無いため、任意のレジストリサービスを準備します。  
以下の例では Dockerhub を使用しています。ユーザ情報を書き換えて以下のコマンドを実行してください。

```
kubectl create secret docker-registry secrets-dockerhub --docker-server=https://index.docker.io/v1/ --docker-username=YOUR_USER_NAME --docker-password=YOUR_PASSWORD --docker-email=YOUR_EMAIL_ADDRESS --dry-run=client -o yaml > dockerhub-secret.yaml
```

以下のファイルを生成されたファイルで上書きしてください。

```
k8s/linode/overlays/dev0/shared/secrets/dockerhub-secrets.yaml
```

また、以下のファイルの `__YOUR_DOCKERHUB_REPOSITORY_NAME__` をリポジトリ名で上書きしてください。

```
k8s/linode/overlays/dev0/kustomization.yaml
```

## インフラ構築手順3 - LKE クラスタの作成

Diarkis を動作させる LKE クラスタを構築します。  
Diarkis を動作させるのにはそれぞれの Node が PublicIP を持つように構築する必要があります。
Linode では Metadata Service API が使用可能な Region を選択する必要があります。

> **Note**  
> Metadata Service API が使用可能な Region
> https://www.linode.com/docs/products/compute/compute-instances/guides/metadata/#availability

1. Linode にログインする。
1. Kubernetes のメニューから “Create Cluster” を押下
1. クラスタ名・Region・CPU タイプを選択して “Create Cluster”
   1. Metadata Service API が使用可能な Region を選択してください。
1. Node が立ち上がる
1. Firewall のメニューから手順１で作成した Firewall を選択して、Linodes タブから作成された Node の Linode を選択する。
   1. **注意**：現状、Firewall の自動紐づけ機能は無いため、スケールによって新たに作成された Node への Firewall 紐づけは手動でしなければならない。

## インフラ構築手順4 - LKE クラスタへの接続

1. Kubernetes メニューから作成されたクラスタを選択して、Kubeconfig をダウンロードする。
2. `kubectl` コマンド発行時に `--kubeconfig` でダウンロードファイルを指定するか、ローカルの kubeconfig (`~/.kube/config`) にダウンロードしたファイルをマージする。

## インフラ構築手順5 - Diarkis CLIを使ってイメージをビルド

まずlinux用にクロスコンパイルを行う準備をする。

```
cp build.yml build.linux.yml
```

とし、build.linux.yml を下記のように設定する。

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

## インフラ構築手順6 - Container イメージを作成

生成したプロジェクトで、dockerイメージを作成します。

```
export DOCKER_REPOSITORY=__YOUR_DOCKERHUB_REPOSITORY_NAME__
docker build -t $DOCKER_REPOSITORY/udp ./remote_bin -f docker/udp/Dockerfile
docker build -t $DOCKER_REPOSITORY/tcp ./remote_bin -f docker/tcp/Dockerfile
docker build -t $DOCKER_REPOSITORY/http ./remote_bin -f docker/http/Dockerfile
docker build -t $DOCKER_REPOSITORY/mars ./remote_bin -f docker/mars/Dockerfile
```

image を作成した dockerhub に push します。

```
docker push $DOCKER_REPOSITORY/udp
docker push $DOCKER_REPOSITORY/tcp
docker push $DOCKER_REPOSITORY/http
docker push $DOCKER_REPOSITORY/mars
```

## インフラ構築手順6 - LKE クラスタへのマニフェスト反映

kustomize を使用し、LKE クラスタにマニフェストの反映をします

```
kustomize build overlays/dev0/ | kubectl apply -f -
```

## インフラ構築手順7 - Diarkis 起動確認

作成したクラスタの稼働を確認するために、Diarkis の認証エンドポイントに HTTP リクエストを送信し、動作の確認をします。
下記の結果がレスポンスされれば正常に起動が完了しております。
(下記の例では xxx や yyy 等で伏せ字をしております)

```
# http loadbalancer のエンドポイント取得
$ kubectl get svc -n dev0

$ curl {http service の EXTERNAL-IP}:7000/auth/1 # 本来は、"GET /auth/{user id}"でリクエストいただく想定です。今回は動作確認のため"/auth/1"で挙動確認しております。
{"TCP":"x.x.x.x:7200","UDP":"x.x.x.x:7100","sid":"xxxxx","encryptionKey":"xxxxx","encryptionIV":"xxxxx","encryptionMacKey":"xxxxx"}
```
