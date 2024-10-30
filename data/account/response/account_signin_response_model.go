package data

import "time"

type AccountSignInResponseModel struct {
	Status  int                       `json:"status"`
	Message string                    `json:"message"`
	Data    *AccountSignInDataModelV2 `json:"Data"`
}

type AccountSignInDataModel struct {
	ID           uint   `json:"id" gorm:"primary_key"`
	UserStamp    string `json:"user_stamp"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type AccountSignInDataModelV2 struct {
	ID                     uint      `json:"id" gorm:"primary_key"`
	UserStamp              string    `json:"user_stamp"`
	Name                   string    `json:"name"`
	Email                  string    `json:"email"`
	Phone                  string    `json:"phone"`
	AccessToken            string    `json:"access_token"`
	AccessTokenExpiryTime  time.Time `json:"access_token_expiry_time"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time `json:"refresh_token_expiry_time"`
}
