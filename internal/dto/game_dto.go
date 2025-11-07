package dto

type CreateGameRequest struct {
	CategoryID        uint    `json:"category_id" validate:"required"`
	Name              string  `json:"name" validate:"required,min=3"`
	Description       string  `json:"description,omitempty"`
	Platform          string  `json:"platform,omitempty"`
	Stock             int     `json:"stock" validate:"required,min=0"`
	RentalPricePerDay float64 `json:"rental_price_per_day" validate:"required,min=0"`
	SecurityDeposit   float64 `json:"security_deposit" validate:"required,min=0"`
	Condition         string  `json:"condition" validate:"required,oneof=excellent good fair"`
}

type UpdateGameRequest struct {
	CategoryID        uint    `json:"category_id,omitempty"`
	Name              string  `json:"name,omitempty"`
	Description       string  `json:"description,omitempty"`
	Platform          string  `json:"platform,omitempty"`
	Stock             int     `json:"stock,omitempty"`
	RentalPricePerDay float64 `json:"rental_price_per_day,omitempty"`
	SecurityDeposit   float64 `json:"security_deposit,omitempty"`
	Condition         string  `json:"condition,omitempty"`
}
