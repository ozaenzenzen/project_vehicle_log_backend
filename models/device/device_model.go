package models

import "time"

type DeviceModel struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	DeviceStamp string    `json:"device_stamp" gorm:"not null;uniqueIndex"` // UUID for device
	DeviceID    string    `json:"device_id" gorm:"not null;uniqueIndex"`    // UUID for device
	UserStamp   string    `json:"user_stamp"`
	DeviceName  string    `json:"device_name" gorm:"not null"` // e.g., "John's iPhone"
	DeviceType  string    `json:"device_type" gorm:"not null"` // e.g., "iOS", "Android", "Web"
	OSVersion   string    `json:"os_version" gorm:"not null"`  // e.g., "iOS 16.1", "Android 13"
	AppVersion  string    `json:"app_version" gorm:"not null"` // App version running on the device
	PushToken   string    `json:"push_token"`                  // Token for push notifications
	LastActive  time.Time `json:"last_active"`                 // Timestamp for last active session
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	// UserID     uint      `json:"user_id" gorm:"not null;index"`
	// RegisteredAt time.Time `json:"registered_at"`               // Timestamp for device registration
}
