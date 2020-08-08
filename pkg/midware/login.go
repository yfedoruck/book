package midware

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

// Process is the middleware function.
func Login(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")

		if username == "" || password == "" {
			return c.Redirect(http.StatusFound, "/login")
		}

		if err := next(c); err != nil {
			c.Error(err)
		}

		return next(c)
	}
}

// Process is the middleware function.
func Register(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		email := c.FormValue("email")
		username := c.FormValue("username")
		password := c.FormValue("password")

		if len(username) < 3 || len(password) < 1 || len(email) < 3 {
			return c.Redirect(http.StatusFound, "/register")
		}

		if err := next(c); err != nil {
			c.Error(err)
		}

		return next(c)
	}
}
