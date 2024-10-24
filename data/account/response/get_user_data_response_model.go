package data

type GetUserDataResponseModel struct {
	Status  int               `json:"status"`
	Message string            `json:"message"`
	Data    *GetUserDataModel `json:"Data"`
}
type GetUserDataModel struct {
	ID             uint   `json:"id" gorm:"primary_key"`
	UserStamp      string `json:"user_stamp"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	ProfilePicture string `json:"profile_picture"`
}
