package data

type CreateLogVehicleRequestModel struct {
	// UserId              uint   `gorm:"not null" json:"user_id" validate:"required"`
	VehicleId           uint   `gorm:"not null" json:"vehicle_id" validate:"required"`
	MeasurementTitle    string `gorm:"not null" json:"measurement_title" validate:"required"`
	CurrentOdo          string `gorm:"not null" json:"current_odo" validate:"required"`
	EstimateOdoChanging string `gorm:"not null" json:"estimate_odo_changing" validate:"required"`
	AmountExpenses      string `gorm:"not null" json:"amount_expenses" validate:"required"`
	CheckpointDate      string `gorm:"not null" json:"checkpoint_date" validate:"required"`
	Notes               string `gorm:"not null" json:"notes" validate:"required"`
}
