package data

type AccountSignUpResponseModel struct {
	Status  int                     `json:"status"`
	Message string                  `json:"message"`
	Data    *AccountSignUpDataModel `json:"Data"`
}

type AccountSignUpDataModel struct {
	UserId    uint   `json:"user_id"`
	UserStamp string `json:"user_stamp"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
