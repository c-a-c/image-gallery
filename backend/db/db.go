// Author: asakura
// Description: DBの接続と切断を行う関数を定義する。

package db

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBの接続
func ConnectDB() *gorm.DB {
	// 環境変数の読み込み
	if os.Getenv(("GO_ENV")) == "dev" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalln(err)
		}
	}
	
	// urlの指定
	url := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"), 
		os.Getenv("POSTGRES_PW"), 
		os.Getenv("POSTGRES_HOST"), 
		os.Getenv("POSTGRES_PORT"), 
		os.Getenv("POSTGRES_DB"))
	
	// DB接続
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected")
	return db
}

// DBの切断
func CloseDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if err := sqlDB.Close(); err != nil {
		log.Fatalln(err)
	}
}