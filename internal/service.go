package internal

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/handlers"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jcserv/go-api-template/internal/repository"
	"github.com/jcserv/go-api-template/internal/service"
	_http "github.com/jcserv/go-api-template/internal/transport/http"
	v1 "github.com/jcserv/go-api-template/internal/transport/http/v1"
	"github.com/jcserv/go-api-template/internal/utils/log"
)

type Service struct {
	api *_http.API
	cfg *Configuration
}

func NewService() (*Service, error) {
	cfg, err := NewConfiguration()
	if err != nil {
		return nil, err
	}

	s := &Service{cfg: cfg}

	conn, err := s.ConnectDB(context.Background())
	if err != nil {
		return nil, err
	}
	repo := repository.New(conn)

	s.api = _http.NewAPI(v1.NewDependencies(service.NewBookService(repo)))
	return s, nil
}

func (s *Service) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()
		s.StartHTTP(ctx)
	}(ctx)

	wg.Wait()
	return nil
}

func (s *Service) StartHTTP(ctx context.Context) error {
	log.Info(ctx, fmt.Sprintf("Starting HTTP server on port %s", s.cfg.HTTPPort))
	r := s.api.RegisterRoutes()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	origins := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methods := handlers.AllowedMethods([]string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions})
	exposedHeaders := handlers.ExposedHeaders([]string{"X-Requested-With", "Content-Type", "X-Total-Count", "X-Token"})

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.cfg.HTTPPort),
		Handler: handlers.CORS(origins, headers, methods, exposedHeaders)(r),
	}

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Error(ctx, fmt.Sprintf("Error shutting down HTTP server: %v", err))
		}
	}()

	srv.ListenAndServe()
	return nil
}

func (s *Service) ConnectDB(ctx context.Context) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(s.cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 25
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 30 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}
	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}
	log.Info(ctx, "Connection established to database")

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	return pool, nil
}
