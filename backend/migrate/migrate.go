// Author: asakura
// Description: マイグレーションを行う関数を定義する。
// マイグレーション: スキーマ変更をDBに対して適用すること。

package main

import (
	"backend/db"
	"backend/model"
	"fmt"
)

func main() {
	dbConn := db.ConnectDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.Users{}, &model.Posts{})
}
