package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/handler"
	"github.com/mp1947/ya-url-shortener/internal/service"
)

func CreateRouter(c config.Config, s service.Service) *gin.Engine {
	r := gin.Default()

	h := handler.HandlerService{Service: s, Cfg: c}

	r.Any("/", h.ShortenURL)
	r.Any("/:id", h.GetOriginalURLByID)

	return r
}
