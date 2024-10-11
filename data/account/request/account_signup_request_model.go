package data

type AccountSignUpRequestModel struct {
	Name            string `gorm:"not null" json:"name"  binding:"required,max=30"`
	Email           string `gorm:"not null" json:"email" binding:"required"`
	Phone           string `gorm:"not null" json:"phone"  binding:"required,max=14"`
	Password        string `gorm:"not null" json:"password" binding:"required"`
	ConfirmPassword string `gorm:"not null" json:"confirmPassword" binding:"required"`
}
