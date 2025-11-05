package dto

import "github.com/Yoochan45/go-game-rental-api/internal/model"

type CreateDisputeRequest struct {
	Type        model.DisputeType `json:"type" validate:"required,oneof=payment item_condition late_return no_show other"`
	Title       string            `json:"title" validate:"required,min=5"`
	Description string            `json:"description" validate:"required,min=10"`
}
