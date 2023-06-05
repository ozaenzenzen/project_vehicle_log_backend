package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	account "project_vehicle_log_backend/models/account"
	notif "project_vehicle_log_backend/models/notification"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type GetNotificationResponse struct {
	Status       int                   `json:"status"`
	Message      string                `json:"message"`
	Notification *[]notif.Notification `json:"notification"`
}

func GetNotificationByUserId(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	var notificationList []notif.Notification

	//--------check id--------check id--------check id--------

	iduint64, err := strconv.ParseUint(c.Param("id"), 10, 32)

	if err != nil {
		c.JSON(http.StatusBadRequest, GetNotificationResponse{
			Status:  500,
			Message: "error parsing",
		})
		return
	}
	iduint := uint(iduint64)

	checkID := db.Table("account_user_models").Where("id = ?", c.Param("id")).Find(&account.AccountUserModel{
		ID: iduint,
	})

	if checkID.Error != nil {
		c.JSON(http.StatusBadRequest, GetNotificationResponse{
			Status:  400,
			Message: checkID.Error.Error(),
		})
		return
	}

	//--------check id--------check id--------check id--------

	result := db.Table("notifications").Where("user_id = ?", c.Param("id")).Find(&notificationList)

	if result.Value == nil {
		log.Println(fmt.Sprintf("error log notification1: %s", result.Error))
		c.JSON(http.StatusInternalServerError, GetNotificationResponse{
			Status:       500,
			Message:      "get notification failed 1",
			Notification: &[]notif.Notification{},
		})
		return
	}

	c.JSON(http.StatusOK, GetNotificationResponse{
		Status:       200,
		Message:      "get notification success",
		Notification: &notificationList,
	})
}
