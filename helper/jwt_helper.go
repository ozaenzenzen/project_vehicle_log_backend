package helper

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

var key string = "ozaenzenzen"

func GenerateJWTToken(uid string, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"uid":   uid,
		"email": email,
		// "exp":   time.Now().Add(time.Hour * 168).Unix(), // Token expires in 168 hour or 1 week
		"exp": time.Now().Add(time.Minute * 60).Unix(), // Token expires in 168 hour or 1 week
	})

	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTTokenOld() (string, error) {
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

func DecodeJWTToken(tokenString string) (jwt.MapClaims, error) {
	// return token.Raw, err
	hmacSecret := []byte(key)
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err != nil {
		return nil, nil
	}
	return token.Claims.(jwt.MapClaims), err
}

func DecodeJWTTokenForEmail(tokenString string) (string, error) {
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

func ValidateTokenJWT(c *gin.Context, db *gorm.DB, headertoken string) (bool, error) {
	if headertoken == "" {
		sampleErr := errors.New("token empty")
		return false, sampleErr
	}
	isValid, err := VerifyToken(headertoken)
	return isValid, err
}

func GetDataTokenJWT(headertoken string, isEmail bool) string {
	tokenRaw, err := DecodeJWTToken(headertoken)
	// fmt.Printf("\ntoken raw %v", tokenRaw)
	if err != nil {
		return ""
	}

	emails := tokenRaw["email"].(string)
	uid := tokenRaw["uid"].(string)

	if isEmail == true {
		return emails
	} else {
		return uid
	}
}
