package main

import (
	"fmt"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
	"github.com/mp1947/ya-url-shortener/internal/router"
)

func main() {

	storage := inmemory.InitStorage()

	cfg := config.Config{}
	cfg.ParseFlags()

	r := router.CreateRouter(cfg, storage)

	if err := r.Run(*cfg.ListenAddr); err != nil {
		fmt.Printf("error on server start: %v\n", err)
	}
}
