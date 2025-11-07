package service

import (
	"errors"
	"time"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"github.com/Yoochan45/go-game-rental-api/internal/repository"
)

var (
	ErrBookingNotFound       = errors.New("booking not found")
	ErrBookingNotOwned       = errors.New("you don't own this booking")
	ErrBookingInvalidDate    = errors.New("invalid booking dates")
	ErrBookingCannotCancel   = errors.New("cannot cancel booking in current status")
	ErrGameStockInsufficient = errors.New("insufficient stock")
)

type BookingService interface {
	// Customer
	Create(userID uint, bookingData *model.Booking) error
	GetUserBookings(userID uint, limit, offset int) ([]*model.Booking, int64, error)
	GetByID(userID uint, bookingID uint) (*model.Booking, error)
	Cancel(userID uint, bookingID uint) error

	// Admin
	GetAll(requestorRole model.UserRole, limit, offset int) ([]*model.Booking, int64, error)
	UpdateStatus(requestorRole model.UserRole, bookingID uint, status model.BookingStatus) error

	// System (for payment)
	ConfirmPayment(bookingID uint) error
	FailPayment(bookingID uint) error
}

type bookingService struct {
	bookingRepo repository.BookingRepository
	gameRepo    repository.GameRepository
}

func NewBookingService(
	bookingRepo repository.BookingRepository,
	gameRepo repository.GameRepository,
) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
		gameRepo:    gameRepo,
	}
}

func (s *bookingService) Create(userID uint, bookingData *model.Booking) error {
	game, err := s.gameRepo.GetByID(bookingData.GameID)
	if err != nil {
		return ErrGameNotFound
	}

	if !game.IsActive {
		return errors.New("game is not available for booking")
	}

	if bookingData.StartDate.After(bookingData.EndDate) || bookingData.StartDate.Before(time.Now().Truncate(24*time.Hour)) {
		return ErrBookingInvalidDate
	}

	available, err := s.gameRepo.CheckAvailability(game.ID)
	if err != nil {
		return err
	}
	if !available {
		return ErrGameStockInsufficient
	}

	rentalDays := int(bookingData.EndDate.Sub(bookingData.StartDate).Hours()/24) + 1
	totalRentalPrice := float64(rentalDays) * game.RentalPricePerDay
	totalAmount := totalRentalPrice + game.SecurityDeposit

	bookingData.UserID = userID
	bookingData.RentalDays = rentalDays
	bookingData.DailyPrice = game.RentalPricePerDay
	bookingData.TotalRentalPrice = totalRentalPrice
	bookingData.SecurityDeposit = game.SecurityDeposit
	bookingData.TotalAmount = totalAmount
	bookingData.Status = model.BookingPending

	err = s.gameRepo.ReserveStock(game.ID)
	if err != nil {
		return err
	}

	return s.bookingRepo.Create(bookingData)
}

func (s *bookingService) GetUserBookings(userID uint, limit, offset int) ([]*model.Booking, int64, error) {
	bookings, err := s.bookingRepo.GetUserBookings(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.bookingRepo.CountUserBookings(userID)
	return bookings, count, err
}

func (s *bookingService) GetByID(userID uint, bookingID uint) (*model.Booking, error) {
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return nil, ErrBookingNotFound
	}

	if booking.UserID != userID {
		return nil, ErrBookingNotOwned
	}

	return booking, nil
}

func (s *bookingService) Cancel(userID uint, bookingID uint) error {
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return ErrBookingNotFound
	}

	if booking.UserID != userID {
		return ErrBookingNotOwned
	}

	if booking.Status != model.BookingPending && booking.Status != model.BookingConfirmed {
		return ErrBookingCannotCancel
	}

	err = s.gameRepo.ReleaseStock(booking.GameID)
	if err != nil {
		return err
	}

	return s.bookingRepo.UpdateStatus(bookingID, model.BookingCancelled)
}

func (s *bookingService) GetAll(requestorRole model.UserRole, limit, offset int) ([]*model.Booking, int64, error) {
	if !s.canManageBookings(requestorRole) {
		return nil, 0, ErrInsufficientPermission
	}

	bookings, err := s.bookingRepo.GetAllBookings(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.bookingRepo.Count()
	return bookings, count, err
}

func (s *bookingService) UpdateStatus(requestorRole model.UserRole, bookingID uint, status model.BookingStatus) error {
	if !s.canManageBookings(requestorRole) {
		return ErrInsufficientPermission
	}

	return s.bookingRepo.UpdateStatus(bookingID, status)
}

func (s *bookingService) ConfirmPayment(bookingID uint) error {
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return ErrBookingNotFound
	}

	if booking.Status != model.BookingPending {
		return errors.New("booking is not in pending status")
	}

	return s.bookingRepo.UpdateStatus(bookingID, model.BookingConfirmed)
}

func (s *bookingService) FailPayment(bookingID uint) error {
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return ErrBookingNotFound
	}

	err = s.gameRepo.ReleaseStock(booking.GameID)
	if err != nil {
		return err
	}

	return s.bookingRepo.UpdateStatus(bookingID, model.BookingCancelled)
}

func (s *bookingService) canManageBookings(role model.UserRole) bool {
	return role == model.RoleAdmin || role == model.RoleSuperAdmin
}
