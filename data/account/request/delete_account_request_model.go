package data

type DeleteAccountRequestModel struct {
	Reason *string `gorm:"not null" json:"reason"  binding:"required"`
}
