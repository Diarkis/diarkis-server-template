# GKE Autopilot 構築手順と概要

Diarkis サーバークラスターを Google Kubernetes Engine (GKE) Autopilot にデプロイするためのマニフェスト群です。

## 概要

GKE Autopilot 環境に最適化された Diarkis サーバー（HTTP, TCP, UDP, Mars）をデプロイします。
Kustomize を使用して、共通設定（base）と環境ごとの設定（overlays/dev0）を管理しています。

## ディレクトリ構成

- `base/` : 全環境共通の基本マニフェスト
- `overlays/dev0/` : 開発環境用の設定パッチ

## コンポーネント

- **HTTP**: `Deployment` としてデプロイ。HPA によるオートスケーリングが有効。
- **TCP/UDP**: `Pod` として個別にデプロイ。GKE Autopilot の `hostPort` 割り当て機能を利用して外部からの通信を受け付けます。
- **Mars**: `StatefulSet` としてデプロイ。サーバークラスター内のノード管理と同期を行います。
- **Monitoring**: Google Cloud Managed Service for Prometheus 用の `PodMonitoring` 設定が含まれています。

## 構築・デプロイ手順

### 1. 前提条件

- GKE Autopilot クラスターが起動していること。
- Google Artifact Registry 等に Diarkis の各コンポーネント（HTTP, TCP, UDP, Mars）のイメージがプッシュされていること。
- `kubectl` および `kustomize` がローカル環境にインストールされていること。

### 2. イメージ情報の修正

`src/k8s/gcp-autopilot/base/kustomization.yaml` の `images` セクションを環境に合わせて修正します。

```yaml
images:
  - name: udp
    newName: [YOUR_REGION]-docker.pkg.dev/[PROJECT_ID]/[REPO_NAME]/udp
    newTag: [TAG]
  # ... 他のコンポーネントも同様
```

### 3. デプロイの実行

以下のコマンドで、名前空間の作成から各リソースのデプロイまでを一括で行います。

```bash
kubectl apply -k src/k8s/gcp-autopilot/overlays/dev0
```

### 4. 動作確認

デプロイされた Pod が正常に稼働しているか確認します。

```bash
kubectl get pods -n dev0
```

## 注意事項

- **UDP/TCP のポート開放**: GKE Autopilot では `hostPort` を使用する場合、アノテーション `autopilot.gke.io/host-port-assignment` を使用してポート範囲を指定する必要があります。本テンプレートでは設定済みです。
- **リソース制限**: Autopilot の制限により、最小リソース要件（CPU 0.25 vCPU 等）があります。マニフェストのリソース指定はこれに従って調整してください。
