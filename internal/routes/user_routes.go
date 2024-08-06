package routes

import (
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"kructer.com/internal/application/user"
)

func InitUserRoutes(e *echo.Echo, db *gorm.DB) {
	userCtrl := user.NewUserController(db)

	userGroup := e.Group("/api/users")
	userGroup.POST("", userCtrl.CreateUser)
	userGroup.POST("/login", userCtrl.Login)
}