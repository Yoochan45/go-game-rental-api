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

type PaymentListResponse struct {
	Payments   []*PaymentDTO `json:"payments"`
	TotalCount int64         `json:"total_count"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
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

func (d *PaymentDTO) GetProvider() model.PaymentProvider {
	return d.Provider
}

func (d *PaymentDTO) GetProviderPaymentID() *string {
	return d.ProviderPaymentID
}

func (d *PaymentDTO) GetAmount() float64 {
	return d.Amount
}

func (d *PaymentDTO) GetStatus() model.PaymentStatus {
	return d.Status
}

func (d *PaymentDTO) GetPaymentMethod() *string {
	return d.PaymentMethod
}

func (d *PaymentDTO) GetPaidAt() *time.Time {
	return d.PaidAt
}

func (d *PaymentDTO) GetFailedAt() *time.Time {
	return d.FailedAt
}

func (d *PaymentDTO) GetFailureReason() *string {
	return d.FailureReason
}

func (d *PaymentDTO) GetCreatedAt() time.Time {
	return d.CreatedAt
}

func (d *PaymentDTO) GetBookingID() uint {
	return d.BookingID
}

func (d *PaymentDTO) GetID() uint {
	return d.ID
}
