package helper

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key string = "ozaenzenzen"

func GenerateJWTToken() (string, error) {
	// Create a new token object
	token := jwt.New(jwt.SigningMethodHS256)

	// Set the claims for the token
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = "JohnDoe"
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // Token expires in 1 hour

	// Generate the token string
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func DecodeJWTToken(tokenString string) (string, error) {
	// return token.Raw, err
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		return "", nil
	}
	// ... error handling

	// // do something with decoded claims
	// for key, val := range claims {
	// 	fmt.Printf("Token: %v\n", token)
	// 	fmt.Printf("Key: %v, value: %v\n", key, val)
	// }

	return token.Claims.(jwt.MapClaims)["email"].(string), err
}

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verify the signing algorithm is HMAC with SHA-256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// Provide the secret key used for signing the token
		return []byte(key), nil
	})

	if err != nil {
		return false, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return true, nil
	}

	return false, nil
}
