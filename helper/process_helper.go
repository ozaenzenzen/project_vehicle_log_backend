package helper

import (
	"fmt"
	"net/http"
	baseResp "project_vehicle_log_backend/data"
	account "project_vehicle_log_backend/models/account"

	// platform "project_vehicle_log_backend/models/platform"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/jinzhu/gorm"
)

func CustomValidatorWithRefreshTokenOld(c *gin.Context, isRefreshToken bool) (*gorm.DB, *string, *account.AccountUserModel, *baseResp.BaseResponseModel) {
	var header_token string
	var header_refreshtoken string

	header_token = c.Request.Header.Get("token")
	if header_token == "" {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: "invalid credential3",
			Data:    nil,
		}
	}
	fmt.Println("header_token", header_token)

	isValidToken, errVerifyToken := VerifyUserToken(header_token)
	if errVerifyToken != nil || !isValidToken {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusUnauthorized,
			Message: "expired",
			Data:    nil,
		}
	}

	if isRefreshToken {
		header_refreshtoken = c.Request.Header.Get("refreshToken")
		if header_refreshtoken == "" {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusBadRequest,
				Message: "invalid credential5",
				Data:    nil,
			}
		}

		isValidRefreshToken, responseMessage, errVerifyRefreshToken := VerifyRefreshToken(header_refreshtoken)
		if errVerifyRefreshToken != nil || !isValidRefreshToken {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusUnauthorized,
				Message: responseMessage,
				Data:    nil,
			}
		}
	}

	var tokenRaw jwt.MapClaims
	var refreshTokenRaw jwt.MapClaims
	var errDecodeUserToken error
	var errDecodeRefreshToken error

	// tokenRaw, errDecodeUserToken = DecodeUserTokenWithoutSignature(header_token)
	tokenRaw, errDecodeUserToken = DecodeUserToken(header_token)
	if errDecodeUserToken != nil {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Failed Decode",
			Data:    nil,
		}
	}
	fmt.Println("tokenRaw1", tokenRaw)

	if isRefreshToken {
		refreshTokenRaw, errDecodeRefreshToken = DecodeRefreshToken(header_refreshtoken)
		if errDecodeRefreshToken != nil {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusBadRequest,
				Message: "Failed Decode2",
				Data:    nil,
			}
		}
		fmt.Println("refreshTokenRaw", refreshTokenRaw)
	}

	var userStamp string
	if resultUserStamp, ok := tokenRaw["user_stamp"].(string); ok {
		userStamp = resultUserStamp
	} else {
		fmt.Println("tokenRaw2", tokenRaw["user_stamp"])
		userStamp = ""
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "failed parsing1",
			Data:    nil,
		}
	}

	var userTokenExpires *float64
	if resultTokenExpires, ok := tokenRaw["exp"].(float64); ok {
		userTokenExpires = &resultTokenExpires
	} else {
		userTokenExpires = nil
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "failed parsing2",
			Data:    nil,
		}
	}

	// Convert the Unix timestamp to a time.Time
	parsedTime := time.Unix(int64(*userTokenExpires), 0)
	if parsedTime.Before(time.Now()) {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusUnauthorized,
			Message: "access token expired",
			Data:    nil,
		}
	}

	if isRefreshToken {
		var refreshTokenExpires *float64
		if resultTokenExpires, ok := refreshTokenRaw["exp"].(float64); ok {
			refreshTokenExpires = &resultTokenExpires
		} else {
			refreshTokenExpires = nil
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusInternalServerError,
				Message: "failed parsing3",
				Data:    nil,
			}
		}

		// Convert the Unix timestamp to a time.Time
		refreshTokenParsedTime := time.Unix(int64(*refreshTokenExpires), 0)
		if refreshTokenParsedTime.Before(time.Now()) {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusUnauthorized,
				Message: "refresh token expired",
				Data:    nil,
			}
		}
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusInternalServerError,
			Message: db.Error.Error(),
			Data:    nil,
		}
	}

	//--------check id--------check id--------check id--------
	var dataAccount account.AccountUserModel
	if !isRefreshToken {
		if errDataAccount := db.Table("account_user_models").
			Where("user_stamp = ?", userStamp).
			First(&dataAccount).Error; errDataAccount != nil {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusBadRequest,
				Message: "User Data Not Found",
				Data:    nil,
			}
		}
	} else {
		if errDataAccount := db.Table("account_user_models").
			Where("user_stamp = ?", userStamp).
			Where("refresh_token = ?", header_refreshtoken).
			First(&dataAccount).Error; errDataAccount != nil {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusUnauthorized,
				Message: "Expired",
				Data:    nil,
			}
		}
	}

	return db, &userStamp, &dataAccount, nil
}

