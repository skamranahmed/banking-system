package main

import (
	"database/sql"
	"log"

	"github.com/skamranahmed/banking-system/api"
	db "github.com/skamranahmed/banking-system/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://postgres:password@localhost:5432/bank?sslmode=disable"
	serverAddress = "localhost:8080"
)

func main() {
	// open the db connection
	conn, err := sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("unable to connect to the db, error: %s", err)
	}

	// close the db connection
	defer conn.Close()

	// instantiate dependencies
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)
	if err != nil {
		log.Fatalf("unable to start server, error: %s", err)
	}
}
