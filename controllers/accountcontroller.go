package controllers

import (
	"net/http"
	baseResp "project_vehicle_log_backend/data"
	req "project_vehicle_log_backend/data/account/request"
	resp "project_vehicle_log_backend/data/account/response"
	helper "project_vehicle_log_backend/helper"

	// account "project_vehicle_log_backend/models/account"
	account "project_vehicle_log_backend/models/account"

	otpService "project_vehicle_log_backend/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func ResendOTP(c *gin.Context) {
	baseResponse := resp.ResendOTPResponseModel{}

	var reqBody req.ResendOTPRequestModel

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Data tidak lengkap1"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Data tidak lengkap2"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = db.Error.Error()
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	var dataOTP account.OTPModel
	if err := db.Table("otp_models").
		Where("resend_otp_key = ?", reqBody.ResendOTPKey).
		Last(&dataOTP).
		Update(account.OTPModel{ResendOTPKey: ""}).
		Error; err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Not yet add otp"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	var dataAccount account.AccountUserModel
	if err := db.Table("account_user_models").
		Where("email = ?", dataOTP.Email).
		First(&dataAccount).Error; err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "User Data Not Found"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	otpKey, resendOtpKey, resultResp := SendEmailAndStoreOTPHelper(db, dataAccount.Email)
	if resultResp != nil {
		baseResponse.Status = resultResp.Status
		baseResponse.Message = resultResp.Message
		baseResponse.Data = nil
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	baseResponse.Status = http.StatusOK
	baseResponse.Message = "Resend OTP Success"
	baseResponse.Data = &resp.ResendOTPResponseDataModel{
		OTPKey:       *otpKey,
		ResendOTPKey: *resendOtpKey,
	}
	c.JSON(baseResponse.Status, baseResponse)
}

func OTPValidation(c *gin.Context) {
	baseResponse := resp.OTPVerificationResponseModel{}

	var reqBody req.OTPValidationRequestModel

	if err := c.ShouldBindJSON(&reqBody); err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Data tidak lengkap1"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	validate := validator.New()
	if err := validate.Struct(reqBody); err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Data tidak lengkap2"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = db.Error.Error()
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	var dataAccount account.AccountUserModel
	if err := db.Table("account_user_models").
		Where("email = ?", reqBody.Email).
		First(&dataAccount).Error; err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "User Data Not Found"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	if dataAccount.IsActivated == 1 {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Already Activated"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	if dataAccount.StatusAccount != 1 {
		baseResponse.Status = http.StatusUnauthorized
		baseResponse.Message = "Account Disabled"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	var dataOTP account.OTPModel
	if err := db.Table("otp_models").
		Where("email = ? AND otp_key = ?", reqBody.Email, reqBody.OTPKey).
		Last(&dataOTP).
		Update(account.OTPModel{OTPKey: ""}).
		Error; err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "Failed OTP Process, Please Resend OTP"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	// if dataOTP.Count >= 5 {
	// 	baseResponse.Status = http.StatusLocked
	// 	baseResponse.Message = "Maximum Hit"
	// 	c.JSON(baseResponse.Status, baseResponse)
	// 	return
	// }
	// dataCreatedAt := dataOTP.CreatedAt
	// // Add 24 hours to createdAt
	// expiryTime := dataCreatedAt.Add(24 * time.Hour)

	// // Get the current time
	// currentTime := time.Now().UTC()

	// // Check if 24 hours have passed
	// if currentTime.After(expiryTime) {
	// 	fmt.Println("24 hours have passed since createdAt.")
	// } else {
	// 	fmt.Println("24 hours have not yet passed since createdAt.")
	// }

	isOTPVerified := otpService.VerifyOTP(
		dataOTP.Email,
		reqBody.OTP,
		dataOTP.OTP,
		dataOTP.ExpiryAt,
	)
	if !isOTPVerified {
		baseResponse.Status = http.StatusUnauthorized
		baseResponse.Message = "Failed Verify OTP"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	// dataAccount.IsActivated = 1
	if err := db.Table("account_user_models").
		Where("email = ?", reqBody.Email).
		First(&dataAccount).Update(&account.AccountUserModel{IsActivated: 1}).Error; err != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = "User Data Not Found"
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	baseResponse.Status = http.StatusOK
	baseResponse.Message = "Success OTP Verified"
	c.JSON(baseResponse.Status, baseResponse)
}

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

	accessToken, accessTokenExpiryTime, refreshToken, refreshTokenExpiryTime, errGenerateJWT := helper.GenerateUserTokenV3(*userStamp) // using stamp
	// accessToken, refreshToken, errGenerateJWT := helper.GenerateUserTokenV2(*userStamp) // using stamp
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
	baseResponse.Data = &resp.RefreshTokenDataModelV2{
		AccessToken:            accessToken,
		AccessTokenExpiryTime:  *accessTokenExpiryTime,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryTime: *refreshTokenExpiryTime,
	}
	c.JSON(http.StatusOK, baseResponse)
}

func SendEmailAndStoreOTPHelper(db *gorm.DB, inputEmail string) (*string, *string, *baseResp.BaseResponseModel) {
	baseResponse := baseResp.BaseResponseModel{}

	email, otp, expiration, errorData := otpService.SendEmailRegisterUser(inputEmail)
	if errorData != nil {
		baseResponse.Status = http.StatusBadRequest
		baseResponse.Message = *errorData
		baseResponse.Data = nil
		return nil, nil, &baseResponse
	}

	otpKey := helper.RandomHash(inputEmail + "OTPKey")
	resendOtpKey := helper.RandomHash(inputEmail + "ResendOTPKey")

	otpModel := account.OTPModel{
		Email:        *email,
		OTP:          *otp,
		ExpiryAt:     *expiration,
		OTPKey:       *otpKey,
		ResendOTPKey: *resendOtpKey,
	}

	resultStore := db.Create(&otpModel)
	if resultStore.Error != nil {
		baseResponse.Status = 400
		baseResponse.Message = resultStore.Error.Error()
		baseResponse.Data = nil
		return nil, nil, &baseResponse
	}

	return otpKey, resendOtpKey, nil
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

	otpKey, resendOtpKey, resultResp := SendEmailAndStoreOTPHelper(db, signUpReq.Email)
	if resultResp != nil {
		baseResponse.Status = resultResp.Status
		baseResponse.Message = resultResp.Message
		baseResponse.Data = nil
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	baseResponse.Status = http.StatusCreated
	baseResponse.Message = "Account Created Successfully"
	baseResponse.Data = &resp.AccountSignUpDataModel{
		UserId:       insertDBPayload.ID,
		UserStamp:    insertDBPayload.UserStamp,
		Name:         signUpReq.Name,
		Email:        signUpReq.Email,
		Phone:        signUpReq.Phone,
		OTPKey:       *otpKey,
		ResendOTPKey: *resendOtpKey,
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

	if tableAccount.IsActivated == 0 {
		baseResponse.Status = http.StatusUnauthorized
		baseResponse.Message = "User unverified"
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

	accessToken, accessTokenExpiryTime, refreshToken, refreshTokenExpiryTime, errGenerateJWT := helper.GenerateUserTokenV3(tableAccount.UserStamp) // using stamp
	// userToken, refreshToken, errGenerateJWT := helper.GenerateUserTokenV2(tableAccount.UserStamp) // using stamp
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
	baseResponse.Data = &resp.AccountSignInDataModelV2{
		ID:                     tableAccount.ID,
		Name:                   tableAccount.Name,
		UserStamp:              tableAccount.UserStamp,
		Email:                  signInReq.Email,
		Phone:                  tableAccount.Phone,
		AccessToken:            accessToken,
		AccessTokenExpiryTime:  *accessTokenExpiryTime,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryTime: *refreshTokenExpiryTime,
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
		c.JSON(baseResponse.Status, baseResponse)
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
		c.JSON(baseResponse.Status, baseResponse)
		return
	}

	baseResponse.Status = http.StatusAccepted
	baseResponse.Message = "Edit Profile Successfully"
	c.JSON(http.StatusOK, baseResponse)

}
