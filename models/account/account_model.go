package models

import "time"

type AccountUserModel struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	Name            string    `json:"name"`
	ProfilePicture  string    `json:"profile_picture"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
