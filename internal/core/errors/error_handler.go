package errors

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func HTTPErrorHandler(err error, c echo.Context) {
	c.Logger().Error(err)
	code := http.StatusInternalServerError

	switch v := err.(type) {
	case *echo.HTTPError:
		err := c.JSON(v.Code, v)
		if err != nil {
			c.Logger().Error("error handler: json encoding", err)
		}
	default:
		e := NewBoom(InternalError, "Bad implementation", nil)
		err := c.JSON(code, e)
		if err != nil {
			c.Logger().Error("error handler: json encoding", err)
		}
	}
}
