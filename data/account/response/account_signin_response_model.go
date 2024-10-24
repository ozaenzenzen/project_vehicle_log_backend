package data

type AccountSignInResponseModel struct {
	Status  int                     `json:"status"`
	Message string                  `json:"message"`
	Data    *AccountSignInDataModel `json:"Data"`
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
