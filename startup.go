package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"kructer.com/config"
	"kructer.com/internal/context"
	"kructer.com/internal/core"
	"kructer.com/internal/core/errors"
	midd "kructer.com/internal/middleware"
	"kructer.com/internal/routes"
)

func Bootstrap() {

	config, err := config.New()
	if err != nil {
		log.Fatalf("%+v\n", err)
	}

	server := core.NewServer(config)

	cc := context.KructerContext{
		Config: config,
		DB:     server.GetDB(),
	}

	server.SetHTTPErrorHandler(errors.HTTPErrorHandler)

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

	routes.InitRoutes(server.Server, server.GetDB())

	server.Server.HideBanner = true

	go func() {
		if err := server.Run(); err != nil {
			server.Server.Logger.Fatal("shutting down the server")
		}
	}()

	server.GracefulShutdown()
}
