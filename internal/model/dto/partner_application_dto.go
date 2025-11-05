package dto

import (
	"time"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
)

type PartnerApplicationDTO struct {
	ID                  uint                    `json:"id"`
	UserID              uint                    `json:"user_id"`
	User                *UserDTO                `json:"user,omitempty"`
	BusinessName        string                  `json:"business_name"`
	BusinessAddress     string                  `json:"business_address"`
	BusinessPhone       *string                 `json:"business_phone,omitempty"`
	BusinessDescription *string                 `json:"business_description,omitempty"`
	Status              model.ApplicationStatus `json:"status"`
	RejectionReason     *string                 `json:"rejection_reason,omitempty"`
	SubmittedAt         time.Time               `json:"submitted_at"`
	DecidedAt           *time.Time              `json:"decided_at,omitempty"`
	DecidedBy           *uint                   `json:"decided_by,omitempty"`
	Decider             *UserDTO                `json:"decider,omitempty"`
}

type CreatePartnerApplicationRequest struct {
	BusinessName        string `json:"business_name" validate:"required,min=2"`
	BusinessAddress     string `json:"business_address" validate:"required,min=10"`
	BusinessPhone       string `json:"business_phone,omitempty" validate:"omitempty,min=10"`
	BusinessDescription string `json:"business_description,omitempty"`
}



func ToPartnerApplicationDTO(app *model.PartnerApplication) *PartnerApplicationDTO {
	if app == nil {
		return nil
	}

	return &PartnerApplicationDTO{
		ID:                  app.ID,
		UserID:              app.UserID,
		User:                ToUserDTO(&app.User),
		BusinessName:        app.BusinessName,
		BusinessAddress:     app.BusinessAddress,
		BusinessPhone:       app.BusinessPhone,
		BusinessDescription: app.BusinessDescription,
		Status:              app.Status,
		RejectionReason:     app.RejectionReason,
		SubmittedAt:         app.SubmittedAt,
		DecidedAt:           app.DecidedAt,
		DecidedBy:           app.DecidedBy,
		Decider:             ToUserDTO(app.Decider),
	}
}

func ToPartnerApplicationDTOList(apps []*model.PartnerApplication) []*PartnerApplicationDTO {
	result := make([]*PartnerApplicationDTO, len(apps))
	for i, app := range apps {
		result[i] = ToPartnerApplicationDTO(app)
	}
	return result
}


