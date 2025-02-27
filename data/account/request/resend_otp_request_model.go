package data

type ResendOTPRequestModel struct {
	ResendOTPKey string `gorm:"not null" json:"resend_otp_key" binding:"required"`
}
