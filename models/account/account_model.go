package models

import "time"

type AccountUserModel struct {
	ID              uint      `json:"id" gorm:"primary_key"`
	Name            string    `json:"name"`
	ProfilePicture  string    `json:"profile_picture"`
	Email           string    `json:"email"`
	Phone           string    `json:"phone"`
	UserStamp       string    `json:"user_stamp"`
	RefreshToken    string    `json:"refresh_token"`
	Password        string    `json:"password"`
	ConfirmPassword string    `json:"confirm_password"`
	IsActivated     int       `gorm:"default:0" json:"is_activated"`
	TypeAccount     int       `gorm:"default:0" json:"type_account"`
	StatusAccount   int       `gorm:"default:1" json:"status_account"`
	CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Is Activated
// 0: not yet otp verified
// 1: verified otp

// Type Account
// 0: basic
// 1: plus

// Status Account
// 0: disabled
// 1: normal
