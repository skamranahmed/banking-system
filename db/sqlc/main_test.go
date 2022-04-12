package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:password@localhost:5432/bank_test?sslmode=disable"
)

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	var err error
	testDB, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatalf("unable to connect to the db, error: %s", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
