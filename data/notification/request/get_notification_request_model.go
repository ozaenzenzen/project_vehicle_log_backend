package data

type GetNotificationaRequestModel struct {
	CurrentPage int     `json:"current_page" binding:"required"`
	Limit       int     `json:"limit" binding:"required"`
	SortOrder   *string `json:"sort_order"` // Optional (asc/desc)
}
