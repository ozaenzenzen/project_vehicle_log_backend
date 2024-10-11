package helper

import (
	"net/http"
	baseResp "project_vehicle_log_backend/data"
	account "project_vehicle_log_backend/models/account"

	// platform "project_vehicle_log_backend/models/platform"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// func GetPlatformNameFromHeader(c *gin.Context, header_platformkey string) (*string, *baseResp.BaseResponseModel) {
// 	tokenPlatformRaw, err := DecodePlatformToken(header_platformkey)
// 	if err != nil {
// 		return nil, &baseResp.BaseResponseModel{
// 			Status:  http.StatusBadRequest,
// 			Message: err.Error(),
// 			Data:    nil,
// 		}
// 	}

// 	var platformName string
// 	if resultPlatformName, ok := tokenPlatformRaw["platform_name"].(string); ok {
// 		platformName = resultPlatformName
// 	} else {
// 		platformName = ""
// 		return nil, &baseResp.BaseResponseModel{
// 			Status:  http.StatusInternalServerError,
// 			Message: "Error PN",
// 			Data:    nil,
// 		}
// 	}

//		if platformName == "" {
//			return nil, &baseResp.BaseResponseModel{
//				Status:  http.StatusInternalServerError,
//				Message: "Error PN2",
//				Data:    nil,
//			}
//		}
//		return &platformName, nil
//	}
func CustomValidatorWithRefreshToken(c *gin.Context, isRefreshToken bool) (*gorm.DB, *string, *account.AccountUserModel, *baseResp.BaseResponseModel) {
	// header_platformkey := c.Request.Header.Get("platformkey")
	// if header_platformkey == "" {
	// 	return nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "invalid credential1",
	// 		Data:    nil,
	// 	}
	// }

	// isValidPlatformKey, errVerifyPlatformKey := VerifyPlatformToken(header_platformkey)
	// if errVerifyPlatformKey != nil || !isValidPlatformKey {
	// 	return nil, nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "invalid credential2",
	// 		Data:    nil,
	// 	}
	// }

	var header_token string
	var header_refreshtoken string
	if !isRefreshToken {
		header_token = c.Request.Header.Get("token")
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
				Status:  http.StatusBadRequest,
				Message: "invalid credential4",
				Data:    nil,
			}
		}
	} else {
		header_refreshtoken = c.Request.Header.Get("refreshToken")
		if header_refreshtoken == "" {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusBadRequest,
				Message: "invalid credential5",
				Data:    nil,
			}
		}

		isValidToken, errVerifyToken := VerifyUserToken(header_refreshtoken)
		if errVerifyToken != nil || !isValidToken {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusBadRequest,
				Message: "invalid credential6",
				Data:    nil,
			}
		}
	}

	// platformName, errorResp := GetPlatformNameFromHeader(c, header_platformkey)
	// if errorResp != nil {
	// 	return nil, nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusBadRequest,
	// 		Message: errorResp.Message,
	// 		Data:    nil,
	// 	}
	// }

	var tokenRaw jwt.MapClaims
	var errDecodeUserToken error
	if !isRefreshToken {
		tokenRaw, errDecodeUserToken = DecodeUserToken(header_token)
		if errDecodeUserToken != nil {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusBadRequest,
				Message: errDecodeUserToken.Error(),
				Data:    nil,
			}
		}
	} else {
		tokenRaw, errDecodeUserToken = DecodeUserToken(header_refreshtoken)
		if errDecodeUserToken != nil {
			return nil, nil, nil, &baseResp.BaseResponseModel{
				Status:  http.StatusBadRequest,
				Message: errDecodeUserToken.Error(),
				Data:    nil,
			}
		}
	}

	var userStamp string
	if resultUserStamp, ok := tokenRaw["user_stamp"].(string); ok {
		userStamp = resultUserStamp
	} else {
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
			Status:  http.StatusInternalServerError,
			Message: "refresh token expired",
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

	//--------check platform--------check platform--------check platform--------
	// var dataPlatform platform.PlatformModel
	// if err := db.Table("platform_models").
	// 	Where("platform_name = ?", platformName).
	// 	First(&dataPlatform).Error; err != nil {
	// 	return nil, nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "User Data Not Found",
	// 		Data:    nil,
	// 	}
	// }

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
				Status:  http.StatusBadRequest,
				Message: "Expired",
				Data:    nil,
			}
		}
	}

	return db, &userStamp, &dataAccount, nil
}

func CustomValidatorAC(c *gin.Context) (*gorm.DB, *string, *account.AccountUserModel, *baseResp.BaseResponseModel) {
	// header_platformkey := c.Request.Header.Get("platformkey")
	// if header_platformkey == "" {
	// 	return nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "invalid credential1",
	// 		Data:    nil,
	// 	}
	// }

	// isValidPlatformKey, errVerifyPlatformKey := VerifyPlatformToken(header_platformkey)
	// if errVerifyPlatformKey != nil || !isValidPlatformKey {
	// 	return nil, nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "invalid credential2",
	// 		Data:    nil,
	// 	}
	// }

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
			Status:  http.StatusBadRequest,
			Message: "invalid credential4",
			Data:    nil,
		}
	}

	// platformName, errorResp := GetPlatformNameFromHeader(c, header_platformkey)
	// if errorResp != nil {
	// 	return nil, nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusBadRequest,
	// 		Message: errorResp.Message,
	// 		Data:    nil,
	// 	}
	// }

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

	//--------check platform--------check platform--------check platform--------
	// var dataPlatform platform.PlatformModel
	// if err := db.Table("platform_models").
	// 	Where("platform_name = ?", platformName).
	// 	First(&dataPlatform).Error; err != nil {
	// 	return nil, nil, nil, nil, &baseResp.BaseResponseModel{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "User Data Not Found",
	// 		Data:    nil,
	// 	}
	// }

	//--------check id--------check id--------check id--------
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

	return db, &userStamp, &dataAccount, nil
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
