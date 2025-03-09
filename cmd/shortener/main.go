package main

import (
	"fmt"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
	"github.com/mp1947/ya-url-shortener/internal/router"
	"github.com/mp1947/ya-url-shortener/internal/service"
)

func main() {

	storage := &inmemory.Memory{}
	storage.Init()

	service := service.ShortenService{Storage: storage}

	cfg := config.Config{}
	cfg.ParseFlags()

	r := router.CreateRouter(cfg, service)

	if err := r.Run(*cfg.ListenAddr); err != nil {
		fmt.Printf("error on server start: %v\n", err)
	}
}
