package controllers

import (
	"net/http"
	user "project_vehicle_log_backend/models/account"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type AccountSingUpResponse struct {
	Status  int                           `json:"status"`
	Message string                        `json:"message"`
	Data    *AccountUserDataResponseModel `json:"account_data"`
}

// ID       uint   `json:"id"`
type AccountUserData struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Phone           string `json:"phone"`
	Link            string `json:"link"`
	Typeuser        uint   `json:"typeuser"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirmPassword"`
}

type AccountUserDataResponseModel struct {
	UserId   uint   `json:"user_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Link     string `json:"link"`
	Typeuser uint   `json:"typeuser"`
}

func SignUpAccount(c *gin.Context) {
	var accountInput AccountUserData
	if err := c.ShouldBindJSON(&accountInput); err != nil {
		c.JSON(http.StatusBadRequest, AccountSingUpResponse{
			Status:  500,
			Message: err.Error(),
			Data:    &AccountUserDataResponseModel{},
		})
		return
	}
	accountResponsePayload := user.AccountUserModel{
		Name:            accountInput.Name,
		Email:           accountInput.Email,
		Phone:           accountInput.Phone,
		Password:        accountInput.Password,
		ConfirmPassword: accountInput.ConfirmPassword,
	}

	db := c.MustGet("db").(*gorm.DB)
	if db.Error != nil {
		c.JSON(http.StatusBadRequest, AccountSingUpResponse{
			Status:  400,
			Message: db.Error.Error(),
			Data:    nil,
		})
		return
	}

	result := db.FirstOrCreate(&accountResponsePayload, user.AccountUserModel{Email: accountInput.Email})

	if result.Value == nil && result.RowsAffected == 0 {
		// log.Println(fmt.Sprintf("log SignUpAccountV2 Value: %s", result.Value))
		// log.Println(fmt.Sprintf("log SignUpAccountV2 RowsAffected: %d", result.RowsAffected))
		c.JSON(http.StatusBadRequest, AccountSingUpResponse{
			Status:  400,
			Message: "Record found",
			Data:    nil,
		})
		return
	}

	createAccountResponse := AccountSingUpResponse{
		Status:  201,
		Message: "Account created successfully",
		Data: &AccountUserDataResponseModel{
			UserId:   accountResponsePayload.ID,
			Name:     accountInput.Name,
			Email:    accountInput.Email,
			Phone:    accountInput.Phone,
			Link:     accountInput.Link,
			Typeuser: accountInput.Typeuser,
		},
	}

	c.JSON(http.StatusCreated, createAccountResponse)
}

type UserDataModelSignIn struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Link     string `json:"link"`
	Typeuser uint   `json:"typeuser"`
}

type AccountUserSignInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AccountUserSignInResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	// Typeuser *uint  `json:"typeuser"`
	UserData *UserDataModelSignIn `json:"userdata"`
}

func SignInAccount(c *gin.Context) {
	var dataUser user.AccountUserModel
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		c.JSON(http.StatusBadRequest, AccountUserSignInResponse{
			Status:  500,
			Message: err.Error(),
		})
		return
	}
	db := c.MustGet("db").(*gorm.DB)

	if err := db.Where("email = ?", dataUser.Email).Where("password = ?", dataUser.Password).First(&dataUser).Error; err != nil {
		c.JSON(http.StatusNotFound, AccountUserSignInResponse{
			Status:  404,
			Message: "Account SignIn Failed",
			// Typeuser: nil,
			UserData: nil,
		})
		return
	}

	accountSignInResponse := AccountUserSignInResponse{
		Status:  200,
		Message: "Account SignIn Successfully",
		// Typeuser: &dataUser.Typeuser,
		UserData: &UserDataModelSignIn{
			ID:    dataUser.ID,
			Name:  dataUser.Name,
			Email: dataUser.Email,
			Phone: dataUser.Phone,
		},
	}

	c.JSON(http.StatusOK, accountSignInResponse)
}

type UserDataModel struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Link     string `json:"link"`
	Typeuser uint   `json:"typeuser"`
}

type AccountUserGetUserResponse struct {
	Status   int            `json:"status"`
	Message  string         `json:"message"`
	UserData *UserDataModel `json:"userdata"`
}

func GetUserData(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var userData user.AccountUserModel

	if err := db.Where("id = ?", c.Param("id")).First(&userData).Error; err != nil {
		c.JSON(http.StatusBadRequest, AccountUserGetUserResponse{
			Status:   400,
			Message:  "User Data Not Found",
			UserData: nil,
		})
		return
	}

	c.JSON(http.StatusOK, AccountUserGetUserResponse{
		Status:  200,
		Message: "get user data success",
		UserData: &UserDataModel{
			ID:    userData.ID,
			Name:  userData.Name,
			Email: userData.Email,
			Phone: userData.Phone,
		},
	})
}
