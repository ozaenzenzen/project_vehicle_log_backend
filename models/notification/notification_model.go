package models

import "time"

type Notification struct {
	NotificationId uint `json:"notification_id" gorm:"primary_key"`
	UserId         uint `gorm:"not null" json:"user_id" validate:"required"`
	// OrganizationId          uint      `gorm:"not null" json:"organization_id" validate:"required"`
	// EventId                 uint      `gorm:"not null" json:"event_id" validate:"required"`
	NotificationTitle       string    `gorm:"not null" json:"notification_title" validate:"required"`
	NotificationDescription string    `gorm:"not null" json:"notification_description" validate:"required"`
	NotificationStatus      uint      `gorm:"not null" json:"notification_status" validate:"required"`
	NotificationType        uint      `gorm:"not null" json:"notification_type" validate:"required"`
	CreatedAt               time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt               time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Notification Status
// 0: available
// 1: read
// 2: close

// VolunteersStatus uint      `json:"volunteers_status" validate:"required"`
// EventStatus      uint      `json:"event_status" validate:"required"`
// EventStatus
// 0: waiting
// 1: approved
// 2: rejected
// 3: finished
// 4: submitted volunteer
