package data

import "time"

type GetLogVehicleDataResponseModelV2 struct {
	Status  int                       `json:"status"`
	Message string                    `json:"message"`
	Data    *GetLogVehicleDataModelV2 `json:"Data"`
}

type GetLogVehicleDataModelV2 struct {
	CurrentPage int                    `json:"current_page"`
	NextPage    int                    `json:"next_page"`
	TotalPages  int                    `json:"total_pages"`
	TotalItems  int                    `json:"total_items"`
	Data        *[]DataGetLogVehicleV2 `json:"list_data"`
}

type DataGetLogVehicleV2 struct {
	Id                  uint      `json:"id" gorm:"primary_key"`
	UserId              uint      `json:"user_id"`
	UserStamp           string    `json:"user_stamp"`
	VehicleId           uint      `json:"vehicle_id"`
	MeasurementTitle    string    `json:"measurement_title"`
	CurrentOdo          string    `json:"current_odo"`
	EstimateOdoChanging string    `json:"estimate_odo_changing"`
	AmountExpenses      string    `json:"amount_expenses"`
	CheckpointDate      string    `json:"checkpoint_date"`
	Notes               string    `json:"notes"`
	CreatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
