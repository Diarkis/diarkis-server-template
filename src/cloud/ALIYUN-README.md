# Diarkis in GCP

## Overview

alibaba cloud(aliyun)　で　Diarkis　を動作させるためのインフラの構築手順及び、Kubernetes manifest を記述しております。

## インフラ構築手順

下記の手順でインフラの構築手順を述べていきます。

1. VPCの作成
2. ACK(Alibaba Cloud Container Service for Kubernetes)で Kubernetes クラスタを作成
3. firewall の作成
4. ACKへの接続
5. manifest の修正
6. manifest の apply

## インフラ構築手順0 - 事前準備

alibaba cloud への権限のあるアカウントで alibaba cloud console へログインしてください。
途中一般ユーザーだと規約に同意できない箇所があるので初回は Admin ユーザーアカウントでログインをしていただければと思います。

## インフラ構築手順1 - VPC の作成

VPC 作成画面に移動し、create VPC ボタンから VPC を作成する。
![VPC-0](img/vpc-0.png)
下記のように、VPCの名前の設定、IP CIDR、vStwitchの作成等を行います。
画像の例では 10.0.0.0/8 で作成していますが、他のVPCとVPNでつなぐ必要があったり、組織で使用CIDRを管理している場合はそれに習ってください。
また vSwitch についても、リリース予定のサービス規模などに応じて設定してください。
![VPC-1](img/vpc-1.png)

## インフラ構築手順2 - ACK による kuberentes cluster の作成

## インフラ構築手順3 - firewall の作成

Diarkis は、外部と通信する際、

```
TCP: 7000-8000
UDP: 7000-8000
```

のポートレンジを使用しますのでこの範囲のアクセスを ACK の firewall によって許可する必要があります。

```
export PROJECT_NAME=YOUR_AWESOME_PROJECT_NAME # 適宜修正してください
gcloud compute --project=$PROJECT_NAME firewall-rules create diarkis-ingress-allow --direction=INGRESS --priority=1000 --network=default --action=ALLOW --rules=tcp:7000-8000,udp:7000-8000 --source-ranges=0.0.0.0/0 --target-tags=diarkis
```

network tag は diarkis という名前で設定していきます

## インフラ構築手順2 - GCR の有効化及び設定

GCR は Diarki のコンテナイメージを配置する場所として今回使用しますので、有効化と Docker の設定を行います。
またCloudDNSはkubernetesクラスターで使用するので、有効化します。
まず https://console.cloud.google.com/flows/enableapi?apiid=containerregistry.googleapis.com にアクセスして対象プロジェクトの GCR API を有効化します。
次に、https://console.cloud.google.com/flows/enableapi?apiid=dns にアクセスしてCloudDNSを有効化します。

```
gcloud auth configure-docker # docker が gcloud を使って認証するように設定する
```

## インフラ構築手順3 - GKE クラスタの作成

Diarkis を動作させる GKE クラスタを構築します。
Diarkis を動作させるのにはそれぞれの Node が PublicIP を持つように構築する必要があります。
また、手順１で作成した firewall と結びつけるために diarkis tag を付与しています。

```
export CLUSTER_NAME=YOUR_AWESOME_CLUSTER_NAME
gcloud beta container --project $PEOJECT_NAME clusters create $CLUSTER_NAME --zone "asia-northeast1-a" --no-enable-basic-auth --cluster-version "1.19.9-gke.1900" --release-channel "regular" --machine-type "c2-standard-4" --image-type "COS_CONTAINERD" --disk-type "pd-standard" --disk-size "50" --metadata disable-legacy-endpoints=true --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" --max-pods-per-node "110" --num-nodes "3" --enable-stackdriver-kubernetes --enable-ip-alias --network "projects/gcgs-poc/global/networks/default" --subnetwork "projects/gcgs-poc/regions/asia-northeast1/subnetworks/default" --no-enable-intra-node-visibility --default-max-pods-per-node "110" --enable-autoscaling --min-nodes "0" --max-nodes "3" --enable-dataplane-v2 --no-enable-master-authorized-networks --addons HorizontalPodAutoscaling,HttpLoadBalancing,NodeLocalDNS,GcePersistentDiskCsiDriver --enable-autoupgrade --enable-autorepair --max-surge-upgrade 1 --max-unavailable-upgrade 0 --enable-shielded-nodes --tags "diarkis" --node-locations "asia-northeast1-a" --cluster-dns clouddns --cluster-dns-scope vpc --cluster-dns-domain diarkis.cluster
```

--cluster-version や、--machine-type などに関してはプロジェクトによって適宜調整していただければと思います。

## インフラ構築手順4 - GKE クラスタへの接続

kubectlに認証を通します。

```
gcloud container clusters get-credentials $CLUSTER_NAME --project $PROJECT_NAME
```

## インフラ構築手順5 - GKE クラスタへのマニフェスト反映

kustomize を使用し、GKE クラスタにマニフェストの反映をします

```
kustomize build overlays/dev0/ | sed -e "s/__GCP_PROJECT_ID__/${PROJECT_NAME}/g" | kubectl apply -f -
```

## インフラ構築手順6 - Diarkis 起動確認

作成したクラスタの稼働を確認するために、Diarkis の認証エンドポイントに HTTP リクエストを送信し、動作の確認をします。
下記の結果がレスポンスされれば正常に起動が完了しております。
(下記の例では xxx や yyy 等で伏せ字をしております)

```
# http loadbalancer のエンドポイント取得
$ kubectl get service http -n dev0

$ curl {http service の EXTERNAL-IP}:7000/auth/1 # 本来は、"GET /auth/{user id}"でリクエストいただく想定です。今回は動作確認のため"/auth/1"で挙動確認しております。
{"TCP":"xxxx.bc.googleusercontent.com:7200","UDP":"yyyy.bc.googleusercontent.com:7100","sid":"xxxxx","encryptionKey":"xxxxx","encryptionIV":"xxxxx","encryptionMacKey":"xxxxx"}
```
