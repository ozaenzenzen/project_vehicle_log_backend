package helper

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func GenerateJWTToken() (string, error) {
	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "JohnDoe"
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Token expires in 1 hour

	// Generate the token string
	tokenString, err := token.SignedString([]byte("ozaenzenzen"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing algorithm is HMAC with SHA-256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the secret key used for signing the token
		return []byte("ozaenzenzen"), nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, nil
}
