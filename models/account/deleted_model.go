package models

import "time"

type DeletedModel struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Email     string    `json:"email"`
	UserStamp string    `json:"user_stamp"`
	Reason    string    `json:"reason"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
