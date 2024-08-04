package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"kructer.com/config"
	"kructer.com/internal/context"
	"kructer.com/internal/core"
	midd "kructer.com/internal/middleware"
)

func Bootstrap() {

	config, err := config.New()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	server := core.NewServer(config)

	cc := context.KructerContext{
		Config: config,
	}

	server.AddMiddleware(middleware.Logger())
	server.AddMiddleware(middleware.Recover())
	server.AddMiddleware(middleware.BodyLimit("5M"))
	server.AddMiddleware(middleware.Secure())
	server.AddMiddleware(middleware.RequestID())
	server.AddMiddleware(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	server.AddMiddleware(midd.KructerContextMiddleware(&cc))

	server.Run()
}
