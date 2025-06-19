package handler_test

import (
	"context"
	"log"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/handler"
	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
	"github.com/mp1947/ya-url-shortener/internal/service"
)

const (
	testURL         = "https://console.yandex.cloud/"
	testJSONRequest = `{"url": "https://console.yandex.cloud"}`
)

var listenAddr = ":8080"
var baseURL = "http://localhost:8080"
var fileStoragePath = "./test.out"
var cfg = config.Config{
	ListenAddr:      &listenAddr,
	BaseURL:         &baseURL,
	FileStoragePath: &fileStoragePath,
}
var storage = &inmemory.Memory{}
var l, _ = logger.InitLogger()
var storageInitializedErr = storage.Init(context.Background(), cfg, l)
var hs = initTestHandlerService()

func initTestHandlerService() handler.HandlerService {

	if storageInitializedErr != nil {
		log.Fatalf("error initializing storage: %v", storageInitializedErr)
	}

	service := service.ShortenService{Storage: storage, Logger: l, Cfg: &cfg}

	return handler.HandlerService{Service: &service}
}
