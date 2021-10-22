# diarkis-infra-template

## Overview
AWS EKS上にDiarkisクラスターを構築するための手順を示す。
プリミティブな設定になっているので、適宜修正したい点があれば修正していただければと思います。
1. install eksctl
2. Create ECR
3. Create EKS
4. Apply k8s manifest to EKS

## prerequisites
- 課金が有効になっているawsアカウント
- awsコマンドの認証が通っていること
- kustomizeが使用可能であること

## 1. install eksctl
## 2. create ECR for diarkis images
Diarkis構成コンポーネントをpushするためのregistryを準備
alpineなどもsampleで使用しているが、それに関してはdocker hubから取得している。
```
aws ecr create-repository --repository-name http
export HTTP_URI=$(aws ecr describe-repositories --repository-names http | jq '.repositories[].repositoryUri')
aws ecr create-repository --repository-name udp
export UDP_URI=$(aws ecr describe-repositories --repository-names http | jq '.repositories[].repositoryUri')
aws ecr create-repository --repository-name tcp
export TCP_URI=$(aws ecr describe-repositories --repository-names http | jq '.repositories[].repositoryUri')
aws ecr create-repository --repository-name mars
export MARS_URI=$(aws ecr describe-repositories --repository-names http | jq '.repositories[].repositoryUri')
aws ecr create-repository --repository-name ws
export WS=$(aws ecr describe-repositories --repository-names http | jq '.repositories[].repositoryUri')
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
http, udp, tcp, mars, wsと名付けたdeploy用のコンテナイメージがローカルに準備されている想定
build時に直接tagを付けていただいてももちろん問題ありません。
```
export AWS_ACCOUNT_NUM=$(aws sts get-caller-identity | jq .Account -r)
docker tag http ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/http
docker tag udp ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/udp
docker tag tcp ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/tcp
docker tag mars ${AWS_ACCOUNT_NUM}.dkr.ecr.ap-northeast-1.amazonaws.com/mars
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
```

取得できたEXTERNAL-IPに対してAPIを叩く
```
curl ${EXTERNAI_IP}/auth/1
```
下記の様なレスポンスが返ってくれば正常です。
```
{"TCP":"ec2-xx-xx-xx-xx.ap-northeast-1.compute.amazonaws.com:7201","UDP":"ec2-xx-xx-xx-xx.ap-northeast-1.compute.amazonaws.com:7101","sid":"6a970e7a66a24d1e998fe5211e11890b","encryptionKey":"59ccc205e9a94e11a17a59c601669102","encryptionIV":"0167b0e1c1e24ff39d3150dae640f67f","encryptionMacKey":"197dc161f4c44f829ff9712805ab6b36"}%
```
抜けている項目等があれば、何かのコンポーネントに異常をきたしていると思われる。

