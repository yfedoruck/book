package cookie

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func Write(c echo.Context, data map[string]string) {
	for name, value := range data {
		cookie := new(http.Cookie)
		cookie.Name = name
		cookie.Value = value
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)
	}
}

func Exists(c echo.Context, key string) bool {
	_, err := c.Cookie(key)
	if err != nil {
		return false
	}
	return true
}

func Delete(c echo.Context, key string) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.SetCookie(cookie)
}
