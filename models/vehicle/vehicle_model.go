package models

import "time"

type VehicleModel struct {
	Id             uint      `json:"id" gorm:"primary_key"`
	UserId         uint      `json:"user_id"`
	UserStamp      string    `json:"user_stamp"`
	VehicleName    string    `json:"vehicle_name"`
	VehicleImage   string    `json:"vehicle_image"`
	Year           string    `json:"year"`
	EngineCapacity string    `json:"engine_capacity"`
	TankCapacity   string    `json:"tank_capacity"`
	Color          string    `json:"color"`
	MachineNumber  string    `json:"machine_number"`
	ChassisNumber  string    `json:"chassis_number"`
	VehicleStamp   string    `json:"vehicle_stamp"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
