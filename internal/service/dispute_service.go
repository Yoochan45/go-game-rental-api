package service

import (
	"errors"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"github.com/Yoochan45/go-game-rental-api/internal/repository"
)

var (
	ErrCannotCreateDispute = errors.New("cannot create dispute for this booking")
)

type DisputeService interface {
	// Customer/Partner methods
	CreateDispute(reporterID uint, bookingID uint, disputeData *model.Dispute) error

	// Admin methods
	GetAllDisputes(requestorRole model.UserRole, limit, offset int) ([]*model.Dispute, error)
}

type disputeService struct {
	disputeRepo repository.DisputeRepository
	bookingRepo repository.BookingRepository
}

func NewDisputeService(disputeRepo repository.DisputeRepository, bookingRepo repository.BookingRepository) DisputeService {
	return &disputeService{
		disputeRepo: disputeRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *disputeService) CreateDispute(reporterID uint, bookingID uint, disputeData *model.Dispute) error {
	// Validate booking exists
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return ErrBookingNotFound
	}

	// Check if user is involved in the booking (either customer or partner)
	if booking.UserID != reporterID && booking.PartnerID != reporterID {
		return ErrCannotCreateDispute
	}

	// Can only create dispute for confirmed, active, or completed bookings
	validStatuses := []model.BookingStatus{
		model.BookingConfirmed,
		model.BookingActive,
		model.BookingCompleted,
	}

	isValidStatus := false
	for _, status := range validStatuses {
		if booking.Status == status {
			isValidStatus = true
			break
		}
	}

	if !isValidStatus {
		return ErrCannotCreateDispute
	}

	// Set dispute details
	disputeData.BookingID = bookingID
	disputeData.ReporterID = reporterID
	disputeData.Status = model.DisputeOpen

	// Update booking status to disputed
	err = s.bookingRepo.UpdateStatus(bookingID, model.BookingDisputed)
	if err != nil {
		return err
	}

	return s.disputeRepo.Create(disputeData)
}

func (s *disputeService) GetAllDisputes(requestorRole model.UserRole, limit, offset int) ([]*model.Dispute, error) {
	if !s.canManageDisputes(requestorRole) {
		return nil, ErrInsufficientPermission
	}

	return s.disputeRepo.GetAllDisputes(limit, offset)
}

func (s *disputeService) canManageDisputes(role model.UserRole) bool {
	return role == model.RoleAdmin || role == model.RoleSuperAdmin
}
