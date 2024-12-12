package models

import "time"

type OTPModel struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	UserStamp string    `json:"user_stamp"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
