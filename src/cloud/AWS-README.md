# diarkis-infra-template

## Overview
AWS EKS上にDiarkisクラスターを構築するための手順です。
プリミティブな設定になっているので、適宜修正したい点があれば修正していただければと思います。

## prerequisites
- 課金が有効になっているawsアカウント
- awsコマンドの認証が通っていること
- kustomize@v4.5.7が使用可能であること

## 1. install eksctl
https://catalog.us-east-1.prod.workshops.aws/workshops/f5abb693-2d87-43b5-a439-77454f28e2e7/ja-JP/020-create-cluster/10-install-eksctl

## 2. create ECR for diarkis images
Diarkis構成コンポーネントをpushするためのregistryを準備
alpineなどもsampleで使用しているが、それに関してはdocker hubから取得
```
aws ecr create-repository --repository-name http
export HTTP_URI=$(aws ecr describe-repositories --repository-names http | jq '.repositories[].repositoryUri')
aws ecr create-repository --repository-name udp
export UDP_URI=$(aws ecr describe-repositories --repository-names http | jq '.repositories[].repositoryUri')
aws ecr create-repository --repository-name tcp
export TCP_URI=$(aws ecr describe-repositories --repository-names http | jq '.repositories[].repositoryUri')
aws ecr create-repository --repository-name mars
export MARS_URI=$(aws ecr describe-repositories --repository-names http | jq '.repositories[].repositoryUri')
```

## 3. Create EKS for diarkis
```
eksctl create cluster --node-type c5n.xlarge --name diarkis --nodes 1 --alb-ingress-access --asg-access # about 10 minutes
```

## 4. connect to eks
```
aws eks --region ap-northeast-1 update-kubeconfig --name diarkis # get credetial for k8s
```

## 5. Open EKS firewall
EKSのNodeに対して0.0.0.0/0からtcp,udpの7000-8000を開放する

## 6. tagging the server image and push
server-templateから生成した project の root から下記を実行
※ 詳細は[こちら](https://help.diarkis.io/ja/running-diarkis-server-on-local)をご覧ください。
```
make build-local
```
./remote_bin にサーバーの実行ファイル郡(udp, tcp, http, mars)が生成された後、コンテナイメージをビルド
```
export AWS_ACCOUNT_NUM=$(aws sts get-caller-identity | jq .Account -r)
docker build -t ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/udp ./remote_bin -f docker/udp/Dockerfile
docker build -t ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/tcp ./remote_bin -f docker/tcp/Dockerfile
docker build -t ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/http ./remote_bin -f docker/http/Dockerfile
docker build -t ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/mars ./remote_bin -f docker/mars/Dockerfile
```
dockerに認証を通す
```
aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com
```
imageをpush
```
docker push ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/http
docker push ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/udp
docker push ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/tcp
docker push ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/mars
```

## 6. apply manifest
```
kustomize build k8s/overlays/dev0 | sed -e "s/__AWS_ACCOUNT_NUM__/${AWS_ACCOUNT_NUM}/g" | kubectl apply -f -
```
下記のように4つのコンポーネントが立ち上がっていればOKです。
```
$ kubectl get po -n dev0
NAME                    READY   STATUS    RESTARTS   AGE
http-5c7dbbb6d7-lhjlm   1/1     Running   0          3d14h
mars-0                  1/1     Running   0          3d14h
tcp-88dc5f97d-7sqk9     1/1     Running   0          3d14h
udp-fdc6bbccc-dwc5w     1/1     Running   0          3d14h
```
## 7. check diarkis cluster
まずpublic endpointを取得する
```
EXTERNAL_IP=$(kubectl get svc http -o json -n dev0 | jq '.status.loadBalancer.ingress[].hostname')
kubectl get svc -n dev0 -o wide # このコマンドで表示されるEXTERNAL IPと同一なのでどちらで見ていただいても構いません。
```

取得できたEXTERNAL-IPに対してAPIを叩く
```
curl ${EXTERNAI_IP}/auth/1
```
下記の様なレスポンスが返ってくればOK
```
{"TCP":"ec2-52-197-27-16.ap-northeast-1.compute.amazonaws.com:7201","UDP":"ec2-52-197-27-16.ap-northeast-1.compute.amazonaws.com:7101","sid":"6a970e7a66a24d1e998fe5211e11890b","encryptionKey":"59ccc205e9a94e11a17a59c601669102","encryptionIV":"0167b0e1c1e24ff39d3150dae640f67f","encryptionMacKey":"197dc161f4c44f829ff9712805ab6b36"}%
```
抜けている項目等があれば、何かのコンポーネントに異常をきたしている可能性があるため、お問合せください。

