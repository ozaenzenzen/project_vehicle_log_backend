package data

import "time"

type CheckDeviceResponseModel struct {
	Status  int                   `json:"status"`
	Message string                `json:"message"`
	Data    *CheckDeviceDataModel `json:"Data"`
}

type CheckDeviceDataModel struct {
	// UserID     uint      `json:"user_id" gorm:"not null;index"`
	DeviceID    string    `json:"device_id" gorm:"not null;uniqueIndex"` // UUID for device
	UserStamp   string    `json:"user_stamp"`
	DeviceStamp string    `json:"device_stamp"`
	DeviceName  string    `json:"device_name" gorm:"not null"` // e.g., "John's iPhone"
	DeviceType  string    `json:"device_type" gorm:"not null"` // e.g., "iOS", "Android", "Web"
	OSVersion   string    `json:"os_version" gorm:"not null"`  // e.g., "iOS 16.1", "Android 13"
	AppVersion  string    `json:"app_version" gorm:"not null"` // App version running on the device
	PushToken   string    `json:"push_token"`                  // Token for push notifications
	LastActive  time.Time `json:"last_active"`                 // Timestamp for last active session
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
