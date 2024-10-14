package data

type RefreshTokenResponseModel struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    *RefreshTokenDataModel `json:"Data"`
}

type RefreshTokenDataModel struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
