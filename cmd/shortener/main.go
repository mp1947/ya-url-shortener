package main

import (
	"fmt"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/handler"
	"github.com/mp1947/ya-url-shortener/internal/router"
)

func main() {

	urls := &handler.Urls{ShortToOriginal: map[string]string{}}

	cfg := config.Config{}
	cfg.ParseFlags()

	r := router.CreateRouter(cfg, *urls)

	if err := r.Run(*cfg.ListenAddr); err != nil {
		fmt.Printf("error on server start: %v\n", err)
	}
}
