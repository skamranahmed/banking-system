package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"

	"github.com/skamranahmed/banking-system/config"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// load config
	config.Load("../../config")

	var err error
	testDB, err = sql.Open(config.TestDbDriver, config.TestDbHost)
	if err != nil {
		log.Fatalf("unable to connect to the db, error: %s", err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
