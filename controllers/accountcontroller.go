package controllers

import (
	"net/http"
	jwthelper "project_vehicle_log_backend/helper"
	account "project_vehicle_log_backend/models/account"
	user "project_vehicle_log_backend/models/account"
	notif "project_vehicle_log_backend/models/notification"
	"strconv"

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
	ID    uint   `json:"id" gorm:"primary_key"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	// Link     string `json:"link"`
	// Typeuser uint   `json:"typeuser"`
	Token string `json:"token"`
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
			Message: "Account not match",
			// Typeuser: nil,
			UserData: nil,
		})
		return
	}

	// token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	// 	"uid":   dataUser.ID,
	// 	"email": dataUser.Email,
	// 	"exp":   time.Now().Add(time.Hour * 168).Unix(), // Token expires in 168 hour or 1 week
	// })

	// tokenString, err := token.SignedString([]byte("ozaenzenzen"))

	tokenString, err := jwthelper.GenerateJWTToken(strconv.FormatUint(uint64(dataUser.ID), 10), dataUser.Email)

	if err != nil {
		c.JSON(http.StatusNotFound, AccountUserSignInResponse{
			Status:   http.StatusInternalServerError,
			Message:  "Failed to generate token",
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
			Token: tokenString,
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
	headertoken := c.Request.Header.Get("token")

	if headertoken == "" {
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  400,
			Message: "token empty",
		})
		return
	}
	isValid, err := jwthelper.VerifyToken(headertoken)

	if isValid == true {
		if err != nil {
			c.JSON(http.StatusBadRequest, AccountUserGetUserResponse{
				Status:   http.StatusBadRequest,
				Message:  err.Error(),
				UserData: nil,
			})
			return
		}

		var userData user.AccountUserModel

		tokenRaw, err := jwthelper.DecodeJWTToken(headertoken)
		// fmt.Printf("\ntoken raw %v", tokenRaw)
		if err != nil {
			c.JSON(http.StatusBadRequest, AccountUserGetUserResponse{
				Status:   http.StatusBadRequest,
				Message:  err.Error(),
				UserData: nil,
			})
			return
		}

		emails := tokenRaw["email"].(string)

		// if err := db.Where("id = ?", c.Param("id")).First(&userData).Error; err != nil {
		if err := db.Where("email = ?", emails).First(&userData).Error; err != nil {
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
		return
	} else {
		c.JSON(http.StatusBadRequest, AccountUserGetUserResponse{
			Status:   http.StatusBadRequest,
			Message:  "invalid token",
			UserData: nil,
		})
		return
	}

}

type EditProfileRequest struct {
	ProfilePicture string `json:"profile_picture"`
	Name           string `json:"name"`
}

type EditProfileResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func EditProfile(c *gin.Context) {
	headerid := c.Request.Header.Get("usd")
	if headerid == "" {
		c.JSON(http.StatusBadRequest, EditProfileResponse{
			Status:  400,
			Message: "headers empty",
		})
		return
	}

	headertoken := c.Request.Header.Get("token")

	if headertoken == "" {
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  400,
			Message: "token empty",
		})
		return
	}

	var editProfileRequest EditProfileRequest
	if err := c.ShouldBindJSON(&editProfileRequest); err != nil {
		c.JSON(http.StatusBadRequest, EditProfileResponse{
			Status:  500,
			Message: err.Error(),
		})
		return
	}
	db := c.MustGet("db").(*gorm.DB)

	//--------check id--------check id--------check id--------

	iduint64, err := strconv.ParseUint(headerid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, EditProfileResponse{
			Status:  500,
			Message: "error parsing",
		})
		return
	}
	iduint := uint(iduint64)

	checkID := db.Table("account_user_models").Where("id = ?", headerid).Find(&account.AccountUserModel{
		ID: iduint,
	})

	if checkID.Error != nil {
		c.JSON(http.StatusBadRequest, EditProfileResponse{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}

	//--------check id--------check id--------check id--------

	isValid, err := jwthelper.VerifyToken(headertoken)
	if err != nil {
		c.JSON(http.StatusBadRequest, EditProfileResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	if db.Error != nil {
		c.JSON(http.StatusBadRequest, EditProfileResponse{
			Status:  400,
			Message: db.Error.Error(),
		})
		return
	}
	// result := db.Create(&vehicleDataOutput)
	result := db.Table("account_user_models").Where("id = ?", headerid).Update(&editProfileRequest)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, EditProfileResponse{
			Status:  400,
			Message: result.Error.Error(),
		})
		return
	}

	inputNotifModel := notif.Notification{
		UserId:                  iduint,
		NotificationTitle:       "Edit Profile",
		NotificationDescription: "Anda Telah Mengubah Data Profile",
		NotificationStatus:      0,
		NotificationType:        0,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		c.JSON(http.StatusBadRequest, EditVehicleResponse{
			Status:  400,
			Message: result.Error.Error() + "Notif error",
		})
		return
	}

	editProfileResponse := EditVehicleResponse{
		Status:  201,
		Message: "Edit profile success",
	}

	if isValid == true {
		c.JSON(http.StatusOK, editProfileResponse)
		return
	}
}
