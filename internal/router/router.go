package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/handler"
)

func CreateRouter(c config.Config, urls handler.Urls) *gin.Engine {
	r := gin.Default()

	r.Any("/", urls.HandleOriginalURL(c))
	r.Any("/:id", urls.HandleShortURL)

	return r
}
