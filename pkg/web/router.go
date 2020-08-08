package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yfedoruck/book/pkg/cookie"
	"github.com/yfedoruck/book/pkg/crypto"
	"github.com/yfedoruck/book/pkg/midware"
	"github.com/yfedoruck/book/pkg/pg"
	"log"
	"net/http"
	"strconv"
)

type Router struct {
	e  *echo.Echo
	db *pg.Postgres
}

func NewRouter(echo *echo.Echo, db *pg.Postgres) *Router {
	return &Router{
		e:  echo,
		db: db,
	}
}

func (r Router) Routes() {
	r.e.GET("/", RootHandler)
	r.e.GET("/books", BooksTplHandler)
	r.e.GET("/register", RegisterTplHandler)
	r.e.GET("/login", LoginTplHandler)
	r.e.GET("/logout", LogoutHandler)
	r.e.POST("/inner/register", r.RegisterHandler, midware.Register)
	r.e.POST("/inner/login", r.LoginHandler, midware.Login)
	r.e.Static("/static", "static")
	r.e.Static("/book", "data")
}

func (r Router) LoginHandler(c echo.Context) error {
	log.Println("call LoginHandler")
	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}

	log.Println("user:", u.Name, u.Password)
	id, encryptedPwd, err := r.db.LoginUser(u.Name)

	if err != nil {
		log.Println("error:", err)
		return c.Redirect(http.StatusFound, "/login")
	}

	err = crypto.Compare(encryptedPwd, u.Password)
	if err != nil {
		log.Println("error:", err)
		return c.Redirect(http.StatusFound, "/login")
	}

	cookie.Write(c, map[string]string{
		"username": u.Name,
		"id":       strconv.Itoa(id),
	})
	return c.Redirect(http.StatusFound, "/books")
}

func LoginTplHandler(c echo.Context) error {
	if cookie.Exists(c, "username") {
		return c.Redirect(http.StatusFound, "/books")
	}
	return c.Render(http.StatusOK, "login", map[string]interface{}{
		"Title": "Login",
		"Css":   "/static/css/style.css",
	})
}

func LogoutHandler(c echo.Context) error {
	if cookie.Exists(c, "username") {
		cookie.Delete(c, "username")
	}
	return c.Redirect(http.StatusFound, "/login")
}

func RootHandler(c echo.Context) error {
	if !cookie.Exists(c, "username") {
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

	if !cookie.Exists(c, "username") {
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

func (r Router) RegisterHandler(c echo.Context) error {
	log.Println("call RegisterHandler")
	if cookie.Exists(c, "username") {
		return c.Redirect(http.StatusFound, "/books")
	}

	u := new(User)
	if err := c.Bind(u); err != nil {
		return err
	}
	if u.Name == "" {
		return c.Redirect(http.StatusFound, "/register")
	}
	cookie.Write(c, map[string]string{
		"username": u.Name,
	})

	username := c.FormValue("username")
	password := c.FormValue("password")
	email := c.FormValue("email")

	fmt.Println("register user here")
	id := r.db.RegisterUser(username, crypto.Generate(password), email)

	cookie.Write(c, map[string]string{
		"username": u.Name,
		"id":       strconv.Itoa(id),
	})

	return c.Redirect(http.StatusFound, "/books")
	//return  c.String(http.StatusOK, name)
	//return  c.JSON(http.StatusOK, u)
}

func RegisterMiddleware() echo.MiddlewareFunc {
	return middleware.KeyAuth(func(username string, c echo.Context) (bool, error) {
		return username == "qwe", nil
	})
}

func LoginMiddleware() echo.MiddlewareFunc {
	//if cookie.Exists(c, "username") {
	//	return c.Redirect(http.StatusFound, "/books")
	//}
	return middleware.KeyAuth(func(username string, c echo.Context) (bool, error) {
		return username == "", nil
	})
}

func RegisterTplHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "register", map[string]interface{}{
		"Title": "Registration",
		"Css":   "/static/css/style.css",
	})
}

//TODO: Dangerous zone. Wrong `json:key` is hard to find.
type User struct {
	Name     string `json:"username" form:"username" query:"username"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}
