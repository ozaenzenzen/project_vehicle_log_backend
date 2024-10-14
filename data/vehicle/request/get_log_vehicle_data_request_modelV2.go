package data

type GetLogVehicleDataRequestModelV2 struct {
	CurrentPage      int     `json:"current_page" binding:"required"`
	Limit            int     `json:"limit" binding:"required"`
	VehicleID        *string `json:"vehicle_id"`
	MeasurementTitle *string `json:"measurement_title"`
	SortOrder        *string `json:"sort_order"` // Optional (asc/desc)
}
