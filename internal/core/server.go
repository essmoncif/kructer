package core

import (
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
