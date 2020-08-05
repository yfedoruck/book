package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type Router struct {
	e *echo.Echo
}

func NewRouter(echo *echo.Echo) *Router {
	return &Router{
		e: echo,
	}
}

func (r Router) Routes() {
	r.e.GET("/", RootHandler)
	r.e.GET("/books", BooksTplHandler)
	r.e.GET("/register", RegisterTplHandler)
	r.e.GET("/login", LoginTplHandler)
	r.e.POST("/inner/register", RegisterHandler)
	r.e.POST("/inner/login", LoginHandler)
	r.e.Static("/static", "static")
	r.e.Static("/book", "data")
}

func LoginHandler(c echo.Context) error {
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if u.Name == "" {
		return c.Redirect(http.StatusFound, "/register")
	}
	writeCookie(c, map[string]string{
		"username": u.Name,
	})
	return c.Redirect(http.StatusFound, "/books")
}

func LoginTplHandler(c echo.Context) error {
	if cookieExists(c, "username") {
		return c.Redirect(http.StatusFound, "/books")
	}
	return c.Render(http.StatusOK, "login", map[string]interface{}{
		"Title": "Login",
		"Css":   "/static/css/style.css",
	})
}

func RootHandler(c echo.Context) error {
	if !cookieExists(c, "username") {
		return c.Redirect(http.StatusFound, "/login")
	}

	return c.Render(http.StatusOK, "hello", map[string]interface{}{
		"LibTitle": "Books library",
		"Css":      "/static/css/style.css",
		"Books": []Book{
			{"Author 1", "Book 1", "book/123.txt"},
			{"Author 2", "Book 2", "book/456.txt"},
		},
	})
}

func BooksTplHandler(c echo.Context) error {

	if !cookieExists(c, "username") {
		return c.Redirect(http.StatusFound, "/login")
	}

	return c.Render(http.StatusOK, "hello", map[string]interface{}{
		"LibTitle": "Books library",
		"Css":      "/static/css/style.css",
		"Books": []Book{
			{"Author 1", "Book 1", "book/123.txt"},
			{"Author 2", "Book 2", "book/456.txt"},
		},
	})
}

func RegisterHandler(c echo.Context) error {
	if cookieExists(c, "username") {
		return c.Redirect(http.StatusFound, "/books")
	}

	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if u.Name == "" {
		return c.Redirect(http.StatusFound, "/register")
	}
	writeCookie(c, map[string]string{
		"username": u.Name,
	})

	return c.Redirect(http.StatusFound, "/books")
	//return  c.String(http.StatusOK, name)
	//return  c.JSON(http.StatusOK, u)
	//return  c.JSON(http.StatusOK, u)
}

func RegisterMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuth(func(username string, c echo.Context) (bool, error) {
		fmt.Println(username)
		return username == "qwe", nil
	})
}

func LoginMiddleware() echo.MiddlewareFunc {
	//if cookieExists(c, "username") {
	//	return c.Redirect(http.StatusFound, "/books")
	//}
	return middleware.KeyAuth(func(username string, c echo.Context) (bool, error) {
		fmt.Println(username)
		return username == "", nil
	})
}

func RegisterTplHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "register", map[string]interface{}{
		"Title": "Registration",
		"Css":   "/static/css/style.css",
		"Books": []Book{
			{"Author 1", "Book 1", "book/123.txt"},
			{"Author 2", "Book 2", "book/456.txt"},
		},
	})
}

type User struct {
	Name     string `json:"username" form:"username" query:"username"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"email" form:"email" query:"email"`
}
