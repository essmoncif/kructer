package core

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"gopkg.in/go-playground/validator.v9"
	"kructer.com/config"
)

type Server struct {
	Server   *echo.Echo
	config   *config.Configuration
	db       *gorm.DB
	registry *ModelRegistry
}

func NewServer(config *config.Configuration) *Server {
	server := &Server{}
	server.config = config
	server.registry = NewRegistry()
	if err := server.registry.OpenWithConfig(config); err != nil {
		log.Fatalf("gorm: could not connect to db %q", err)
	}

	server.db = server.registry.DB
	server.Server = echo.New()
	server.Server.Validator = &Validator{validator: validator.New()}

	return server
}

func (s *Server) GetConfig() *config.Configuration {
	return s.config
}

func (s *Server) GetRegistry() *ModelRegistry {
	return s.registry
}

func (s *Server) GetDB() *gorm.DB {
	return s.db
}

func (s *Server) AddMiddleware(middleware echo.MiddlewareFunc) {
	s.Server.Use(middleware)
}

func (s *Server) SetHTTPErrorHandler(handler echo.HTTPErrorHandler) {
	s.Server.HTTPErrorHandler = handler
}

func (s *Server) Run() error {
	return s.Server.Start(s.config.Address)
}

func (s *Server) GracefulShutdown() {
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// close database connection
	if s.db != nil {
		dErr := s.db.Close()
		if dErr != nil {
			s.Server.Logger.Fatal(dErr)
		}
	}

	// shutdown http server
	if err := s.Server.Shutdown(ctx); err != nil {
		s.Server.Logger.Fatal(err)
	}
}
