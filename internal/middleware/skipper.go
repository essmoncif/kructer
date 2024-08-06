package middleware

import "github.com/labstack/echo/v4"

func SkipperJwtFn(skipURLs []string) func(echo.Context) bool {
	return func(context echo.Context) bool {
		for _, url := range skipURLs {
			if url == context.Request().URL.String() {
				return true
			}
		}
		return false
	}
}
