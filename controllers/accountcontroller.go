package controllers

import (
	"net/http"
	req "project_vehicle_log_backend/data/account/request"
	resp "project_vehicle_log_backend/data/account/response"
	helper "project_vehicle_log_backend/helper"

	// account "project_vehicle_log_backend/models/account"
	account "project_vehicle_log_backend/models/account"
	notif "project_vehicle_log_backend/models/notification"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func RefreshToken(c *gin.Context) {
	db, userStamp, _, errorResp := helper.CustomValidatorWithRefreshToken(c, true)
	if errorResp != nil {
		c.JSON(errorResp.Status, errorResp)
		return
	}

	accessToken, refreshToken, errGenerateJWT := helper.GenerateUserTokenV2(*userStamp) // using stamp
	if errGenerateJWT != nil {
		c.JSON(http.StatusInternalServerError, resp.RefreshTokenResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate token",
			Data:    nil,
		})
		return
	}

	// store refresh token
	storeRefreshToken := db.Table("account_user_models").
		Where("user_stamp = ?", *userStamp).
		Update(&req.RefreshTokenRequestModel{RefreshToken: refreshToken})
	if storeRefreshToken.Error != nil {
		c.JSON(http.StatusInternalServerError, resp.RefreshTokenResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "error storing",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp.RefreshTokenResponseModel{
		Status:  http.StatusOK,
		Message: "Refresh Token Success",
		Data: &resp.RefreshTokenDataModel{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	})
}

func SignUpAccount(c *gin.Context) {
	var signUpReq req.AccountSignUpRequestModel
	if errBindJSON := c.ShouldBindJSON(&signUpReq); errBindJSON != nil {
		c.JSON(http.StatusBadRequest, resp.AccountSignUpResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Data tidak lengkap",
			Data:    nil,
		})
		return
	}

	validate := validator.New()
	if errValidate := validate.Struct(signUpReq); errValidate != nil {
		c.JSON(http.StatusBadRequest, resp.AccountSignUpResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Data tidak lengkap",
			Data:    nil,
		})
		return
	}

	hashPw, errPw := helper.HashPassword(signUpReq.Password)
	if errPw != nil {
		c.JSON(http.StatusInternalServerError, resp.AccountSignUpResponseModel{
			Status:  http.StatusInternalServerError,
			Message: errPw.Error(),
			Data:    nil,
		})
		return
	}

	hashCpw, errCpw := helper.HashPassword(signUpReq.ConfirmPassword)
	if errCpw != nil {
		c.JSON(http.StatusInternalServerError, resp.AccountSignUpResponseModel{
			Status:  http.StatusInternalServerError,
			Message: errCpw.Error(),
			Data:    nil,
		})
		return
	}

	stampToken := uuid.New().String()
	// userStamp := signUpReq.Name + stampToken

	insertDBPayload := account.AccountUserModel{
		Name:            signUpReq.Name,
		Email:           signUpReq.Email,
		UserStamp:       stampToken,
		Phone:           signUpReq.Phone,
		Password:        hashPw,
		ConfirmPassword: hashCpw,
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		c.JSON(http.StatusInternalServerError, resp.AccountSignUpResponseModel{
			Status:  http.StatusInternalServerError,
			Message: db.Error.Error(),
			Data:    nil,
		})
		return
	}

	result := db.FirstOrCreate(
		&insertDBPayload,
		account.AccountUserModel{
			Email: signUpReq.Email,
		},
	)
	if result.Value == nil && result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, resp.AccountSignUpResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Record found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, resp.AccountSignUpResponseModel{
		Status:  http.StatusCreated,
		Message: "Account Created Successfully",
		Data: &resp.AccountSignUpDataModel{
			UserId:    insertDBPayload.ID,
			UserStamp: insertDBPayload.UserStamp,
			Name:      signUpReq.Name,
			Email:     signUpReq.Email,
			Phone:     signUpReq.Phone,
		},
	})
}

