package sqlc

import (
	"database/sql"
	"log"
	"os"
	"simple_shop/db/util"
	"testing"

	_ "github.com/lib/pq"
)

var testQueries *Queries
var testBD *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config, ", err)
	}

	testBD, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	testQueries = New(testBD)

	os.Exit(m.Run())
}