func CustomValidatorWithRefreshToken(c *gin.Context, isRefreshToken bool) (*gorm.DB, *string, *account.AccountUserModel, *baseResp.BaseResponseModel) {
	var header_token string
	var header_refreshtoken string

	header_token = c.Request.Header.Get("token")
	if header_token == "" {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: "invalid credential3",
			Data:    nil,
		}
	}
	fmt.Println("header_token", header_token)

	// isValidToken, errVerifyToken := VerifyUserToken(header_token)
	// if errVerifyToken != nil || !isValidToken {
	// 	return nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusUnauthorized,
	// 		Message: "expired",
	// 		Data:    nil,
	// 	}
	// }

	if isRefreshToken {
		header_refreshtoken = c.Request.Header.Get("refreshToken")
		if header_refreshtoken == "" {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusBadRequest,
				Message: "invalid credential5",
				Data:    nil,
			}
		}

		isValidRefreshToken, responseMessage, errVerifyRefreshToken := VerifyRefreshToken(header_refreshtoken)
		if errVerifyRefreshToken != nil || !isValidRefreshToken {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusUnauthorized,
				Message: responseMessage,
				Data:    nil,
			}
		}
	}

	var tokenRaw jwt.MapClaims
	var refreshTokenRaw jwt.MapClaims
	var errDecodeUserToken error
	var errDecodeRefreshToken error

	tokenRaw, errDecodeUserToken = DecodeUserTokenWithoutSignature(header_token)
	// tokenRaw, errDecodeUserToken = DecodeUserToken(header_token)
	if errDecodeUserToken != nil {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Failed Decode",
			Data:    nil,
		}
	}
	fmt.Println("tokenRaw1", tokenRaw)

	if isRefreshToken {
		refreshTokenRaw, errDecodeRefreshToken = DecodeRefreshToken(header_refreshtoken)
		if errDecodeRefreshToken != nil {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusBadRequest,
				Message: "Failed Decode2",
				Data:    nil,
			}
		}
		fmt.Println("refreshTokenRaw", refreshTokenRaw)
	}

	var userStamp string
	if resultUserStamp, ok := tokenRaw["user_stamp"].(string); ok {
		userStamp = resultUserStamp
	} else {
		fmt.Println("tokenRaw2", tokenRaw["user_stamp"])
		userStamp = ""
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "failed parsing1",
			Data:    nil,
		}
	}

	// var userTokenExpires *float64
	// if resultTokenExpires, ok := tokenRaw["exp"].(float64); ok {
	// 	userTokenExpires = &resultTokenExpires
	// } else {
	// 	userTokenExpires = nil
	// 	return nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusInternalServerError,
	// 		Message: "failed parsing2",
	// 		Data:    nil,
	// 	}
	// }

	// // Convert the Unix timestamp to a time.Time
	// parsedTime := time.Unix(int64(*userTokenExpires), 0)
	// if parsedTime.Before(time.Now()) {
	// 	return nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusUnauthorized,
	// 		Message: "access token expired",
	// 		Data:    nil,
	// 	}
	// }

	if isRefreshToken {
		var refreshTokenExpires *float64
		if resultTokenExpires, ok := refreshTokenRaw["exp"].(float64); ok {
			refreshTokenExpires = &resultTokenExpires
		} else {
			refreshTokenExpires = nil
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusInternalServerError,
				Message: "failed parsing3",
				Data:    nil,
			}
		}

		// Convert the Unix timestamp to a time.Time
		refreshTokenParsedTime := time.Unix(int64(*refreshTokenExpires), 0)
		if refreshTokenParsedTime.Before(time.Now()) {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusUnauthorized,
				Message: "refresh token expired",
				Data:    nil,
			}
		}
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusInternalServerError,
			Message: db.Error.Error(),
			Data:    nil,
		}
	}

	//--------check id--------check id--------check id--------
	var dataAccount account.AccountUserModel
	if !isRefreshToken {
		if errDataAccount := db.Table("account_user_models").
			Where("user_stamp = ?", userStamp).
			First(&dataAccount).Error; errDataAccount != nil {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusBadRequest,
				Message: "User Data Not Found",
				Data:    nil,
			}
		}
	} else {
		if errDataAccount := db.Table("account_user_models").
			Where("user_stamp = ?", userStamp).
			Where("refresh_token = ?", header_refreshtoken).
			First(&dataAccount).Error; errDataAccount != nil {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusUnauthorized,
				Message: "Expired",
				Data:    nil,
			}
		}
	}

	return db, &userStamp, &dataAccount, nil
}

