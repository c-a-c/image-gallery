// migrate.go

package main

import (
	"backend/infrastructure/db"
	"backend/model"
	"fmt"
)

func main() {
	dbConn := db.ConnectDB()
	defer fmt.Println("Successfully Migrated")
	defer db.CloseDB(dbConn)
	dbConn.AutoMigrate(&model.Users{}, &model.Posts{})
}
