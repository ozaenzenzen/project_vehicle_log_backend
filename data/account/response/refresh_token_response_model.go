package data

import "time"

type RefreshTokenResponseModel struct {
	Status  int                      `json:"status"`
	Message string                   `json:"message"`
	Data    *RefreshTokenDataModelV2 `json:"Data"`
}

type RefreshTokenDataModel struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenDataModelV2 struct {
	AccessToken            string    `json:"access_token"`
	AccessTokenExpiryTime  time.Time `json:"access_token_expiry_time"`
	RefreshToken           string    `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time `json:"refresh_token_expiry_time"`
}
