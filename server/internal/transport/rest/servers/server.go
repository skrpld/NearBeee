package servers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/skrpld/NearBeee/internal/core/logger"
	"github.com/skrpld/NearBeee/internal/core/repository"
	"github.com/skrpld/NearBeee/internal/transport/rest/middlewares"
	"github.com/skrpld/NearBeee/internal/transport/rest/routers"
)

type HttpServerConfig struct {
	Host   string `env:"SERVER_HOST" env-default:"localhost" mapstructure:"SERVER_HOST"`
	Port   int    `env:"SERVER_PORT" env-default:"50050" mapstructure:"SERVER_PORT"`
	Secret string `env:"SECRET" env-default:"secret" mapstructure:"SECRET"`
}

type HttpServer struct {
	cfg    HttpServerConfig
	server *http.Server
	logger logger.Logger
}

func NewHttpServer(cfg HttpServerConfig, postgresRepo *repository.PostgresRepository, mongodbRepo *repository.MongodbRepository, logger logger.Logger) (*HttpServer, error) {
	mainMux := http.NewServeMux()

	authRouter, authSrv := routers.NewAuthRouter(postgresRepo, cfg.Secret)
	postsRouter := routers.NewPostsRouter(postgresRepo)
	messagesRouter := routers.NewMessagesRouter(mongodbRepo)

	authMiddleware := middlewares.NewAuthMiddlewareHandler(authSrv).AuthMiddleware

	apiMux := http.NewServeMux()
	apiMux.Handle("/auth/", authRouter)
	apiMux.Handle("/posts/", authMiddleware(postsRouter))
	apiMux.Handle("/messages/", authMiddleware(messagesRouter))

	handler := middlewares.LoggerMiddleware(logger)(
		middlewares.GlobalMiddleware(
			http.StripPrefix("/api", apiMux),
		),
	)

	mainMux.Handle("/api/", handler)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler: mainMux,
	}

	return &HttpServer{
		cfg:    cfg,
		server: server,
		logger: logger,
	}, nil
}

func (s *HttpServer) Start() error {
	s.logger.With(logger.Time("started_at", time.Now())).Info(fmt.Sprintf("Server started on %s:%d", s.cfg.Host, s.cfg.Port))
	return s.server.ListenAndServe()
}

func (s *HttpServer) Stop() error {
	s.logger.With(logger.Time("stopped_at", time.Now())).Info("Server stopped")
	return s.server.Shutdown(context.Background())
}
