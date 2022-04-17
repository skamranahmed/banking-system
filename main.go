package main

import (
	"database/sql"
	"log"

	"github.com/skamranahmed/banking-system/api"
	"github.com/skamranahmed/banking-system/config"

	db "github.com/skamranahmed/banking-system/db/sqlc"

	_ "github.com/lib/pq"
)

func main() {
	// load config
	config.Load("./config")

	// open the db connection
	log.Printf("🔌 Connecting to database....")

	conn, err := sql.Open(config.DbDriver, config.DbHost)
	if err != nil {
		log.Fatalf("❌ unable to connect to the db, error: %s", err)
	}

	err = conn.Ping()
	if err != nil {
		log.Fatalf("❌ unable to establish db connection, error: %s", err)
	}

	log.Printf("✅ Database connection successful")

	// close the db connection
	defer conn.Close()

	// instantiate dependencies
	store := db.NewStore(conn)
	server, err := api.NewServer(store)
	if err != nil {
		log.Fatalf("unable to instantiate server, error: %v", err)
	}

	err = server.Start(config.ServerPort)
	if err != nil {
		log.Fatalf("unable to start server, error: %v", err)
	}
}