func SignInAccount(c *gin.Context) {
	var signInReq req.AccountSignInRequestModel
	if errBindJSON := c.ShouldBindJSON(&signInReq); errBindJSON != nil {
		c.JSON(http.StatusBadRequest, resp.AccountSignInResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Data belum lengkap",
			Data:    nil,
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		c.JSON(http.StatusInternalServerError, resp.AccountSignInResponseModel{
			Status:  http.StatusInternalServerError,
			Message: db.Error.Error(),
			Data:    nil,
		})
		return
	}

	var tableAccount account.AccountUserModel
	resultCheckEmail := db.Where("email = ?", signInReq.Email).
		First(&tableAccount)
	if resultCheckEmail.Error != nil {
		c.JSON(http.StatusUnauthorized, resp.AccountSignInResponseModel{
			Status:  http.StatusUnauthorized,
			Message: "Invalid user email or password",
			Data:    nil,
		})
		return
	}

	checkHashPw := helper.CheckPasswordHash(signInReq.Password, tableAccount.Password)
	if !checkHashPw {
		c.JSON(http.StatusUnauthorized, resp.AccountSignInResponseModel{
			Status:  http.StatusUnauthorized,
			Message: "Invalid user email or password",
		})
		return
	}

	userToken, refreshToken, errGenerateJWT := helper.GenerateUserTokenV2(tableAccount.UserStamp) // using stamp
	if errGenerateJWT != nil {
		c.JSON(http.StatusNotFound, resp.AccountSignInResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "Failed to generate token",
			Data:    nil,
		})
		return
	}

	// store refresh token
	storeRefreshToken := db.Table("account_user_models").
		Where("user_stamp = ?", tableAccount.UserStamp).
		Update(&req.RefreshTokenRequestModel{RefreshToken: refreshToken})

	if storeRefreshToken.Error != nil {
		c.JSON(http.StatusUnauthorized, resp.AccountSignInResponseModel{
			Status:  http.StatusUnauthorized,
			Message: "Failed Storing",
			Data:    nil,
		})
		return
	}

	// userToken, errGenerateAccessToken := helper.GenerateJWTToken(
	// 	strconv.FormatUint(uint64(tableAccount.ID), 10),
	// 	signInReq.Email,
	// )
	// if errGenerateAccessToken != nil {
	// 	c.JSON(http.StatusInternalServerError, resp.AccountSignInResponseModel{
	// 		Status:  http.StatusInternalServerError,
	// 		Message: "Failed to generate token",
	// 		Data:    nil,
	// 	})
	// 	return
	// }

	// result := db.Where("email = ?", signInReq.Email).
	// 	Where("password = ?", tableAccount.Password).
	// 	First(&tableAccount)
	// if result.Error != nil {
	// 	c.JSON(http.StatusNotFound, resp.AccountSignInResponseModel{
	// 		Status:  http.StatusNotFound,
	// 		Message: "Failed authorization",
	// 		Data:    nil,
	// 	})
	// 	return
	// }

	c.JSON(http.StatusOK, resp.AccountSignInResponseModel{
		Status:  http.StatusOK,
		Message: "Account SignIn Successfully",
		Data: &resp.AccountSignInDataModel{
			ID:           tableAccount.ID,
			Name:         tableAccount.Name,
			UserStamp:    tableAccount.UserStamp,
			Email:        signInReq.Email,
			Phone:        tableAccount.Phone,
			Token:        userToken,
			RefreshToken: refreshToken,
		},
	})
}

