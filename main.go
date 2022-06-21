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
	"log"
	"net/http"

	"github.com/nfk93/rating-service/generated/api"
	"github.com/nfk93/rating-service/internal/endpoints"
)

func main() {
	log.Printf("Server started")

	DefaultApiService := endpoints.NewDefaultApiService()
	DefaultApiController := api.NewDefaultApiController(DefaultApiService)

	router := api.NewRouter(DefaultApiController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
