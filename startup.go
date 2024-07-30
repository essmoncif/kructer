package main

import (
	"log"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	midd "kructer.com/internal/middleware"
)

func Run() {
	e := echo.New()

	// Middleware setup
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(midd.KructerContextMiddleware())


	port := ":8080"
	log.Fatal(e.Start(port))
}
