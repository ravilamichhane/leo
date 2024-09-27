package leo

import "github.com/labstack/echo/v4"

type Route struct {
	Path        string
	Method      string
	Middlewares []echo.MiddlewareFunc
	Handler     echo.HandlerFunc
}
