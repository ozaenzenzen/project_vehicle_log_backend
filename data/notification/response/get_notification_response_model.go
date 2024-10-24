package data

import "time"

type GetNotificationResponseModel struct {
	Status  int                        `json:"status"`
	Message string                     `json:"message"`
	Data    *GetNotificationPagination `json:"Data"`
}

type GetNotificationPagination struct {
	CurrentPage int                         `json:"current_page"`
	NextPage    int                         `json:"next_page"`
	TotalPages  int                         `json:"total_pages"`
	TotalItems  int                         `json:"total_items"`
	Data        *[]GetNotificationDataModel `json:"list_data"`
}

type GetNotificationDataModel struct {
	NotificationId          uint      `json:"notification_id"`
	UserId                  uint      `json:"user_id"`
	UserStamp               string    `json:"user_stamp"`
	NotificationTitle       string    `json:"notification_title"`
	NotificationDescription string    `json:"notification_description"`
	NotificationStatus      uint      `json:"notification_status"`
	NotificationType        uint      `json:"notification_type"`
	NotificationStamp       string    `json:"notification_stamp"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
}
