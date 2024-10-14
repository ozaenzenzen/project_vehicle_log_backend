package data

type GetAllVehicleDataResponseModelV2 struct {
	Status  int                       `json:"status"`
	Message string                    `json:"message"`
	Data    *GetAllVehicleDataModelV2 `json:"Data"`
}

type GetAllVehicleDataModelV2 struct {
	CurrentPage int                    `json:"current_page"`
	NextPage    int                    `json:"next_page"`
	TotalPages  int                    `json:"total_pages"`
	TotalItems  int                    `json:"total_items"`
	Data        *[]DataGetAllVehicleV2 `json:"list_data"`
}

type DataGetAllVehicleV2 struct {
	Id             uint   `json:"id" gorm:"primary_key"`
	UserId         uint   `json:"user_id"`
	UserStamp      string `json:"user_stamp"`
	VehicleName    string `json:"vehicle_name"`
	VehicleImage   string `json:"vehicle_image"`
	Year           string `json:"year"`
	EngineCapacity string `json:"engine_capacity"`
	TankCapacity   string `json:"tank_capacity"`
	Color          string `json:"color"`
	MachineNumber  string `json:"machine_number"`
	ChassisNumber  string `json:"chassis_number"`
}
