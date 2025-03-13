package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/handler"
	"github.com/mp1947/ya-url-shortener/internal/middleware"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"go.uber.org/zap"
)

func CreateRouter(
	c config.Config,
	s service.Service,
	l *zap.Logger,
) *gin.Engine {

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(l))

	h := handler.HandlerService{Service: s, Cfg: c}

	r.Any("/", h.ShortenURL)
	r.Any("/:id", h.GetOriginalURLByID)

	api := r.Group("/api")
	api.POST("/shorten", h.JSONShortenURL)

	return r
}
