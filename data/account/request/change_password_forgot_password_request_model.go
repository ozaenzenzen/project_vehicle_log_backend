package data

type ChangePasswordForgotPasswordRequestModel struct {
	ForgotKey          string `gorm:"not null" json:"forgot_key" binding:"required"`
	NewPassword        string `gorm:"not null" json:"new_password" binding:"required"`
	ConfirmNewPassword string `gorm:"not null" json:"confirm_new_password" binding:"required"`
}
