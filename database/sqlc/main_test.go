package database

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/glebarez/go-sqlite"
)

var testQueries *Queries

func TestMain(m *testing.M) {

	sqliteLocation := "../../chat.db"

	conn, err := sql.Open("sqlite", sqliteLocation)

	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	testQueries = New(conn)

	os.Exit(m.Run())
}
