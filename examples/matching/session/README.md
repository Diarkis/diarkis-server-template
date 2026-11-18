# 概要

このサンプルは **(3)** つのサーバーと **(1)** つのクライアントで構成されています。

- **MARS** サーバー: サーバー同士をつなぎ、全体をまとめます。
- **HTTP** サーバー: ユーザーの認証を行います。
- **UDP** サーバー: クライアントの接続を受け、Session の処理を行います。

このサンプルでは、Diarkis の **Session** の使い方と、Session の **コールバック** が
いつ呼ばれるかを確認できます。

## Session とは

Session は、ユーザーをひとつのグループにまとめる機能です。
メンバーが同じサーバーに接続している必要はありません。
別々のサーバーに接続しているユーザー同士でも、同じ Session に集まることができます。

- Session を作ったユーザーがオーナーになります。
- 他のユーザーは、招待を受けるか、Session ID を指定して参加することができます。
- メンバー全員、または特定のメンバーにメッセージを送れます。
- オーナーが抜けると、残ったメンバーからランダムで新しいオーナーが選ばれます。
- メンバーがいなくなった Session は、TTL (秒) がたつと自動で削除されます。

**NOTE**: このサンプルは、わかりやすさのため UDP サーバー **(1)** 台で動かしています。

### Room との違い

|                  | Session                                        | Room                                                         |
| ---------------- | ---------------------------------------------- | ------------------------------------------------------------ |
| メンバーの接続先 | 別々のサーバーのままでよい                     | Room があるサーバーに集まる (必要なら自動で接続し直す)       |
| 参加のしかた     | 招待を受ける、または Session ID を指定して参加 | Room ID を指定して参加                                       |
| メッセージ       | 全員への broadcast、特定のメンバーへの message | broadcast、message、relay など。まとめて送る間隔を設定できる |
| 状態の同期       | プロパティを共有できる                         | プロパティに加えて、状態の自動同期 (State Sync) がある       |

## Session のコールバック

サーバー側では、Session ごとに次のコールバックを登録できます。
このサンプルでは、どのコールバックも UDP サーバーのログに出力するようにしています。

| コールバック | 呼ばれるタイミング                                             |
| ------------ | -------------------------------------------------------------- |
| `OnJoin`     | ユーザーが Session に入る **前**。                             |
| `OnJoined`   | ユーザーが Session に入った **後**。                           |
| `OnLeft`     | ユーザーが Session から抜けた後。kick されたときも呼ばれます。 |
| `OnDeleted`  | Session が削除された後。                                       |

このサンプルでは、kick されてから 60 秒以内のユーザーを拒否しています。

**NOTE**: Session を作ったユーザー自身には `OnJoin` / `OnJoined` は呼ばれません。

## ビルド方法

付属の Mage スクリプトで、**(3)** つのサーバーとクライアントをまとめてビルドできます。

### Linux / macOS

```sh
./run-mage.sh build:local
```

### Windows

```sh
.\run-mage.bat build:local
```

**MARS**、**HTTP**、**UDP** サーバーとクライアントのバイナリが `remote_bin` ディレクトリに作られます。

## 起動方法

クライアントを使う前に、**(3)** つのサーバーをすべて起動してください。

### 1. MARS サーバーを起動する

```sh
./run-mage.sh server mars
```

### 2. HTTP サーバーを起動する

```sh
./run-mage.sh server http
```

### 3. UDP サーバーを起動する

```sh
./run-mage.sh server udp
```

サーバーがすべて起動したら、別々のターミナルでクライアントを **(3)** つ起動します。

```sh
./remote_bin/cli -uid user1 -puffer true
```

```sh
./remote_bin/cli -uid user2 -puffer true
```

```sh
./remote_bin/cli -uid user3 -puffer true
```

各コマンドで聞かれる項目 (TCP/UDP、Session type など) は、
特に書いていなければ何も入力せず Enter (デフォルト) で進めます。

