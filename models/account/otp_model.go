package models

import "time"

type OTPModel struct {
	ID           uint      `json:"id" gorm:"primary_key"`
	Email        string    `json:"email"`
	OTP          string    `json:"otp"`
	ExpiryAt     time.Time `json:"expiry_time"`
	Count        int       `gorm:"default:1" json:"count"`
	OTPKey       string    `json:"otp_key"`
	ResendOTPKey string    `json:"resend_otp_key"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
