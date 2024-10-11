package data

type BaseResponseModel struct {
	Status  int          `json:"status"`
	Message string       `json:"message"`
	Data    *interface{} `json:"Data"`
}