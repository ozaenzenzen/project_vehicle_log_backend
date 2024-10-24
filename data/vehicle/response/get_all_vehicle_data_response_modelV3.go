package data

import "time"

type GetAllVehicleDataResponseModelV3 struct {
	Status  int                         `json:"status"`
	Message string                      `json:"message"`
	Data    *[]GetAllVehicleDataModelV3 `json:"Data"`
}

type GetAllVehicleDataModelV3 struct {
	Id                          uint                                      `json:"id" gorm:"primary_key"`
	UserId                      uint                                      `json:"user_id"`
	VehicleName                 string                                    `json:"vehicle_name"`
	VehicleImage                string                                    `json:"vehicle_image"`
	Year                        string                                    `json:"year"`
	EngineCapacity              string                                    `json:"engine_capacity"`
	TankCapacity                string                                    `json:"tank_capacity"`
	Color                       string                                    `json:"color"`
	MachineNumber               string                                    `json:"machine_number"`
	ChassisNumber               string                                    `json:"chassis_number"`
	MeasurementTitle            *[]string                                 `json:"measurment_title"`
	VehicleMeasurementLogModels *[]GetAllVehicleMeasurementLogDataModelV3 `json:"measurement_data" gorm:"foreignKey:vehicle_id;references:Id"`
}

type GetAllVehicleMeasurementLogDataModelV3 struct {
	Id                  uint      `json:"id" `
	UserId              uint      `json:"user_id"`
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
