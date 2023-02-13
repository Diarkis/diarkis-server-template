# Diarkis in GCP
## Overview
GKE で Diarkis を動作させるためのインフラの構築手順及び、Kubernetes manifest を格納しております。

## インフラ構築手順
下記の手順でインフラの構築手順を述べていきます。

1. firewall の作成
2. GCRとCloudDNS の有効化
3. GKE クラスタの作成
4. GKE への接続
5. manifest の修正
6. manifest の apply

なお、手順については以下のバージョンで検証しました。
```
CLUSTER_VERSION="1.24.8-gke.2000"
```


## インフラ構築手順1 - firewall の作成
Diarkis は、外部と通信する際、
```
TCP: 7000-8000
UDP: 7000-8000
```
のポートレンジを使用しますのでこの範囲のアクセスを GCP firewall によって許可する必要があります。
```
export PROJECT_NAME=YOUR_AWESOME_PROJECT_NAME # 適宜修正してください
gcloud compute --project=$PROJECT_NAME firewall-rules create diarkis-ingress-allow --direction=INGRESS --priority=1000 --network=default --action=ALLOW --rules=tcp:7000-8000,udp:7000-8000 --source-ranges=0.0.0.0/0 --target-tags=diarkis
```
network tag は diarkis という名前で設定していきます

## インフラ構築手順2 - GCR の有効化及び設定
GCR は Diarkis のコンテナイメージを配置する場所として今回使用しますので、有効化と Docker の設定を行います。  
また CloudDNSはkubernetes クラスターで使用するので、有効化します。  
まず https://console.cloud.google.com/flows/enableapi?apiid=containerregistry.googleapis.com にアクセスして対象プロジェクトの GCR API を有効化します。  
次に、https://console.cloud.google.com/apis/enableflow?apiid=dns.googleapis.com にアクセスして CloudDNS を有効化します。

```
gcloud auth configure-docker # docker が gcloud を使って認証するように設定する
```

## インフラ構築手順3 - GKE クラスタの作成
Diarkis を動作させる GKE クラスタを構築します。
Diarkis を動作させるのにはそれぞれの Node が PublicIP を持つように構築する必要があります。

必要なオプション
|  option  |  値  | 詳細 |
| ---- | ---- | ---- |
| --tags | "diarkis" | 手順１で作成した firewall との紐付けに使用します。 |
| --enable-stackdriver-kubernetes |  | モニタリングに使用するため有効化します。 |
| --enable-autoscaling |  | Node を固定する必要が無ければ有効化してください。 |
| --enable-dataplane-v2 |  | パフォーマンス向上のため有効化します。 |
| --cluster-dns | clouddns | DNS の可用性向上のため有効化します。 | 
| --cluster-dns-scope | vpc | VPC の用途によって切り替えてください。 |
| --cluster-dns-domain | diarkis.cluster | 各 DNS レコードの suffix を設定します。 |

コマンド例：
```
export CLUSTER_NAME="my-diarkis-cluster"
export CLUSTER_VERSION="1.24.8-gke.2000"
gcloud beta container clusters create $CLUSTER_NAME　\
  --project $PROJECT_NAME \
  --zone "asia-northeast1-a" \
  --no-enable-basic-auth \
  --cluster-version "1.19.9-gke.1900" \
  --release-channel "regular" \
  --machine-type "c2-standard-4" \
  --image-type "COS_CONTAINERD" \
  --disk-type "pd-standard" \
  --disk-size "50" \
  --metadata disable-legacy-endpoints=true \
  --scopes "https://www.googleapis.com/auth/devstorage.read_only","https://www.googleapis.com/auth/logging.write","https://www.googleapis.com/auth/monitoring","https://www.googleapis.com/auth/servicecontrol","https://www.googleapis.com/auth/service.management.readonly","https://www.googleapis.com/auth/trace.append" \
  --max-pods-per-node "110" \
  --num-nodes "3" \
  --enable-stackdriver-kubernetes \
  --enable-ip-alias \
  --network "projects/gcgs-poc/global/networks/default" \
  --subnetwork "projects/gcgs-poc/regions/asia-northeast1/subnetworks/default" \
  --no-enable-intra-node-visibility \
  --default-max-pods-per-node "110" \
  --enable-autoscaling \
  --min-nodes "0" \
  --max-nodes "3" \
  --enable-dataplane-v2 \
  --no-enable-master-authorized-networks \
  --addons HorizontalPodAutoscaling,HttpLoadBalancing,NodeLocalDNS,GcePersistentDiskCsiDriver \
  --enable-autoupgrade \
  --enable-autorepair \
  --max-surge-upgrade 1 \
  --max-unavailable-upgrade 0 \
  --enable-shielded-nodes \
  --tags "diarkis" \
  --node-locations "asia-northeast1-a" \
  --cluster-dns clouddns \
  --cluster-dns-scope vpc \
  --cluster-dns-domain diarkis.cluster
```
--cluster-version や、--machine-type などに関してはプロジェクトによって適宜調整していただければと思います。

