package data

import "time"

type CheckDeviceRequestModel struct {
	// UserStamp  string    `json:"user_stamp"`
	DeviceID   string    `json:"device_id" gorm:"not null"`   // UUID for device
	DeviceName string    `json:"device_name" gorm:"not null"` // e.g., "John's iPhone"
	DeviceType string    `json:"device_type" gorm:"not null"` // e.g., "iOS", "Android", "Web"
	OSVersion  string    `json:"os_version" gorm:"not null"`  // e.g., "iOS 16.1", "Android 13"
	AppVersion string    `json:"app_version" gorm:"not null"` // App version running on the device
	PushToken  string    `json:"push_token"`                  // Token for push notifications
	LastActive time.Time `json:"last_active"`                 // Timestamp for last active session
}
