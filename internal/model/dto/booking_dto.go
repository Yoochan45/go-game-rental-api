package dto

import "time"

type CreateBookingRequest struct {
	GameID    uint      `json:"game_id" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
	Notes     string    `json:"notes,omitempty"`
}
