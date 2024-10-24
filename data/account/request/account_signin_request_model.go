package data

type AccountSignInRequestModel struct {
	Email    string `gorm:"not null" json:"email"  binding:"required"`
	Password string `gorm:"not null" json:"password"  binding:"required"`
}
