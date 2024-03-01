package server

import (
	"boletia/api/internal/obtain"
	"boletia/api/internal/platform/server/handler/currency"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"boletia/api/internal/platform/server/handler/health"
	"boletia/api/internal/platform/server/middleware/logging"
	"boletia/api/internal/platform/server/middleware/recovery"
	"github.com/gin-gonic/gin"
)

type Server struct {
	httpAddr        string
	engine          *gin.Engine
	getCurrency     obtain.CurrencyService
	shutdownTimeout time.Duration
}

func New(ctx context.Context, host string, port uint, shutdownTimeout time.Duration, getCurrency obtain.CurrencyService) (context.Context, Server) {
	//gin.SetMode(gin.ReleaseMode)
	srv := Server{
		engine:   gin.New(),
		httpAddr: fmt.Sprintf("%s:%d", host, port),

		shutdownTimeout: shutdownTimeout,
		getCurrency:     getCurrency,
	}

	srv.registerRoutes()
	return serverContext(ctx), srv
}

func (s *Server) registerRoutes() {
	s.engine.Use(recovery.Middleware(), logging.Middleware())

	s.engine.GET("/health", health.CheckHandler())
	s.engine.GET("/currencies/:currency", currency.ObtainHandler(s.getCurrency))
}

func (s *Server) Run(ctx context.Context) error {
	log.Println("Server running on", s.httpAddr)

	srv := &http.Server{
		Addr:    s.httpAddr,
		Handler: s.engine,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server shut down", err)
		}
	}()

	<-ctx.Done()
	ctxShutDown, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return srv.Shutdown(ctxShutDown)
}

func serverContext(ctx context.Context) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-c
		cancel()
	}()

	return ctx
}
