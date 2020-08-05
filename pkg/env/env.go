package env

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
)

func BasePath() string {
	_, b, _, ok := runtime.Caller(0)
	if !ok {
		log.Panic("Caller error")
	}
	env := filepath.Dir(b)
	pkg := filepath.Dir(env)
	app := filepath.Dir(pkg)
	return app
}

func Domain() string {
	domain := os.Getenv("USERDOMAIN")
	if domain == "home" {
		return "local"
	}
	return "heroku"
}

func Port() string {
	var port = os.Getenv("PORT")
	if port == "" {
		return "1323"
	}
	return port
}