## シナリオ 1: 招待・参加・退出・kick

**1. user1 が Session を作成**

```sh
# user1 (Terminal 1)
> session create
Enter for which protocol to a create a new Session (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
Enter max members [1 - 255] (uint8): (Default: 10)
Enter TTL (seconds) [30 - 255] (uint8): (Default: 30)
Session created. Created: {sessionType: 0, sessionID: bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be}
```

UDP サーバーのログ:

```sh
callbacks registered {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be"}
```

**2. user1 が user2 と user3 を招待**

```sh
# user1 (Terminal 1)
> session invite
Enter for which protocol send a Session invitation(s) (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
Enter invitee UIDs (',' for multiple entries) ([]string): (Example: 0000,0001,0002) user2,user3
Enter message (string): hi
```

```sh
# user2 (Terminal 2) / user3 (Terminal 3)
You have been invited to join a Session. Invitation: {sessionType: 0, sessionID: bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be, msg: hi}
Reply with `invite yes` or `invite no`
```

**3. user2 が招待を受ける**

```sh
# user2 (Terminal 2)
> invite yes
Invite reply sent: joined
A new member has joined the Session. Joined: {sessionType: 0, memberID: user2}
```

UDP サーバーのログ:

```sh
OnJoin {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be","uid":"user2"}
OnJoined {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be","uid":"user2"}
```

**4. user3 が招待を断る**

```sh
# user3 (Terminal 3)
> invite no
Invite reply sent: declined
```

```sh
# user1 (Terminal 1)
user3 declined the invitation. SessionID: bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be
```

> 断られたことは、招待した user1 にだけ通知されます (user2 には届きません)。

**5. user3 が招待なしで参加**

```sh
# user3 (Terminal 3)
> session join
Enter for which protocol to join a Session (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
Enter Session ID (string): bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be
A new member has joined the Session. Joined: {sessionType: 0, memberID: user3}
```

UDP サーバーのログ:

```sh
OnJoin {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be","uid":"user3"}
OnJoined {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be","uid":"user3"}
```

**6. user1 がメンバー全員にメッセージを送る**

```sh
# user1 (Terminal 1)
> session broadcast
Enter for which protocol to broadcast to Session (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
Enter message (string): hi
Session broadcast received. Broadcast: {sessionType: 0, msg: hi}
```

```sh
# user2 (Terminal 2) / user3 (Terminal 3)
Session broadcast received. Broadcast: {sessionType: 0, msg: hi}
```

**7. user1 が user2 にだけメッセージを送る**

```sh
# user1 (Terminal 1)
> session message
Enter for which protocol send a Session message (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
Enter recipient UIDs (',' for multiple entries) ([]string): (Example: 0000,0001,0002) user2
Enter message (string): hi
```

```sh
# user2 (Terminal 2)
Session message received. Message: {sessionType: 0, msg: hi}
```

**8. オーナー (user1) が抜ける**

```sh
# user1 (Terminal 1)
> session leave
Enter for which protocol to leave from a Session (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
```

```sh
# user2 (Terminal 2) / user3 (Terminal 3)
Member left the Session. Left: {sessionType: 0, memberID: user1}
```

UDP サーバーのログ:

```sh
OnLeft {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be","uid":"user1"}
```

**9. user2 が Session の情報を取得**

Session の情報から、新しいオーナー (`ownerID`) を確認できます。

```sh
# user2 (Terminal 2)
> session get info type
Enter for which protocol get Session info (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
Session info: {sessionType: 0, sessionID: bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be, currentMembers: 2/10, MemberIDs: [user2 user3], ownerID: user2}
```

> `ownerID` は user2 と user3 のどちらになるかわかりません。

**10. オーナーがもう一人を kick**

```sh
> session kick
Enter for which protocol to kick member from a Session (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
Enter user ID (string): user2 or user3
Member left the Session. Left: {sessionType: 0, memberID: user3}
```

