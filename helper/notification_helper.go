package helper

import (
	"net/http"
	"project_vehicle_log_backend/data"
	account "project_vehicle_log_backend/models/account"
	notif "project_vehicle_log_backend/models/notification"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func InsertNotification(c *gin.Context, db *gorm.DB, userData *account.AccountUserModel, title string, description string) {
	baseResponse := data.BaseResponseModel{}
	stampToken := uuid.New().String()
	inputNotifModel := notif.Notification{
		UserId:                  userData.ID,
		UserStamp:               userData.UserStamp,
		NotificationTitle:       title,
		NotificationDescription: description,
		NotificationStatus:      0,
		NotificationType:        0,
		NotificationStamp:       stampToken,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		baseResponse.Status = 400
		baseResponse.Message = resultNotif.Error.Error() + "Notif error"
		c.JSON(http.StatusBadRequest, baseResponse)
		return
	}
}
