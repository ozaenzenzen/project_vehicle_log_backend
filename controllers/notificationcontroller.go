package controllers

import (
	"fmt"
	"log"
	"net/http"

	"project_vehicle_log_backend/helper"
	jwthelper "project_vehicle_log_backend/helper"
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
	headertoken := c.Request.Header.Get("token")
	isValid, err := jwthelper.ValidateTokenJWT(c, db, headertoken)
	if err != nil {
		c.JSON(http.StatusBadRequest, GetNotificationResponse{
			Status:  http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}
	if isValid == true {

		var notificationList []notif.Notification

		if c.Param("id") == "" {
			c.JSON(http.StatusBadRequest, GetNotificationResponse{
				Status:  http.StatusBadRequest,
				Message: "Empty parameter",
			})
			return
		}
		returnEmailsOrUid := jwthelper.GetDataTokenJWT(headertoken, false)
		if returnEmailsOrUid != c.Param("id") {
			c.JSON(http.StatusBadRequest, GetNotificationResponse{
				Status:  http.StatusBadRequest,
				Message: "Wrong token",
			})
			return
		}

		//--------check id--------check id--------check id--------

		helper.CheckID(db, c, c.Param("id"))

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
	} else {
		c.JSON(http.StatusBadRequest, GetNotificationResponse{
			Status:  http.StatusBadRequest,
			Message: "invalid token",
		})
		return
	}
}
