package handlehttp_test

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/mp1947/ya-url-shortener/config"
	handler "github.com/mp1947/ya-url-shortener/internal/handler/http"
	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
	"github.com/mp1947/ya-url-shortener/internal/router"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"go.uber.org/zap"
)

const (
	testURL         = "https://console.yandex.cloud/"
	testJSONRequest = `{"url": "https://console.yandex.cloud"}`
)

var listenAddr = ":8080"
var baseURL = "http://localhost:8080"
var fileStoragePath = "./test.out"
var cfg = config.Config{
	HTTPServerAddress: &listenAddr,
	BaseHTTPURL:       &baseURL,
	FileStoragePath:   &fileStoragePath,
}
var storage = &inmemory.Memory{}
var l, _ = logger.InitLogger()
var storageInitializedErr = storage.Init(context.Background(), cfg, l)
var hs = initTestHandlerService()

func initTestHandlerService() handler.HandlerService {

	if storageInitializedErr != nil {
		log.Fatalf("error initializing storage: %v", storageInitializedErr)
	}

	service := service.ShortenService{
		Storage: storage,
		Logger:  l,
		Cfg:     &cfg,
		CommCh:  make(chan model.BatchDeleteShortURLs, 1),
	}

	return handler.HandlerService{Service: &service}
}

func setupTestServer() (string, func()) {
	router := router.CreateRouter(cfg, hs.Service, storage, l)
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		l.Fatal("failed to start test server", zap.Error(err))
	}
	srv := &http.Server{Handler: router}

	go func() {
		if err := srv.Serve(listener); err != nil && err != http.ErrServerClosed {
			l.Fatal("server error", zap.Error(err))
		}
	}()

	return "http://" + listener.Addr().String(), func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		_ = srv.Shutdown(ctx)
	}
}
