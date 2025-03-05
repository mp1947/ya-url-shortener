package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/handler"
	"github.com/mp1947/ya-url-shortener/internal/repository/inmemory"
)

func CreateRouter(c config.Config, s *inmemory.Storage) *gin.Engine {
	r := gin.Default()

	r.Any("/", handler.ShortenURL(c, s))
	r.Any("/:id", handler.HandleShortURL(s))

	return r
}