func GetUserData(c *gin.Context) {
	header_usertoken := c.Request.Header.Get("token")
	if header_usertoken == "" {
		c.JSON(http.StatusBadRequest, resp.GetUserDataResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid Credential1",
			Data:    nil,
		})
		return
	}

	isValidToken, errVerifyToken := helper.VerifyToken(header_usertoken)
	if errVerifyToken != nil || !isValidToken {
		c.JSON(http.StatusBadRequest, resp.GetUserDataResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid Credential2",
			Data:    nil,
		})
		return
	}

	tokenRaw, errDecodeToken := helper.DecodeJWTToken(header_usertoken)
	if errDecodeToken != nil {
		c.JSON(http.StatusBadRequest, resp.GetUserDataResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid Credential3",
			Data:    nil,
		})
		return
	}

	// userEmail := tokenRaw["email"].(string)
	var userStamp string
	if resultUserStamp, ok := tokenRaw["user_stamp"].(string); ok {
		userStamp = resultUserStamp
	} else {
		userStamp = ""
		c.JSON(http.StatusInternalServerError, resp.GetUserDataResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "Failed Parsing",
			Data:    nil,
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		c.JSON(http.StatusInternalServerError, resp.GetUserDataResponseModel{
			Status:  http.StatusInternalServerError,
			Message: db.Error.Error(),
			Data:    nil,
		})
		return
	}

	var tableAccount account.AccountUserModel
	if err := db.Where("user_stamp = ?", userStamp).
		First(&tableAccount).Error; err != nil {
		c.JSON(http.StatusBadRequest, resp.GetUserDataResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Not Found",
			Data:    nil,
		})
		return
	}

	c.JSON(http.StatusOK, resp.GetUserDataResponseModel{
		Status:  http.StatusOK,
		Message: "Get User Data Successfully",
		Data: &resp.GetUserDataModel{
			ID:             tableAccount.ID,
			UserStamp:      tableAccount.UserStamp,
			Name:           tableAccount.Name,
			Email:          tableAccount.Email,
			Phone:          tableAccount.Phone,
			ProfilePicture: tableAccount.ProfilePicture,
		},
	})
}

func EditProfile(c *gin.Context) {
	var editProfileReq req.EditProfileRequesModel
	if err := c.ShouldBindJSON(&editProfileReq); err != nil {
		c.JSON(http.StatusBadRequest, resp.EditProfileResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Data Tidak Lengkap",
		})
		return
	}

	header_usertoken := c.Request.Header.Get("token")
	if header_usertoken == "" {
		c.JSON(http.StatusBadRequest, resp.EditProfileResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid Credential 1",
		})
		return
	}

	isValidUserToken, errVerifyToken := helper.VerifyToken(header_usertoken)
	if errVerifyToken != nil && !isValidUserToken {
		c.JSON(http.StatusBadRequest, resp.EditProfileResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid Credential 2",
		})
		return
	}

	tokenRaw, errDecodeToken := helper.DecodeJWTToken(header_usertoken)
	if errDecodeToken != nil {
		c.JSON(http.StatusBadRequest, resp.EditProfileResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Invalid Credential 3",
		})
		return
	}

	// ids := tokenRaw["uid"].(string)
	var userStamp string
	if resultUserStamp, ok := tokenRaw["user_stamp"].(string); ok {
		userStamp = resultUserStamp
	} else {
		userStamp = ""
		c.JSON(http.StatusInternalServerError, resp.GetUserDataResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "Failed Parsing1",
			Data:    nil,
		})
		return
	}

	// iduint64, errParseUserID := strconv.ParseUint(userID, 10, 32)
	// if errParseUserID != nil {
	// 	c.JSON(http.StatusBadRequest, resp.EditProfileResponseModel{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "Failed Parsing2",
	// 	})
	// 	return
	// }
	// iduint := uint(iduint64)

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		c.JSON(http.StatusBadRequest, resp.EditProfileResponseModel{
			Status:  http.StatusBadRequest,
			Message: db.Error.Error(),
		})
		return
	}

	// checkID := db.Table("account_user_models").
	// 	Where("user_stamp = ?", userStamp).
	// 	Find(&account.AccountUserModel{
	// 		ID: iduint,
	// 	})
	// if checkID.Error != nil {
	// 	c.JSON(http.StatusBadRequest, resp.EditProfileResponseModel{
	// 		Status:  http.StatusBadRequest,
	// 		Message: "Not Found",
	// 	})
	// 	return
	// }

	result := db.Table("account_user_models").
		Where("user_stamp = ?", userStamp).
		Update(&editProfileReq)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, resp.EditProfileResponseModel{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}

	inputNotifModel := notif.Notification{
		// UserId:                  iduint,
		UserStamp:               userStamp,
		NotificationTitle:       "Edit Profile",
		NotificationDescription: "Anda Telah Mengubah Data Profile",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		c.JSON(http.StatusInternalServerError, resp.EditProfileResponseModel{
			Status:  http.StatusInternalServerError,
			Message: "Notif Error",
		})
		return
	}

	c.JSON(http.StatusOK, resp.EditProfileResponseModel{
		Status:  http.StatusAccepted,
		Message: "Edit Profile Successfully",
	})

}
