# Overview
AWS EKS 上に Diarkis クラスターを terraform を用いて構築するための手順です。
ローカルの PC から叩いて構築することを想定しています。
# 使用ツール
- Helm
- Go
- Terraform
- AWS CLI
- kubectl
- Docker
- Git
- kustomize

terraform の version は、v1.10.0 で動作確認済みです。

# 前提
- 構築したい aws アカウントに対して、aws cli でアクセスできるよう認証が通っていて、AWS_PROFILE 環境変数を設定し、向き先を設定しておきます。
- 権限としては、editor 権限で操作できることを想定しています

# ディレクトリ構造
```
.
├── README.md # このドキュメントです。terraform を使った diarkis の動作する EKS の構築方法を記載しています
├── k8s # diarkis や cluster-autoscaler の manifest を入れています
├── prometheus # EKS cluster にインストールしている prometheusServer のマニフェストや、権限付与 scirpt 等を格納
└── terraform # インフラ構築のための terraform ファイルを格納しています
```

# 構築対象
## AWS
- EKS
- VPC
- managed prometheus
- managed grafana
- S3

## EKS 内部
- cluster autoscaler
- prometheus server
- diarkis

# 構築手順
## terraform の backend 用の S3 を AWS コンソールから構築する。

### 設定値
- name: {dev,stg,prd}-diarkis-terraform (dev, stg, prd 等は環境に合わせて適宜設定してください)
- region: ap-northeast-1
- publicAccess: deny
 
といった設定で作成します。

## terraform を使用して、インフラを構築する。
コマンドラインより、下記のように実行します。
```
$ cd terraform
$ terraform init -backend-config="bucket=(dev|stg|mnt|prd)-$YOUR_PROJECT_NAME-diarkis-terraform" -reconfigure
Initializing the backend...
key
  The path to the state file inside the bucket

  Enter a value: ## YOUR_ENV_NAME
$ terraform plan # 差分を確認してください
var.env
  Enter a value: (構築したい環境に合わせて dev, stg, prd などを入力)
$ terraform apply # 実際に構築が始まります。
var. env
  Enter a value: (構築したい環境に合わせて dev, stg, prd などを入力)
...
Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes  
```

## 構築した cluster に接続
`terraform apply` したときに、各種接続コマンドを出力しているので、それを使用して作成したクラスタに接続します。

## cluster autoscaler を install する
cluster autoscaler に必要な権限はすでについているので、
`kubectl apply -f <(curl https://raw.githubusercontent.com/kubernetes/autoscaler/master/cluster-autoscaler/cloudprovider/aws/examples/cluster-autoscaler-autodiscover.yaml | sed 's/<YOUR CLUSTER NAME>/(dev|prd)-diarkis/g')` # 作成している環境に合わせて、dev-diarkis か prd-diarkis にして下さい。
上記を実行していただければ 完了 です。 (各種環境 dev,  prd 等に合わせて変更)

## diarkis application のイメージを作成
生成してプロジェクトのルートに移動し下記を実行
```
make build-container-aws
make push-container-aws
```

## prometheus server を構築する(必要な環境のみ、(dev|prd))
k8s 内に prometheus server を構築し、メトリクスを収集し、メトリクスの保存先としては、Amazon Managed Service for Prometheus (https://aws.amazon.com/jp/prometheus/) を想定しています
```
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add kube-state-metrics https://kubernetes.github.io/kube-state-metrics
helm repo update
kubectl create namespace prometheus
cd prometheus
./install-prometheus.sh
```
設定を変更する場合は、下記を実行していただければ完了です。
```
./update-prometheus.sh
```

## prometheus の動作確認(必要な環境のみ、(dev|prd))
```
cd prometheus
./proxy-promehtues.sh
open http://localhost:9090
```
`Users_UDP_node`` といった Diarkis 固有のメトリクスを含めてメトリクスが取得できれば、適切に設定がなされています。

## grafana の datasource に managed prometheus を追加(必要な環境のみ、(dev|prd))
1. SAML か AWS SSO のどちらかの認証方法を選んで、grafana にログインする。(https://docs.aws.amazon.com/ja_jp/grafana/latest/userguide/authentication-in-AMG.html あたりを参考にして設定していただく。私に権限はないので、設定しておりません。)
2. grafana に対してログイン
3. grafana の左側のペインで、Connections を選択し、prometheus の追加を選択する
4. 右上の add new data source を選択
5. Connection Prometheus Server URL に、 terraform のアウトプットにも出力されている、prometheus-query-endpoint から、`/api/v1/query`を取り除いたものを追加する(最後のパスは自動的に、grafana が追加してくれる)
現状構築した grafana では、datasource の動作確認のために auth0 で認証した Diarkis 社 奥村の ID を登録しています。
6. authentication method は、`SigV4`を選択し、Default Region: `ap-northeast-1`を選択し、test を行い、問題がなければ保存

## diarkis を k8s を用いてデプロイ

```
kustomize build k8s/aws/overlays/dev0 | kubectl apply -f -
```

下記のように 4 つのコンポーネントが立ち上がっていれば OK です。

```
$ kubectl get po -n dev0
NAME                    READY   STATUS    RESTARTS   AGE
http-5c7dbbb6d7-lhjlm   1/1     Running   0          3d14h
mars-0                  1/1     Running   0          3d14h
tcp-88dc5f97d-7sqk9     1/1     Running   0          3d14h
udp-fdc6bbccc-dwc5w     1/1     Running   0          3d14h
```

## check diarkis cluster

まず public endpoint を取得します。

```
EXTERNAL_IP=$(kubectl get svc http -o json -n dev0 | jq -r '.status.loadBalancer.ingress[].hostname')
kubectl get svc -n dev0 -o wide # このコマンドで表示される EXTERNAL IP と同一なのでどちらで見ていただいても構いません。
```

取得できた EXTERNAL-IP に対して HTTP GET リクエストを送信します。

```
curl ${EXTERNAL_IP}/auth/1
```

下記の様なレスポンスが返ってくれば OK です。

```
{"TCP":"ec2-xx-xx-xx-xx.ap-northeast-1.compute.amazonaws.com:7201","UDP":"ec2-yy-yy-yy-yy.ap-northeast-1.compute.amazonaws.com:7101","sid":"xxxxxxxxxx","encryptionKey":"xxxxxxxxxx","encryptionIV":"xxxxxxxxxx","encryptionMacKey":"xxxxxxxxxx"}
```

抜けている項目等があれば、何かのコンポーネントに異常をきたしている可能性があるため、お問合せください。
