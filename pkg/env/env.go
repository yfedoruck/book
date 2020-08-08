package env

import (
	"os"
)

// Set all path as absolute, because compiled program can be in other directory than source files
// and so runtime.Caller will get wrong path.
func BasePath() string {
	//_, b, _, ok := runtime.Caller(0)
	//if !ok {
	//	log.Panic("Caller error")
	//}
	//env := filepath.Dir(b)
	//pkg := filepath.Dir(env)
	//app := filepath.Dir(pkg)
	//return app

	return "/go/src/book"
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
		return "5000" //1323
	}
	return port
}
