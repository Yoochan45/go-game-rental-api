package dto

type CreateBookingRequest struct {
	GameID    uint   `json:"game_id" validate:"required"`
	StartDate string `json:"start_date" validate:"required"` // Ubah jadi string
	EndDate   string `json:"end_date" validate:"required"`   // Ubah jadi string
	Notes     string `json:"notes,omitempty"`
}
