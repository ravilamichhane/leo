package middlewares

import (
	"fmt"
	"strings"

	"github.com/ravilmc/leo/jwt"
	"github.com/ravilmc/leo/types"
	"github.com/ravilmc/leo/web"

	"github.com/labstack/echo/v4"
)

func JwtMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")
		if header == "" {
			return next(c)
		}
		token := strings.Replace(header, "Bearer ", "", 1)

		claims, err := jwt.VerifyJWT(token)

		if err != nil {
			if err == types.ErrTokenExpired {
				return web.NewTrustedError(fmt.Errorf("token expired"), 401)
			}

			return web.NewTrustedError(err, 401)
		}

		web.SetActiveUser(c, claims)

		return next(c)
	}
}

func UserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := web.GetActiveUser(c)
		if user == nil {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		return next(c)
	}
}

func AdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		activeUser := web.GetActiveUser(c)
		if activeUser == nil {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		if web.IsAdmin(activeUser) {
			return next(c)
		}

		return echo.NewHTTPError(403, "Forbidden")
	}
}

func SuperAdminMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		activeUser := web.GetActiveUser(c)
		if activeUser == nil {
			return echo.NewHTTPError(401, "Unauthorized")
		}

		if web.IsSuperAdmin(activeUser) {
			return next(c)
		}
		return echo.NewHTTPError(403, "Forbidden")
	}
}
