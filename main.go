package main

import (
	"log"

	"github.com/Rinlys-sama/AgnosAssignment/config"
	"github.com/Rinlys-sama/AgnosAssignment/routes"
)

func main() {
	// Step 1: Load configuration from environment variables
	cfg := config.LoadConfig()

	// Step 2: Connect to PostgreSQL
	db := config.ConnectDB(cfg)
	defer db.Close() // Close the DB connection when the app shuts down

	// Step 3: Set up routes and start the server
	router := routes.SetupRouter(db, cfg)

	log.Println("Starting server on :8080")
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
