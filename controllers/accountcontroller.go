package controllers

import (
	"net/http"
	req "project_vehicle_log_backend/data/account/request"
	resp "project_vehicle_log_backend/data/account/response"
	helper "project_vehicle_log_backend/helper"

	// account "project_vehicle_log_backend/models/account"
	account "project_vehicle_log_backend/models/account"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func RefreshToken(c *gin.Context) {
	baseResponse := resp.RefreshTokenResponseModel{}

	db, userStamp, _, errorResp := helper.CustomValidatorWithRefreshToken(c, true)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		baseResponse.Data = nil
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	accessToken, refreshToken, errGenerateJWT := helper.GenerateUserTokenV2(*userStamp) // using stamp
	if errGenerateJWT != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = "Failed to generate token"
		baseResponse.Data = nil
		c.JSON(http.StatusInternalServerError, baseResponse)
		return
	}

	// store refresh token
	storeRefreshToken := db.Table("account_user_models").
		Where("user_stamp = ?", *userStamp).
		Update(&req.RefreshTokenRequestModel{RefreshToken: refreshToken})
	if storeRefreshToken.Error != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = "error storing"
		baseResponse.Data = nil
		c.JSON(http.StatusInternalServerError, baseResponse)
		return
	}

	baseResponse.Status = http.StatusOK
	baseResponse.Message = "Refresh Token Success"
	baseResponse.Data = &resp.RefreshTokenDataModel{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	c.JSON(http.StatusOK, baseResponse)
}

func SignUpAccount(c *gin.Context) {
	baseResponse := resp.AccountSignUpResponseModel{}

	var signUpReq req.AccountSignUpRequestModel
	if errBindJSON := c.ShouldBindJSON(&signUpReq); errBindJSON != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Data tidak lengkap"
		baseResponse.Data = nil
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	validate := validator.New()
	if errValidate := validate.Struct(signUpReq); errValidate != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Data tidak lengkap"
		baseResponse.Data = nil
		c.JSON(http.StatusBadRequest, resp.AccountSignUpResponseModel{
			Status:  http.StatusBadRequest,
			Message: "Data tidak lengkap",
			Data:    nil,
		})
		return
	}

	hashPw, errPw := helper.HashPassword(signUpReq.Password)
	if errPw != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = errPw.Error()
		baseResponse.Data = nil
		c.JSON(http.StatusInternalServerError, baseResponse)
		return
	}

	hashCpw, errCpw := helper.HashPassword(signUpReq.ConfirmPassword)
	if errCpw != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = errCpw.Error()
		baseResponse.Data = nil
		c.JSON(http.StatusInternalServerError, baseResponse)
		return
	}

	stampToken := uuid.New().String()

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
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = db.Error.Error()
		baseResponse.Data = nil
		c.JSON(http.StatusInternalServerError, baseResponse)
		return
	}

	result := db.FirstOrCreate(
		&insertDBPayload,
		account.AccountUserModel{
			Email: signUpReq.Email,
		},
	)
	if result.Value == nil && result.RowsAffected == 0 {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Record found"
		baseResponse.Data = nil
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	baseResponse.Status = http.StatusCreated
	baseResponse.Message = "Account Created Successfully"
	baseResponse.Data = &resp.AccountSignUpDataModel{
		UserId:    insertDBPayload.ID,
		UserStamp: insertDBPayload.UserStamp,
		Name:      signUpReq.Name,
		Email:     signUpReq.Email,
		Phone:     signUpReq.Phone,
	}
	c.JSON(http.StatusCreated, baseResponse)
}

func SignInAccount(c *gin.Context) {
	baseResponse := resp.AccountSignInResponseModel{}

	var signInReq req.AccountSignInRequestModel
	if errBindJSON := c.ShouldBindJSON(&signInReq); errBindJSON != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Data belum lengkap"
		baseResponse.Data = nil
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = db.Error.Error()
		baseResponse.Data = nil
		c.JSON(http.StatusInternalServerError, baseResponse)
		return
	}

	var tableAccount account.AccountUserModel
	resultCheckEmail := db.Where("email = ?", signInReq.Email).
		First(&tableAccount)
	if resultCheckEmail.Error != nil {
		baseResponse.Status = http.StatusUnauthorized
		baseResponse.Message = "Invalid user email or password"
		baseResponse.Data = nil
		c.JSON(http.StatusUnauthorized, baseResponse)
		return
	}

	checkHashPw := helper.CheckPasswordHash(signInReq.Password, tableAccount.Password)
	if !checkHashPw {
		baseResponse.Status = http.StatusUnauthorized
		baseResponse.Message = "Invalid user email or password"
		baseResponse.Data = nil
		c.JSON(http.StatusUnauthorized, baseResponse)
		return
	}

	userToken, refreshToken, errGenerateJWT := helper.GenerateUserTokenV2(tableAccount.UserStamp) // using stamp
	if errGenerateJWT != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = "Failed to generate token"
		baseResponse.Data = nil
		c.JSON(http.StatusNotFound, baseResponse)
		return
	}

	// store refresh token
	storeRefreshToken := db.Table("account_user_models").
		Where("user_stamp = ?", tableAccount.UserStamp).
		Update(&req.RefreshTokenRequestModel{RefreshToken: refreshToken})
	if storeRefreshToken.Error != nil {
		baseResponse.Status = http.StatusUnauthorized
		baseResponse.Message = "Failed Storing"
		baseResponse.Data = nil
		c.JSON(http.StatusUnauthorized, baseResponse)
		return
	}

	baseResponse.Status = http.StatusOK
	baseResponse.Message = "Account SignIn Successfully"
	baseResponse.Data = &resp.AccountSignInDataModel{
		ID:           tableAccount.ID,
		Name:         tableAccount.Name,
		UserStamp:    tableAccount.UserStamp,
		Email:        signInReq.Email,
		Phone:        tableAccount.Phone,
		Token:        userToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, baseResponse)
}

func GetUserData(c *gin.Context) {
	baseResponse := resp.GetUserDataResponseModel{}

	_, _, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	baseResponse.Status = http.StatusOK
	baseResponse.Message = "Get User Data Successfully"
	baseResponse.Data = &resp.GetUserDataModel{
		ID:             userData.ID,
		UserStamp:      userData.UserStamp,
		Name:           userData.Name,
		Email:          userData.Email,
		Phone:          userData.Phone,
		ProfilePicture: userData.ProfilePicture,
	}
	c.JSON(http.StatusOK, baseResponse)
}

func EditProfile(c *gin.Context) {
	baseResponse := resp.GetUserDataResponseModel{}

	var editProfileReq req.EditProfileRequesModel
	if err := c.ShouldBindJSON(&editProfileReq); err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Data Tidak Lengkap"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}

	db, _, userData, errorResp := helper.CustomValidatorAC(c)
	if errorResp != nil {
		baseResponse.Status = errorResp.Status
		baseResponse.Message = errorResp.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	result := db.Table("account_user_models").Where("id = ?", userData.ID).Update(editProfileReq)
	if result.Error != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = "Terjadi kesalahan"
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	respNotif := helper.InsertNotification(
		c,
		db,
		userData,
		"Edit Profile",
		"Anda Telah Mengubah Data Profile",
	)
	if respNotif != nil {
		baseResponse.Status = respNotif.Status
		baseResponse.Message = respNotif.Message
		c.JSON(errorResp.Status, baseResponse)
		return
	}

	baseResponse.Status = http.StatusAccepted
	baseResponse.Message = "Edit Profile Successfully"
	c.JSON(http.StatusOK, baseResponse)

}
