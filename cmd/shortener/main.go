package main

import (
	"net/http"

	"github.com/mp1947/ya-url-shortener/internal/app"
)

func main() {

	urls := &app.Urls{IdToUrl: map[string]string{}}

	mux := http.NewServeMux()
	mux.HandleFunc("/", urls.HandleOriginal)
	mux.HandleFunc("/{id}", urls.HandleShort)

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
