package models

import "time"

type QNAModel struct {
	Id         uint      `gorm:"not null" json:"id" gorm:"primary_key"`
	TopicId    uint      `gorm:"not null" json:"topic_id" validate:"required"`
	TopicTitle string    `gorm:"not null" json:"topic_title" validate:"required"`
	Question   string    `gorm:"not null" json:"question" validate:"required"`
	Answer     string    `gorm:"not null" json:"answer" validate:"required"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
