package web

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/yfedoruck/book/pkg/cookie"
	"github.com/yfedoruck/book/pkg/crypto"
	"github.com/yfedoruck/book/pkg/midware"
	"github.com/yfedoruck/book/pkg/pg"
	"github.com/yfedoruck/book/pkg/user"
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
	formData := user.New(r.db)
	if err := c.Bind(&formData); err != nil {
		return err
	}

	log.Println("user:", formData)

	u := formData
	if u.Login() == nil {
		log.Println("error: user not found - ", formData)
		return c.Redirect(http.StatusFound, "/login")
	}

	err := crypto.Compare(u.Password, formData.Password)
	if err != nil {
		log.Println("error:", err, u, formData)
		return c.Redirect(http.StatusFound, "/login")
	}

	cookie.Write(c, map[string]string{
		"username": u.Username,
		"id":       strconv.Itoa(u.Id),
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

	u := user.New(r.db)
	if err := c.Bind(&u); err != nil {
		return err
	}
	if u.Username == "" {
		return c.Redirect(http.StatusFound, "/register")
	}
	cookie.Write(c, map[string]string{
		"username": u.Username,
	})

	fmt.Println("register user here")
	u.Password = crypto.Generate(c.FormValue("password"))
	//id := r.db.RegisterUser(username, crypto.Generate(password), email)
	id := u.Register()

	cookie.Write(c, map[string]string{
		"username": u.Username,
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
	Id       int
	Name     string `json:"username" form:"username" query:"username"`
	Email    string `json:"email" form:"email" query:"email"`
	Password string `json:"password" form:"password" query:"password"`
}