```sh
You have been kicked from the Session. Kick: {sessionType: 0}
```

UDP サーバーのログ:

```sh
OnLeft {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be","uid":"user2 or user3"}
```

**11. kick されたユーザーが60秒以内にもう一度参加**

```sh
> session join
Enter for which protocol to join a Session (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
Enter Session ID (string): bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be
Command failed: SessionJoinError
Failed to join -> Type:0 SessionID:bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be UserID:user3
```

UDP サーバーのログ:

```sh
OnJoin rejected: kicked recently {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be","uid":"user2 or user3"}
```

**12. 60 秒たってから kick されたユーザーがもう一度参加**

```sh
> session join
Enter for which protocol to join a Session (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
Enter Session ID (string): bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be
A new member has joined the Session. Joined: {sessionType: 0, memberID: user2 or user3}
```

UDP サーバーのログ:

```sh
OnJoin {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be","uid":"user2 or user3"}
OnJoined {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be","uid":"user2 or user3"}
```

---

## シナリオ 2: Session 一覧と削除

このシナリオでは、Session が削除されたときに `OnDeleted` が呼ばれ、一覧から消えることを確認します。
シナリオ 1 の後に続ける場合は、全員 `session leave` してから始めてください。

**1. user1 と user2 がそれぞれ Session を作成**

```sh
# user1 (Terminal 1)
> session create
Session created. Created: {sessionType: 0, sessionID: bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be}

# user2 (Terminal 2)
> session create
Session created. Created: {sessionType: 0, sessionID: 2be6ecd74e044c8d93466877ce7ae5ce7f0000011fa58cd04e87}
```

UDP サーバーのログに、Session ごとに `callbacks registered` が出ます。

**2. user3 が一覧を見る**

```sh
# user3 (Terminal 3)
> custom show
Sessions (2):
  type: 0, id: 2be6ecd74e044c8d93466877ce7ae5ce7f0000011fa58cd04e87, members: [user2]
  type: 0, id: bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be, members: [user1]
```

**3. user3 が片方の Session に参加**

```sh
# user3 (Terminal 3)
> session join
Enter for which protocol to join a Session (TCP/UDP): (Default: UDP)
Enter Session type [0 - 255] (uint8): (Default: 0)
Enter Session ID (string): bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be
A new member has joined the Session. Joined: {sessionType: 0, memberID: user3}

> custom show
Sessions (2):
  type: 0, id: 2be6ecd74e044c8d93466877ce7ae5ce7f0000011fa58cd04e87, members: [user2]
  type: 0, id: bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be, members: [user1 user3]
```

**4. 全員が抜ける**

```sh
# user1 (Terminal 1) / user2 (Terminal 2) / user3 (Terminal 3)
> session leave
```

```sh
# user3 (Terminal 3)
> custom show
Sessions (2):
  type: 0, id: 2be6ecd74e044c8d93466877ce7ae5ce7f0000011fa58cd04e87, members: []
  type: 0, id: bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be, members: []
```

> メンバーがいなくなっても、Session はすぐには削除されません。
>
> **注意**: TTL のデフォルトは 30 秒なので、Session は早く削除されます。
> 全員が抜けたら、すぐに `custom show` を実行してください。
> ゆっくり確認したい場合は、手順 1 の `session create` で TTL を長め (例: `120`) に設定してください。

**5. TTL がたってから一覧を見る**

TTL 秒ほど待つと、UDP サーバーのログに Session ごとに `OnDeleted` が出ます。

```sh
OnDeleted {"sessionID":"bcda3bac44a34341b33f930c8c33cf547f0000011fa5c8a0e6be"}
OnDeleted {"sessionID":"2be6ecd74e044c8d93466877ce7ae5ce7f0000011fa58cd04e87"}
```

```sh
# user3 (Terminal 3)
> custom show
Sessions (0):
```
