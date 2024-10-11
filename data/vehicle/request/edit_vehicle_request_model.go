package data

type EditVehicleRequestModel struct {
	// UserId         uint   `gorm:"not null" json:"user_id" validate:"required"`
	VehicleId      uint   `gorm:"not null" json:"vehicle_id" validate:"required"`
	VehicleName    string `gorm:"not null" json:"vehicle_name" validate:"required"`
	VehicleImage   string `json:"vehicle_image"`
	Year           string `gorm:"not null" json:"year" validate:"required"`
	EngineCapacity string `gorm:"not null" json:"engine_capacity" validate:"required"`
	TankCapacity   string `gorm:"not null" json:"tank_capacity" validate:"required"`
	Color          string `gorm:"not null" json:"color" validate:"required"`
	MachineNumber  string `gorm:"not null" json:"machine_number" validate:"required"`
	ChassisNumber  string `gorm:"not null" json:"chassis_number" validate:"required"`
}
