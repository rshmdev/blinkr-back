package hooks

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("secretpassphrase")

func GenerateJWTToken(email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expira em 24 horas

	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	return tokenString, err
}

func ComparePasswords(hashedPassword string, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password)
	return err == nil
}
