package context

import (
	"github.com/labstack/echo/v4"
	"kructer.com/config"
)

type KructerContext struct {
	echo.Context
	Config *config.Configuration
}
