# Overview
This is a terraform files and manifests and manuals for deploying a Diarkis cluster sets in AWS.

# Structure

```
.
├── ProcedureManual.md # インフラの構築手順書
├── README # this file
├── k8s # diarkis や cluster-autoscaler の manifest を入れています
├── prometheus # EKS cluster にインストールしている prometheusServer のマニフェストや、権限付与 scirpt 等を格納
└── terraform # インフラ構築のための terraform ファイルを格納しています
```
