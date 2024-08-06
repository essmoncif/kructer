package context

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"kructer.com/config"
	"kructer.com/internal/core"
)

type KructerContext struct {
	echo.Context
	Config  *config.Configuration
	DB      *gorm.DB
	Session *core.Session
}
