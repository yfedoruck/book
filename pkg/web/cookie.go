package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func writeCookie(c echo.Context, data map[string]string) {
	for name, value := range data {
		cookie := new(http.Cookie)
		cookie.Name = name
		cookie.Value = value
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)
	}
}

func cookieExists(c echo.Context, key string) bool {
	cookie, err := c.Cookie(key)
	if err != nil {
		return false
	}
	fmt.Println(cookie.Name)
	fmt.Println(cookie.Value)
	return true
}
