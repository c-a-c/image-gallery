package main

import (
	"backend/infrastructure/db"
	"backend/infrastructure/router"
	"log"
	"os"
)

func main() {
	// Connect to database
	database := db.ConnectDB()
	defer db.CloseDB(database)

	// Setup routes
	e := router.SetupRoutes(database)

	// Get port from environment variable
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("Server starting on port %s", port)
	if err := e.Start(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
