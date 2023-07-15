package sqlc

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testBD *sql.DB

const (
	dbDriver = "postgres"
	dbSource = "postgresql://postgres:oneanhiuemlove33@localhost:5432/simple_shop?sslmode=disable"
)

func TestMain(m *testing.M) {
	var err error

	testBD, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	testQueries = New(testBD)

	os.Exit(m.Run())
}
