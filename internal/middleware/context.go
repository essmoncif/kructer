package middleware

import (
	"github.com/labstack/echo/v4"
	"kructer.com/internal/context"
)

func KructerContextMiddleware(cc *context.KructerContext) echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc.Context = c
			return h(cc)
		}
	}
}