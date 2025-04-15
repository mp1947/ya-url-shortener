package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/handler"
	"github.com/mp1947/ya-url-shortener/internal/middleware"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/repository/database"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"go.uber.org/zap"
)

func CreateRouter(
	c config.Config,
	s service.Service,
	repo repository.Repository,
	l *zap.Logger,
) *gin.Engine {

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(middleware.AuthMiddleware(l))
	r.Use(middleware.LoggerMiddleware(l))
	r.Use(middleware.GzipMiddleware())

	h := handler.HandlerService{Service: s, Cfg: c}

	r.Any("/", h.ShortenURL)
	r.Any("/:id", h.GetOriginalURLByID)

	if repo.GetType() == "database" {
		r.GET("/ping", h.Ping(repo.(*database.Database)))
	}

	api := r.Group("/api")
	api.POST("/shorten", h.JSONShortenURL)
	api.POST("/shorten/batch", h.BatchShortenURL)

	api.GET("/user/urls", h.GetUserURLS)
	api.DELETE("/user/urls", h.DeleteUserURLs)

	return r
}
