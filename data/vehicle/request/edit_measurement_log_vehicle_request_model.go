package data

type EditMeasurementLogVehicleRequestModel struct {
	ID                  uint   `json:"id" `
	VehicleId           uint   `json:"vehicle_id" `
	MeasurementTitle    string `json:"measurement_title" `
	CurrentOdo          string `json:"current_odo" `
	EstimateOdoChanging string `json:"estimate_odo_changing" `
	AmountExpenses      string `json:"amount_expenses" `
	CheckpointDate      string `json:"checkpoint_date" `
	Notes               string `json:"notes" `
}
