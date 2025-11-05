package dto

import (
	"time"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
)

type PaymentDTO struct {
	ID                uint                  `json:"id"`
	BookingID         uint                  `json:"booking_id"`
	Provider          model.PaymentProvider `json:"provider"`
	ProviderPaymentID *string               `json:"provider_payment_id,omitempty"`
	Amount            float64               `json:"amount"`
	Status            model.PaymentStatus   `json:"status"`
	PaymentMethod     *string               `json:"payment_method,omitempty"`
	PaidAt            *time.Time            `json:"paid_at,omitempty"`
	FailedAt          *time.Time            `json:"failed_at,omitempty"`
	FailureReason     *string               `json:"failure_reason,omitempty"`
	CreatedAt         time.Time             `json:"created_at"`
}

type CreatePaymentRequest struct {
	Provider model.PaymentProvider `json:"provider" validate:"required,oneof=stripe midtrans"`
}

type PaymentWebhookRequest struct {
	ProviderPaymentID string  `json:"provider_payment_id" validate:"required"`
	Status            string  `json:"status" validate:"required"`
	PaymentMethod     *string `json:"payment_method,omitempty"`
	FailureReason     *string `json:"failure_reason,omitempty"`
}

func ToPaymentDTO(payment *model.Payment) *PaymentDTO {
	if payment == nil {
		return nil
	}

	return &PaymentDTO{
		ID:                payment.ID,
		BookingID:         payment.BookingID,
		Provider:          payment.Provider,
		ProviderPaymentID: payment.ProviderPaymentID,
		Amount:            payment.Amount,
		Status:            payment.Status,
		PaymentMethod:     payment.PaymentMethod,
		PaidAt:            payment.PaidAt,
		FailedAt:          payment.FailedAt,
		FailureReason:     payment.FailureReason,
		CreatedAt:         payment.CreatedAt,
	}
}

func ToPaymentDTOList(payments []*model.Payment) []*PaymentDTO {
	result := make([]*PaymentDTO, len(payments))
	for i, payment := range payments {
		result[i] = ToPaymentDTO(payment)
	}
	return result
}


