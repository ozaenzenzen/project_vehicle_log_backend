package data

type SendOTPForgotPasswordRequestModel struct {
	Email string `gorm:"not null" json:"email" binding:"required"`
}
