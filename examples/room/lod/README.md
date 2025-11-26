# Overview

This project is comprised of **(3)** servers and **(1)** client.

- **MARS** server is a standard (out-of-the-box) Diarkis template. It serves to orchestrates the
  node mesh.

- **HTTP** server is a standard (out-of-the-box) Diarkis template; It serves to authenticate the
  client.

- **UDP** server hosts the client connection and handles all incoming Room and Room LoD commands.

このサンプルのゴールは Room 内で距離によって遠方の同期を減らすことで、クライアントの描画と通信の負荷を減らすことです。

## How to Build

You can build all **(3)** servers and the client binary using the provided Mage build scripts.

### To build on Linux or macOS

```sh
./run-mage.sh build:local
```

### To build on Windows

```sh
.\run-mage.bat build:local
```

This will create the **MARS**, **HTTP**, and **UDP** server binaries, and client binary,
placing them inside the `remote_bin` directory.

## How to Run

This project requires all **(3)** servers to be running before clients can test matchmaking
behavior.

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

## Testing Room LoD

User 1 connects to the UDP server and creates a room.

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

User 2 connects to the UDP server and joins the room.

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

User 1 gets the LoD information of the room.

```sh
> lod getinfo
Get LoD Info successful: MaxDistanceForFar = 40000 | MaxDistanceForNearby = 10000 | SyncIntervalForFar = 2000 | SyncIntervalForNearby = 16
```

User 1 broadcasts the LoD information to the room.

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

User 2 receives data from User 1 at the appropriate synchronization interval.

```sh
Broadcast LoD Push successful:       sync data
```

You can see the LoD commands by running the following command:

```sh
 > help lod
================ Command List ================
lod broadcast - Broadcast the LoD
lod getinfo   - Get the LoD configuration
=============================================
```

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

- 近い位置の相手 ( `maxDistanceForNearby` 以下 ) には頻繁に ( `syncIntervalForNearby` )
- 遠い位置の相手 ( `maxDistanceForFar` 以下 ) には一定時間間隔に ( `syncIntervalForFar` )
- それ以外の相手 ( `maxDistanceForFar` 以上 ) には送信しません。

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

#### Errors

TODO: エラー定義

---

_First created on 2025-11-26. Updated on 2025-11-26._
