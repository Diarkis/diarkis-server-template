# Overview

このプロジェクトは、サーバーとクライアントで構成されています。

- **MARS** ノードが構築しているメッシュネットワークをオーケストレーションします。

- **HTTP** サーバーは、クライアントの認証を行います。

- **UDP** サーバーは、クライアントの接続を管理し、すべての入力 Room と Room LoD コマンドを処理します。

このサンプルのゴールは Room 内で距離によって遠方の同期を減らすことで、クライアントの描画と通信の負荷を減らすことです。

## How to Build

サーバーとクライアントは mage を使用してビルドします。

### To build on Linux or macOS

```sh
./run-mage.sh build:local
```

### To build on Windows

```sh
.\run-mage.bat build:local
```

これにより、**MARS**、**HTTP**、**UDP**サーバーとクライアントバイナリが作成され、`remote_bin`ディレクトリに配置されます。

## How to Run

このプロジェクトは、すべてのサーバーが動作している必要があります。

### 1. First, start the MARS server to orchestrate the node mesh

```sh
./run-mage.sh server mars
```

### 2. Next, start the HTTP server, which holds the custom matchmaking criteria

```sh
./run-mage.sh server http
```

### 3. Then, start the UDP server, to manage incoming client connections

```sh
./run-mage.sh server udp
```

## Limit

このサンプルのサーバーは必要最低限のモジュールだけ有効にしているので、Lod 以外の機能をテストする際には、必要なモジュールを有効にすること、もしくは有効にしたサーバーに接続してください。

## Testing Room LoD

ユーザー 1 が UDP サーバーに接続して Room を作成します。

```sh
./remote_bin/cli -host=127.0.0.1:7000 -uid=user11
Connected UDP

> room create
Enter for which protocol to a create a Room (TCP/UDP): (Default: UDP)
Enter max members [1 - 255] (uint16): (Default: 10)
Enter if allow empty (y/n): (Default: no)
Enter if join on creation (y/n): (Default: yes)
Enter TTL (seconds) [10 - 255] (uint8): (Default: 30)
Enter broadcast interval (milliseconds) (uint8): (Default: 100)
Room created. roomID:cf7b49ea1fc836d87f0000011fa5000000000000000000000000, createdAt:1764084948
```

ユーザー 2 が UDP サーバーに接続して Room に参加します。

```sh
./remote_bin/cli -host=127.0.0.1:7000 -uid=user12
Connected UDP

> room join
Enter for which protocol to a join a Room (TCP/UDP): (Default: UDP)
Enter Room ID (52 characters) (string): cf7b49ea1fc836d87f0000011fa5000000000000000000000000
Enter Room join message (string): (Default: Hello from user12!)
UDP Room cf7b49ea1fc836d87f0000011fa5000000000000000000000000 joined and it was created at 1764084948
New member joined Room. Message: Hello from user12!
```

ユーザー 1 が Room の LoD 情報を取得します。

```sh
> lod getinfo
Get LoD Info successful: MaxDistanceForFar = 40000 | MaxDistanceForNearby = 10000 | SyncIntervalForFar = 2000 | SyncIntervalForNearby = 16
```

ユーザー 1 が Room に LoD 情報をブロードキャストします。

```sh
> lod broadcast
Enter X (int32):
100
Enter Y (int32):
200
Enter Payload (string):
sync data
Broadcast LoD successful:
Broadcast LoD Push successful:  sync data
```

ユーザー 2 がユーザー 1 から適切な同期間隔でデータを受信します。

```sh
Broadcast LoD Push successful:       sync data
```

LoD コマンドを表示するには、以下のコマンドを実行します。

```sh
 > help lod
================ Command List ================
lod broadcast - Broadcast the LoD
lod getinfo   - Get the LoD configuration
=============================================
```

---

### 距離による同期の仕組み

サーバーは送信者と受信者の距離を計算し、以下のルールで同期間隔を決定します：

#### 1. 近距離ユーザー（距離 ≤ MaxDistanceForNearby）

