package data

type GetListLogTypeResponseModel struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    *[]LogDataModel `json:"Data"`
}

type LogDataModel struct {
	MeasurementTitle string `json:"measurement_title"`
}
