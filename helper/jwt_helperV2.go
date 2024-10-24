package helper

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key1 string = "ozaenzenzen"

// var key_platform1 string = "ozaenzenzen_plat"

func GenerateUserTokenV2(userstamp string) (string, string, error) {
	//Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_stamp": userstamp,
		"exp":        time.Now().Add(time.Hour * 168).Unix(), // Token expires in 168 hour or 1 week
	})

	accessTokenString, err := accessToken.SignedString([]byte(key1))
	if err != nil {
		return "", "", err
	}

	//Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_stamp": userstamp,
		"exp":        time.Now().Add(time.Hour * 168 * 2).Unix(), // Token expires in 168 * 4 hour or 2 weeks
	})

	refreshTokenString, err := refreshToken.SignedString([]byte(key1))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}
