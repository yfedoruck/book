package crypto

import (
	"github.com/yfedoruck/book/pkg/fail"
	"golang.org/x/crypto/bcrypt"
)

func Generate(password string) string {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fail.Check(err)
	return string(fromPassword)
}

func Compare(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