- **同期間隔**: `SyncIntervalForNearby`（デフォルト: 16ms）
- **例**: 距離が 10,000cm 以下の場合、高頻度で同期

#### 2. 中距離ユーザー（MaxDistanceForNearby < 距離 ≤ MaxDistanceForFar）

- **同期間隔**: 距離に応じて線形補間
- **計算式**:
  ```
  interval = SyncIntervalForNearby +
             (SyncIntervalForFar - SyncIntervalForNearby) ×
             (距離 - MaxDistanceForNearby) /
             (MaxDistanceForFar - MaxDistanceForNearby)
  ```
- **例**: 距離が 25,000cm の場合、約 1,000ms 間隔で同期

#### 3. 遠距離ユーザー（距離 > MaxDistanceForFar）

- **同期間隔**: 同期なし
- **例**: 距離が 40,000cm を超える場合、データは送信されない

> **Note**: 距離はマンハッタン距離（`|X1 - X2| + |Y1 - Y2|`）で計算されます。

---

# 実践的な使用例 / Practical Use Cases

## シナリオ 1: 3 人のプレイヤーによる距離ベースの同期

このシナリオでは、3人のプレイヤー（Alice, Bob, Charlie）が異なる位置にいる場合の同期動作を確認します。

### 設定

- **MaxDistanceForNearby**: 10,000cm (100m)
- **MaxDistanceForFar**: 40,000cm (400m)
- **SyncIntervalForNearby**: 33ms
- **SyncIntervalForFar**: 2000ms

### プレイヤーの位置

- **Alice**: (0, 0)
- **Bob**: (5000, 3000) → Alice からの距離: 8,000cm（近距離）
- **Charlie**: (20000, 15000) → Alice からの距離: 35,000cm（中距離）

### テスト手順

**1. 各プレイヤーが Room に参加**

```sh
# Alice (Terminal 1)
./remote_bin/cli -host=127.0.0.1:7000 -uid=alice
> room create
> lod getinfo

# Bob (Terminal 2)
./remote_bin/cli -host=127.0.0.1:7000 -uid=bob
> room join
# Enter Room ID from Alice

# Charlie (Terminal 3)
./remote_bin/cli -host=127.0.0.1:7000 -uid=charlie
> room join
# Enter Room ID from Alice
```

**2. Alice が位置情報をブロードキャスト**

```sh
# Alice (Terminal 1)
> lod b
Enter X (int32):
0
Enter Y (int32):
0
Enter Payload (string):
Alice at origin
```

**3. 期待される動作**

- **Bob**: 距離 8,000cm（近距離）→ 約 16ms 間隔で「Alice at origin」を受信
- **Charlie**: 距離 35,000cm（中距離）→ 約 1,666ms 間隔で「Alice at origin」を受信

**4. Bob が移動して遠距離になった場合**

```sh
# Bob (Terminal 2)
> lod b
Enter X (int32):
50000
Enter Y (int32):
0
Enter Payload (string):
Bob moved far away
```

この時、Bob と Alice の距離は 50,000cm（遠距離超過）となり、相互に同期が停止します。

---

## シナリオ 2: 移動するプレイヤーの連続ブロードキャスト

実際のゲームでは、プレイヤーが連続的に位置を更新します。このシナリオでは、移動中のプレイヤーがどのように同期されるかを確認します。

### テスト手順

**1. Player1 が原点に立つ**

```sh
# Player1
> lod b
Enter X (int32):
0
Enter Y (int32):
0
Enter Payload (string):
{"status": "idle", "hp": 100}
```

**2. Player2 が近づきながら複数回ブロードキャスト**

```sh
# Player2 - 位置 1（遠距離）
> lod b
Enter X (int32):
45000
Enter Y (int32):
0
Enter Payload (string):
{"status": "walking", "direction": "west"}

# Player2 - 位置 2（中距離）
> lod b
Enter X (int32):
25000
Enter Y (int32):
0
Enter Payload (string):
{"status": "walking", "direction": "west"}

# Player2 - 位置 3（近距離）
> lod b
Enter X (int32):
5000
Enter Y (int32):
0
Enter Payload (string):
{"status": "walking", "direction": "west"}
```

