package dto

import (
	"time"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
)

type BookingDTO struct {
	ID                  uint                `json:"id"`
	UserID              uint                `json:"user_id"`
	User                *UserDTO            `json:"user,omitempty"`
	GameID              uint                `json:"game_id"`
	Game                *GameDTO            `json:"game,omitempty"`
	PartnerID           uint                `json:"partner_id"`
	Partner             *UserDTO            `json:"partner,omitempty"`
	StartDate           time.Time           `json:"start_date"`
	EndDate             time.Time           `json:"end_date"`
	RentalDays          int                 `json:"rental_days"`
	DailyPrice          float64             `json:"daily_price"`
	TotalRentalPrice    float64             `json:"total_rental_price"`
	SecurityDeposit     float64             `json:"security_deposit"`
	TotalAmount         float64             `json:"total_amount"`
	Status              model.BookingStatus `json:"status"`
	Notes               *string             `json:"notes,omitempty"`
	HandoverConfirmedAt *time.Time          `json:"handover_confirmed_at,omitempty"`
	ReturnConfirmedAt   *time.Time          `json:"return_confirmed_at,omitempty"`
	Payment             *PaymentDTO         `json:"payment,omitempty"`
	Review              *ReviewDTO          `json:"review,omitempty"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
}

type CreateBookingRequest struct {
	GameID    uint      `json:"game_id" validate:"required"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
	Notes     string    `json:"notes,omitempty"`
}

type BookingListResponse struct {
	Bookings   []*BookingDTO `json:"bookings"`
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
}

func ToBookingDTO(booking *model.Booking) *BookingDTO {
	if booking == nil {
		return nil
	}

	return &BookingDTO{
		ID:                  booking.ID,
		UserID:              booking.UserID,
		User:                ToUserDTO(&booking.User),
		GameID:              booking.GameID,
		Game:                ToGameDTO(&booking.Game),
		PartnerID:           booking.PartnerID,
		Partner:             ToUserDTO(&booking.Partner),
		StartDate:           booking.StartDate,
		EndDate:             booking.EndDate,
		RentalDays:          booking.RentalDays,
		DailyPrice:          booking.DailyPrice,
		TotalRentalPrice:    booking.TotalRentalPrice,
		SecurityDeposit:     booking.SecurityDeposit,
		TotalAmount:         booking.TotalAmount,
		Status:              booking.Status,
		Notes:               booking.Notes,
		HandoverConfirmedAt: booking.HandoverConfirmedAt,
		ReturnConfirmedAt:   booking.ReturnConfirmedAt,
		Payment:             ToPaymentDTO(booking.Payment),
		Review:              ToReviewDTO(booking.Review),
		CreatedAt:           booking.CreatedAt,
		UpdatedAt:           booking.UpdatedAt,
	}
}

func ToBookingDTOList(bookings []*model.Booking) []*BookingDTO {
	result := make([]*BookingDTO, len(bookings))
	for i, booking := range bookings {
		result[i] = ToBookingDTO(booking)
	}
	return result
}

func FromCreateBookingRequest(req *CreateBookingRequest) *model.Booking {
	return &model.Booking{
		GameID:    req.GameID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Notes:     &req.Notes,
	}
}
