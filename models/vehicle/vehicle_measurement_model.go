package models

import "time"

type VehicleMeasurementModel struct {
	Id                  uint      `json:"id" gorm:"primary_key"`
	UserId              uint      `json:"user_id"`
	VehicleId           uint      `json:"vehicle_id"`
	MeasurementTitle    string    `json:"measurement_title"`
	CurrentOdo          string    `json:"current_odd"`
	EstimateOdoChanging string    `json:"estimate_odo_changing"`
	AmountExpenses      string    `json:"amount_expenses"`
	CheckpointDate      string    `json:"checkpoint_date"`
	Notes               string    `json:"notes"`
	CreatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt           time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
