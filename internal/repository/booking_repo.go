package repository

import (
	"time"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"gorm.io/gorm"
)

type BookingRepository interface {
	// Basic CRUD
	Create(booking *model.Booking) error
	GetByID(id uint) (*model.Booking, error)
	GetByIDWithRelations(id uint) (*model.Booking, error)
	Update(booking *model.Booking) error
	Delete(id uint) error

	// Query methods
	GetUserBookings(userID uint, limit, offset int) ([]*model.Booking, error)
	GetPartnerBookings(partnerID uint, limit, offset int) ([]*model.Booking, error)
	GetBookingsByStatus(status model.BookingStatus, limit, offset int) ([]*model.Booking, error)

	// Status updates
	UpdateStatus(bookingID uint, status model.BookingStatus) error
	UpdateHandoverConfirmation(bookingID uint) error
	UpdateReturnConfirmation(bookingID uint) error

	// Date conflicts check
	CheckDateConflicts(gameID uint, startDate, endDate time.Time, excludeBookingID *uint) (bool, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) GetByID(id uint) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Where("id = ?", id).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetByIDWithRelations(id uint) (*model.Booking, error) {
	var booking model.Booking
	err := r.db.Preload("User").Preload("Game").Preload("Partner").Preload("Payment").
		Where("id = ?", id).First(&booking).Error
	if err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) Update(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) Delete(id uint) error {
	return r.db.Delete(&model.Booking{}, id).Error
}

func (r *bookingRepository) GetUserBookings(userID uint, limit, offset int) ([]*model.Booking, error) {
	var bookings []*model.Booking
	err := r.db.Preload("Game").Preload("Partner").Preload("Payment").
		Where("user_id = ?", userID).Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetPartnerBookings(partnerID uint, limit, offset int) ([]*model.Booking, error) {
	var bookings []*model.Booking
	err := r.db.Preload("User").Preload("Game").Preload("Payment").
		Where("partner_id = ?", partnerID).Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetBookingsByStatus(status model.BookingStatus, limit, offset int) ([]*model.Booking, error) {
	var bookings []*model.Booking
	err := r.db.Preload("User").Preload("Game").Preload("Partner").
		Where("status = ?", status).Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) UpdateStatus(bookingID uint, status model.BookingStatus) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", bookingID).Update("status", status).Error
}

func (r *bookingRepository) UpdateHandoverConfirmation(bookingID uint) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", bookingID).
		Updates(map[string]interface{}{
			"handover_confirmed_at": gorm.Expr("CURRENT_TIMESTAMP"),
			"status":                model.BookingActive,
		}).Error
}

func (r *bookingRepository) UpdateReturnConfirmation(bookingID uint) error {
	return r.db.Model(&model.Booking{}).Where("id = ?", bookingID).
		Updates(map[string]interface{}{
			"return_confirmed_at": gorm.Expr("CURRENT_TIMESTAMP"),
			"status":              model.BookingCompleted,
		}).Error
}

func (r *bookingRepository) CheckDateConflicts(gameID uint, startDate, endDate time.Time, excludeBookingID *uint) (bool, error) {
	query := r.db.Model(&model.Booking{}).
		Where("game_id = ? AND status IN (?)", gameID, []model.BookingStatus{
			model.BookingConfirmed, model.BookingActive,
		}).
		Where("(start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?)",
			startDate, startDate, endDate, endDate, startDate, endDate)

	if excludeBookingID != nil {
		query = query.Where("id != ?", *excludeBookingID)
	}

	var count int64
	err := query.Count(&count).Error
	return count > 0, err
}


