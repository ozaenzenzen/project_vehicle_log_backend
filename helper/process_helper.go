package helper

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"

	account "project_vehicle_log_backend/models/account"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type BaseResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func CheckID(db *gorm.DB, c *gin.Context, returnEmailsOrUid string) {
	iduint64, err := strconv.ParseUint(returnEmailsOrUid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  500,
			Message: "error parsing",
		})
		return
	}
	iduint := uint(iduint64)

	checkID := db.Table("account_user_models").Where("id = ?", iduint).Find(&account.AccountUserModel{
		ID: iduint,
	})

	if checkID.Error != nil {
		c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}
}

func CheckID2(db *gorm.DB, c *gin.Context, returnEmailsOrUid string, obj1 interface{}, obj2 interface{}) {
	// Get the type of the struct
	t1 := reflect.TypeOf(obj1)

	// Ensure the provided interface is a struct
	if t1.Kind() != reflect.Struct {
		fmt.Println("Provided value is not a struct!")
		return
	}

	t2 := reflect.TypeOf(obj2)

	// Ensure the provided interface is a struct
	if t2.Kind() != reflect.Struct {
		fmt.Println("Provided value is not a struct!")
		return
	}

	iduint64, err := strconv.ParseUint(returnEmailsOrUid, 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  500,
			Message: "error parsing",
		})
		return
	}
	iduint := uint(iduint64)

	checkID := db.Table("account_user_models").Where("id = ?", iduint).Find(&account.AccountUserModel{
		ID: iduint,
	})

	if checkID.Error != nil {
		c.JSON(http.StatusBadRequest, BaseResponse{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}
}
