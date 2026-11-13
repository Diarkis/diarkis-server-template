// © 2019-2025 Diarkis Inc. All rights reserved.

// Package protocol サーバーとクライアントで共通して使う、独自コマンドの定義。
package protocol

// 独自コマンドの ver。ver 0〜1 は Diarkis の標準コマンドが使うので、2 を使う。
const CustomVer uint8 = 2

// ShowSessionsCmd Session の一覧を取得するコマンド
const ShowSessionsCmd uint16 = 100

// SessionInfo 一覧に入れる Session 1つ分の情報
type SessionInfo struct {
	ID        string   `json:"id"`
	Type      uint8    `json:"type"`
	MemberIDs []string `json:"memberIDs"`
}
