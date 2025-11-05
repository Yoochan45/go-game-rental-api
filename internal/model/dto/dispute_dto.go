package dto

import (
	"time"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
)

type DisputeDTO struct {
	ID          uint                `json:"id"`
	BookingID   uint                `json:"booking_id"`
	Booking     *BookingDTO         `json:"booking,omitempty"`
	ReporterID  uint                `json:"reporter_id"`
	Reporter    *UserDTO            `json:"reporter,omitempty"`
	Type        model.DisputeType   `json:"type"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	Status      model.DisputeStatus `json:"status"`
	Resolution  *string             `json:"resolution,omitempty"`
	ResolvedBy  *uint               `json:"resolved_by,omitempty"`
	Resolver    *UserDTO            `json:"resolver,omitempty"`
	CreatedAt   time.Time           `json:"created_at"`
	ResolvedAt  *time.Time          `json:"resolved_at,omitempty"`
}

type CreateDisputeRequest struct {
	Type        model.DisputeType `json:"type" validate:"required,oneof=payment item_condition late_return no_show other"`
	Title       string            `json:"title" validate:"required,min=5"`
	Description string            `json:"description" validate:"required,min=10"`
}

type ResolveDisputeRequest struct {
	Action     string `json:"action" validate:"required,oneof=investigate resolve close"`
	Resolution string `json:"resolution,omitempty" validate:"required_if=Action resolve,required_if=Action close"`
}

type DisputeListResponse struct {
	Disputes   []*DisputeDTO `json:"disputes"`
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
}

func ToDisputeDTO(dispute *model.Dispute) *DisputeDTO {
	if dispute == nil {
		return nil
	}

	return &DisputeDTO{
		ID:          dispute.ID,
		BookingID:   dispute.BookingID,
		Booking:     ToBookingDTO(&dispute.Booking),
		ReporterID:  dispute.ReporterID,
		Reporter:    ToUserDTO(&dispute.Reporter),
		Type:        dispute.Type,
		Title:       dispute.Title,
		Description: dispute.Description,
		Status:      dispute.Status,
		Resolution:  dispute.Resolution,
		ResolvedBy:  dispute.ResolvedBy,
		Resolver:    ToUserDTO(dispute.Resolver),
		CreatedAt:   dispute.CreatedAt,
		ResolvedAt:  dispute.ResolvedAt,
	}
}

func ToDisputeDTOList(disputes []*model.Dispute) []*DisputeDTO {
	result := make([]*DisputeDTO, len(disputes))
	for i, dispute := range disputes {
		result[i] = ToDisputeDTO(dispute)
	}
	return result
}

func FromCreateDisputeRequest(req *CreateDisputeRequest) *model.Dispute {
	return &model.Dispute{
		Type:        req.Type,
		Title:       req.Title,
		Description: req.Description,
	}
}
