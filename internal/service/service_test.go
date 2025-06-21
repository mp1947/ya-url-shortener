package service_test

import (
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/eventlog"
	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/service"
)

var listenAddr = ":8080"
var baseURL = "http://localhost:8080"
var fileStoragePath = "./test.out"
var cfg = config.Config{
	ListenAddr:      &listenAddr,
	BaseURL:         &baseURL,
	FileStoragePath: &fileStoragePath,
}
var l, _ = logger.InitLogger()

func initTestService(r repository.Repository) *service.ShortenService {

	ep, _ := eventlog.NewEventProcessor(cfg)

	service := service.ShortenService{
		Storage: r,
		EP:      *ep,
		Logger:  l,
		Cfg:     &cfg,
		CommCh:  make(chan model.BatchDeleteShortURLs, 1),
	}

	return &service
}
