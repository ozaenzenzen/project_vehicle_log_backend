package helper

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

var key1 string = "ozaenzenzen"

// var key_platform1 string = "ozaenzenzen_plat"

func GenerateUserTokenV3(userstamp string) (string, *time.Time, string, *time.Time, error) {
	expAccessToken := time.Now().Add(time.Minute * 3)
	expRefreshToken := time.Now().Add(time.Hour * 168 * 2)
	//Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_stamp": userstamp,
		"exp":        expAccessToken.Unix(), // Token expires in 168 hour or 1 week
		// "exp":        time.Now().Add(time.Minute * 3).Unix(), // Token expires in 168 hour or 1 week
		// "exp":        time.Now().Add(time.Hour * 168).Unix(), // Token expires in 168 hour or 1 week
	})

	accessTokenString, errSignAccessToken := accessToken.SignedString([]byte(key1))
	if errSignAccessToken != nil {
		return "", nil, "", nil, errSignAccessToken
	}

	//Refresh Token
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_stamp": userstamp,
		"exp":        expRefreshToken.Unix(), // Token expires in 168 * 4 hour or 2 weeks
		// "exp":        time.Now().Add(time.Hour * 168 * 2).Unix(), // Token expires in 168 * 4 hour or 2 weeks
	})

	refreshTokenString, errSignRefreshToken := refreshToken.SignedString([]byte(key1))
	if errSignRefreshToken != nil {
		return "", nil, "", nil, errSignRefreshToken
	}

	return accessTokenString, &expAccessToken, refreshTokenString, &expRefreshToken, nil
}

func GenerateUserTokenV2(userstamp string) (string, string, error) {
	//Access Token
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_stamp": userstamp,
		"exp":        time.Now().Add(time.Minute * 3).Unix(), // Token expires in 168 hour or 1 week
		// "exp":        time.Now().Add(time.Hour * 168).Unix(), // Token expires in 168 hour or 1 week
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
