package routes

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

func InitRoutes(e *echo.Echo, db *gorm.DB) {
	InitUserRoutes(e, db)
}