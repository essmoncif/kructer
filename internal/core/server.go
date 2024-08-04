package core

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
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

func (s *Server) SetValidator(validor echo.Validator) {
	s.Server.Validator = validor
}

func (s *Server) Run() {
	s.Server.Start(s.config.Address)
}