func CustomValidatorAC(c *gin.Context) (*gorm.DB, *string, *account.AccountUserModel, *baseResp.BaseResponseModel) {
	header_token := c.Request.Header.Get("token")
	if header_token == "" {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: "invalid credential3",
			Data:    nil,
		}
	}

	isValidToken, errVerifyToken := VerifyUserToken(header_token)
	if errVerifyToken != nil || !isValidToken {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusUnauthorized,
			Message: "expired",
			Data:    nil,
		}
	}

	tokenRaw, err := DecodeUserToken(header_token)
	if err != nil {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}

	var userStamp string
	if resultUserStamp, ok := tokenRaw["user_stamp"].(string); ok {
		userStamp = resultUserStamp
	} else {
		userStamp = ""
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "failed parsing",
			Data:    nil,
		}
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusInternalServerError,
			Message: db.Error.Error(),
			Data:    nil,
		}
	}

	var dataAccount account.AccountUserModel
	if err := db.Table("account_user_models").
		Where("user_stamp = ?", userStamp).
		First(&dataAccount).Error; err != nil {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: "User Data Not Found",
			Data:    nil,
		}
	}

	if dataAccount.IsActivated == 0 {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusNotFound,
			Message: "Account Unverified",
			Data:    nil,
		}
	}

	if dataAccount.IsActivated == 2 {
		return nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusNotFound,
			Message: "Account Deleted",
			Data:    nil,
		}
	}

	return db, &userStamp, &dataAccount, nil
}

func CustomValidatorWithRequestBody[T any](c *gin.Context, requestModel T) (*T, *gorm.DB, *string, *account.AccountUserModel, *baseResp.BaseResponseModel) {
	if err := c.ShouldBindJSON(&requestModel); err != nil {
		return nil, nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Data tidak lengkap1",
			Data:    nil,
		}
	}

	validate := validator.New()
	if err := validate.Struct(requestModel); err != nil {
		return nil, nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Data tidak lengkap2",
			Data:    nil,
		}
	}

	header_token := c.Request.Header.Get("token")
	if header_token == "" {
		return nil, nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: "invalid credential3",
			Data:    nil,
		}
	}

	isValidToken, errVerifyToken := VerifyUserToken(header_token)
	if errVerifyToken != nil || !isValidToken {
		return nil, nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusUnauthorized,
			Message: "expired",
			Data:    nil,
		}
	}

	tokenRaw, err := DecodeUserToken(header_token)
	if err != nil {
		return nil, nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
			Data:    nil,
		}
	}

	var userStamp string
	if resultUserStamp, ok := tokenRaw["user_stamp"].(string); ok {
		userStamp = resultUserStamp
	} else {
		userStamp = ""
		return nil, nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "failed parsing",
			Data:    nil,
		}
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		return nil, nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusInternalServerError,
			Message: db.Error.Error(),
			Data:    nil,
		}
	}

	//--------check id--------check id--------check id--------
	var dataAccount account.AccountUserModel
	if err := db.Table("account_user_models").
		Where("user_stamp = ?", userStamp).
		First(&dataAccount).Error; err != nil {
		return nil, nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusBadRequest,
			Message: "User Data Not Found",
			Data:    nil,
		}
	}

	if dataAccount.IsActivated == 0 {
		return nil, nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusNotFound,
			Message: "Account Unverified",
			Data:    nil,
		}
	}

	if dataAccount.IsActivated == 2 {
		return nil, nil, nil, nil, &baseResp.BaseResponseModel{
			Status:  http.StatusNotFound,
			Message: "Account Deleted",
			Data:    nil,
		}
	}

	return &requestModel, db, &userStamp, &dataAccount, nil
}

// func HeaderPlatformValidator(c *gin.Context) (*string, *baseResp.BaseResponseModel) {
// 	header_platformkey := c.Request.Header.Get("platformkey")
// 	if header_platformkey == "" {
// 		return nil, &baseResp.BaseResponseModel{
// 			Status:  http.StatusBadRequest,
// 			Message: "invalid credential1",
// 			Data:    nil,
// 		}
// 	}

// 	isValidPlatformKey, errVerifyPlatformKey := VerifyPlatformToken(header_platformkey)
// 	if errVerifyPlatformKey != nil || !isValidPlatformKey {
// 		return nil, &baseResp.BaseResponseModel{
// 			Status:  http.StatusBadRequest,
// 			Message: "invalid credential2",
// 			Data:    nil,
// 		}
// 	}

// 	platformName, errorResp := GetPlatformNameFromHeader(c, header_platformkey)
// 	if errorResp != nil {
// 		return nil, &baseResp.BaseResponseModel{
// 			Status:  http.StatusBadRequest,
// 			Message: errorResp.Message,
// 			Data:    nil,
// 		}
// 	}

// 	return platformName, nil
// }
