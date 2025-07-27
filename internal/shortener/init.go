package shortener

import (
	"context"
	"net/http"

	"github.com/mp1947/ya-url-shortener/config"
	"github.com/mp1947/ya-url-shortener/internal/interceptor"
	"github.com/mp1947/ya-url-shortener/internal/logger"
	"github.com/mp1947/ya-url-shortener/internal/model"
	"github.com/mp1947/ya-url-shortener/internal/repository"
	"github.com/mp1947/ya-url-shortener/internal/router"
	"github.com/mp1947/ya-url-shortener/internal/service"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// InitShortener initializes and configures the URL shortener application.
//
// It sets up the configuration, logger, storage repository, service layer, HTTP router, and optionally a gRPC server.
// The function logs build information and initialization steps. It returns a pointer to a Shortener instance
// containing all initialized components, or an error if any step fails.
//
// Parameters:
//   - ctx: Context for managing lifecycle and cancellation.
//   - buildVersion: The version string of the build.
//   - buildDate: The build date string.
//   - buildCommit: The commit hash of the build.
//
// Returns:
//   - *Shortener: The initialized Shortener application instance.
//   - error: An error if initialization fails.
func InitShortener(
	ctx context.Context,
	buildVersion string,
	buildDate string,
	buildCommit string,
) (*Shortener, error) {

	cfg := config.InitConfig()

	logger, err := logger.InitLogger()

	if err != nil {
		return nil, err
	}

	logger.Info("Build version", zap.String("version", buildVersion))
	logger.Info("Build date", zap.String("date", buildDate))
	logger.Info("Build commit", zap.String("commit", buildCommit))

	logger.Info(
		"initializing web application with config",
		zap.String("host", *cfg.HTTPServerAddress),
		zap.String("base_url", *cfg.BaseHTTPURL),
	)

	storage, err := repository.CreateRepository(logger, *cfg, ctx)
	if err != nil {
		return nil, err
	}

	logger.Info(
		"storage has been initialized",
		zap.String("type", storage.GetType()),
	)

	service := service.ShortenService{
		Cfg:     cfg,
		Logger:  logger,
		CommCh:  make(chan model.BatchDeleteShortURLs),
		Storage: storage,
	}

	r := router.CreateRouter(*cfg, &service, storage, logger)

	logger.Info(
		"router has been created. web server is ready to start",
	)

	srv := &http.Server{
		Addr:    *cfg.HTTPServerAddress,
		Handler: r.Handler(),
	}

	var grpcServer *grpc.Server

	if *cfg.GRPCEnabled {
		grpcServer = grpc.NewServer(grpc.UnaryInterceptor(interceptor.AuthUnaryInterceptor))
	}

	return &Shortener{
		repo:       storage,
		service:    service,
		cfg:        cfg,
		Logger:     logger,
		httpServer: srv,
		grpcServer: grpcServer,
	}, nil
}