**3. 期待される動作**

- 位置 1（45,000cm）: Player1 には同期されない
- 位置 2（25,000cm）: Player1 に約 1,000ms 間隔で同期
- 位置 3（5,000cm）: Player1 に約 33ms 間隔で同期（滑らかな動き）

---

## Tips & Best Practices

### 1. 座標の単位に注意

- 座標の単位は**センチメートル（cm）**です
- 1m = 100cm、10m = 1,000cm、100m = 10,000cm

### 2. データ最適化

- 位置が変わらない場合でも、データが更新されない場合は遠距離の同期間隔が適用されます
- 頻繁に更新する必要があるデータのみペイロードに含める

### 3. デバッグ時のヒント

- `lod getinfo`でサーバー設定を確認してから、テストを開始
- 複数のターミナルを開いて、各プレイヤーの視点で同期を確認

---

# 仕様 / Specification

LoD (Level of Detail) は、距離によって同期の頻度を変えることで、クライアントの描画と通信の負荷を減らすための仕組みです。
このサンプルでは、近い人と遠い人で、それぞれの更新頻度と距離範囲を設定することができます。
その上で、クライアントから送信したデータをサーバー側で自動的に適切な頻度で相手に送信することが可能です。

## 設定

./configs/shared/room.json で以下の設定を行うことができます。

```json
{
  "syncIntervalForNearby": 16,
  "syncIntervalForFar": 2000,
  "maxDistanceForNearby": 1000,
  "maxDistanceForFar": 40000
}
```

| key                   | description                       | default value |
| --------------------- | --------------------------------- | ------------- |
| syncIntervalForNearby | 近距離の同期間隔 (millisecond)    | 16            |
| syncIntervalForFar    | 遠距離の同期間隔 (millisecond)    | 2000          |
| maxDistanceForNearby  | 近距離の最大同期距離 (centimeter) | 10000         |
| maxDistanceForFar     | 遠距離の最大同期距離 (centimeter) | 40000         |

---

## コマンド仕様

### BroadcastLoD

Room に、自分のキャラクターの x, y 座標と同期したいデータを送信します。

自キャラクターの座標を元に、以下の条件で相手にデータを送信します。

- 近い位置の相手 ( `maxDistanceForNearby` 以下 ) には頻繁に送信 ( `syncIntervalForNearby` )
- 遠い位置の相手 ( `maxDistanceForFar` 以下 ) には低頻度で送信 ( `syncIntervalForFar` )
- それ以外の相手 ( `maxDistanceForFar` 以上 ) には送信しません

#### Command

Ver=2, Cmd=1001

#### Parameters

| parameter | type   | description                 |
| --------- | ------ | --------------------------- |
| x         | int32  | 自分のキャラクターの x 座標 |
| y         | int32  | 自分のキャラクターの y 座標 |
| payload   | byte[] | 同期したいデータ            |

#### Response

なし

#### Push Notification

近いユーザーと遠いユーザーで異なる間隔でデータをプッシュ通知で受信します。

**Command**: Ver=2, Cmd=1001

**Parameters**:

| parameter | type   | description |
| --------- | ------ | ----------- |
| payload   | byte[] | 同期データ  |

#### Errors

TODO: エラー定義

- Invalid Payload
- Not in Room

---

### GetLoDInfo

Room の LoD 設定を取得します。

#### Command

Ver=2, Cmd=1002

#### Parameters

なし

#### Response

**Command**: Ver=2, Cmd=1002

**Parameters**:

| parameter             | type   | description                       |
| --------------------- | ------ | --------------------------------- |
| syncIntervalForNearby | uint32 | 近距離の同期間隔 (millisecond)    |
| syncIntervalForFar    | uint32 | 遠距離の同期間隔 (millisecond)    |
| maxDistanceForNearby  | uint32 | 近距離の最大同期距離 (centimeter) |
| maxDistanceForFar     | uint32 | 遠距離の最大同期距離 (centimeter) |

#### Push Notification

なし

---

_First created on 2025-11-26. Updated on 2025-12-04._