## インフラ構築手順4 - GKE クラスタへの接続
kubectlに認証を通します。
```
gcloud container clusters get-credentials $CLUSTER_NAME --project $PROJECT_NAME
```

## インフラ構築手順5 - サーバ用バイナリのビルド
Todo: 各クラウド共通の手順を作成してリンク

## インフラ構築手順6 - イメージの生成・レジストリへのプッシュ
docker イメージをビルドして、クラウド上のレジストリにプッシュします。

```
# タグ名を生成
IMAGE_TAG=`date "+%Y%m%d-%H%M%S"`

# docker imageのビルド
docker build --platform=linux/amd64 -f docker/mars/Dockerfile remote_bin -t asia.gcr.io/$PROJECT_NAME/mars:$IMAGE_TAG
docker build --platform=linux/amd64 -f docker/http/Dockerfile remote_bin -t asia.gcr.io/$PROJECT_NAME/http:$IMAGE_TAG
docker build --platform=linux/amd64 -f docker/tcp/Dockerfile remote_bin -t asia.gcr.io/$PROJECT_NAME/tcp:$IMAGE_TAG
docker build --platform=linux/amd64 -f docker/udp/Dockerfile remote_bin -t asia.gcr.io/$PROJECT_NAME/udp:$IMAGE_TAG

# Container RegistryへPush
docker push asia.gcr.io/$PROJECT_NAME/mars:$IMAGE_TAG
docker push asia.gcr.io/$PROJECT_NAME/http:$IMAGE_TAG
docker push asia.gcr.io/$PROJECT_NAME/tcp:$IMAGE_TAG
docker push asia.gcr.io/$PROJECT_NAME/udp:$IMAGE_TAG
```


## インフラ構築手順7 - GKE クラスタへのマニフェスト反映
kustomize を使用し、GKE クラスタにマニフェストの反映をします
<!-- Todo: fix manifests -->
```
kustomize build ./src/k8s/overlays/gcp/ | sed -e "s/__GCP_PROJECT_ID__/${PROJECT_NAME}/g" | kubectl apply -f -
```

## インフラ構築手順8 - Diarkis 起動確認
作成したクラスタの稼働を確認するために、Diarkis の認証エンドポイントに HTTP リクエストを送信し、動作の確認をします。
下記の結果がレスポンスされれば正常に起動が完了しております。
(下記の例では xxx や yyy 等で伏せ字をしております)
```
# http loadbalancer のエンドポイント取得
$ kubectl get service http -n dev0

$ curl {http service の EXTERNAL-IP}:7000/auth/1 # 本来は、"GET /auth/{user id}"でリクエストいただく想定です。今回は動作確認のため"/auth/1"で挙動確認しております。
{"TCP":"xxxx.bc.googleusercontent.com:7200","UDP":"yyyy.bc.googleusercontent.com:7100","sid":"xxxxx","encryptionKey":"xxxxx","encryptionIV":"xxxxx","encryptionMacKey":"xxxxx"}
```

