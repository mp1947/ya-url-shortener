package router

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/handler"
	im "github.com/mp1947/ya-url-shortener/internal/middleware"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/repository/database"
	"github.com/mp1947/ya-url-shortener/internal/service"
	pm "github.com/mp1947/ya-url-shortener/pkg/middleware"
	"go.uber.org/zap"
)

// CreateRouter initializes and configures a new Gin router with the provided configuration, service, repository, and logger.
// It sets up middleware for recovery, authentication, logging, and gzip compression.
// The function registers HTTP handlers for URL shortening, retrieval, batch operations, user-specific endpoints, and health checks.
// If the repository type is "database", a /ping endpoint is added for database connectivity checks.
// The function also registers pprof endpoints for profiling and debugging.
// Returns the configured *gin.Engine instance.
func CreateRouter(
	c config.Config,
	s service.Service,
	repo repository.Repository,
	l *zap.Logger,
) *gin.Engine {

	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(im.AuthMiddleware(l))
	r.Use(pm.LoggerMiddleware(l))
	r.Use(pm.GzipMiddleware())

	h := handler.HandlerService{Service: s}

	r.Any("/", h.ShortenURL)
	r.Any("/:id", h.GetOriginalURLByID)

	if repo.GetType() == "database" {
		r.GET("/ping", h.Ping(repo.(*database.Database)))
	}

	api := r.Group("/api")
	api.POST("/shorten", h.JSONShortenURL)
	api.POST("/shorten/batch", h.BatchShortenURL)

	api.GET("/user/urls", h.GetUserURLs)
	api.DELETE("/user/urls", h.DeleteUserURLs)

	pprof.Register(r, "debug/pprof")

	return r
}
