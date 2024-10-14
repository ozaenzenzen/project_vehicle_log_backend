package helper

import (
	"net/http"
	baseResp "project_vehicle_log_backend/data"
	account "project_vehicle_log_backend/models/account"
	notif "project_vehicle_log_backend/models/notification"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

func InsertNotification(c *gin.Context, db *gorm.DB, userData *account.AccountUserModel, title string, description string) *baseResp.BaseResponseModel {
	baseResponse := baseResp.BaseResponseModel{}
	notificationStamp := uuid.New().String()
	inputNotifModel := notif.Notification{
		UserId:                  userData.ID,
		UserStamp:               userData.UserStamp,
		NotificationTitle:       title,
		NotificationDescription: description,
		NotificationStatus:      0,
		NotificationType:        0,
		NotificationStamp:       notificationStamp,
	}

	resultNotif := db.Table("notifications").Create(&inputNotifModel)
	if resultNotif.Error != nil {
		baseResponse.Status = http.StatusInternalServerError
		baseResponse.Message = resultNotif.Error.Error() + "Notif error"
		return &baseResponse
	}
	return nil
}
