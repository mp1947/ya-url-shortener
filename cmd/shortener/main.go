// Package main initializes and starts the URL shortener web application,
// setting up configuration, logging, storage, services, and the HTTP server.
package main

import (
	"context"
	"log"

	"github.com/mp1947/ya-url-shortener/internal/shortener"
)

var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sh, err := shortener.InitShortener(
		ctx,
		buildVersion,
		buildDate,
		buildCommit,
	)

	if err != nil {
		log.Printf("error initializing shortener: %v", err)
		return
	}

	sh.Run()

}
