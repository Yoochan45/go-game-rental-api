package dto

type CreatePartnerApplicationRequest struct {
	BusinessName        string `json:"business_name" validate:"required,min=2"`
	BusinessAddress     string `json:"business_address" validate:"required,min=10"`
	BusinessPhone       string `json:"business_phone,omitempty" validate:"omitempty,min=10"`
	BusinessDescription string `json:"business_description,omitempty"`
}


