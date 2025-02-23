package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/app"
)

func setupRouter(c config.Config, urls app.Urls) *gin.Engine {
	r := gin.Default()

	r.Any("/", urls.HandleOriginal(c))
	r.Any("/:id", urls.HandleShort)

	return r
}

func main() {

	urls := &app.Urls{IDToURL: map[string]string{}}

	cfg := config.Config{}
	cfg.ParseFlags()

	r := setupRouter(cfg, *urls)

	if err := r.Run(*cfg.ListenAddr); err != nil {
		fmt.Printf("error on server start: %v\n", err)
	}
}
