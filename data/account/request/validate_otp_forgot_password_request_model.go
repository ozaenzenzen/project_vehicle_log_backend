package data

type ValidateOTPForgotPasswordRequestModel struct {
	Email     string `gorm:"not null" json:"email" binding:"required"`
	OTP       string `gorm:"not null" json:"otp" binding:"required"`
	OTPKey    string `gorm:"not null" json:"otp_key" binding:"required"`
	ForgotKey string `gorm:"not null" json:"forgot_key" binding:"required"`
}
