package data

type ChangePasswordRequestModel struct {
	OldPassword        string `gorm:"not null" json:"old_password"  binding:"required"`
	NewPassword        string `gorm:"not null" json:"new_password"  binding:"required"`
	ConfirmNewPassword string `gorm:"not null" json:"confirm_new_password"  binding:"required"`
}
