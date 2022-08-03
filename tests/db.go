package tests

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	dbUser = "postgres"
	dbPwd  = "password"
	dbName = "test"
)

func SetupTestDB() *sql.DB {
	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=localhost sslmode=disable",
		dbUser, dbPwd, dbName)

	sqldb, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(fmt.Sprintf("error opening db connection: %s", err))
	}

	return sqldb
}
