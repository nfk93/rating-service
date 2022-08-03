/*
 * Rating Service
 *
 * Optional multiline or single-line description in [CommonMark](http://commonmark.org/help/) or HTML.
 *
 * API version: 0.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nfk93/rating-service/generated/api"
	"github.com/nfk93/rating-service/generated/database"
	"github.com/nfk93/rating-service/internal/endpoints"
	"github.com/nfk93/rating-service/internal/logic/user"

	_ "github.com/lib/pq"
)

func main() {
	mustGetenv := func(k string) string {
		v := os.Getenv(k)
		if v == "" {
			log.Fatalf("Warning: %s environment variable not set.\n", k)
		}
		return v
	}

	var (
		dbUser         = mustGetenv("DB_USER")              // e.g. 'my-db-user'
		dbPwd          = mustGetenv("DB_PASS")              // e.g. 'my-db-password'
		unixSocketPath = mustGetenv("INSTANCE_UNIX_SOCKET") // e.g. '/cloudsql/project:region:instance'
		dbName         = mustGetenv("DB_NAME")              // e.g. 'my-database'
	)

	dbURI := fmt.Sprintf("user=%s password=%s database=%s host=%s",
		dbUser, dbPwd, dbName, unixSocketPath)

	sqldb, err := sql.Open("postgres", dbURI)
	if err != nil {
		panic(fmt.Sprintf("error opening db connection: %s", err))
	}

	queries := database.New(sqldb)
	userService := user.NewUserService(queries)

	DefaultApiService := endpoints.NewApiService(userService)
	DefaultApiController := api.NewDefaultApiController(DefaultApiService)

	router := api.NewRouter(DefaultApiController)
	log.Fatal(http.ListenAndServe(":8080", router))
}
