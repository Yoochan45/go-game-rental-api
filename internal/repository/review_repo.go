package repository

import (
	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"gorm.io/gorm"
)

type ReviewRepository interface {
	// Basic CRUD
	Create(review *model.Review) error
	GetByID(id uint) (*model.Review, error)
	GetByIDWithRelations(id uint) (*model.Review, error)
	Update(review *model.Review) error
	Delete(id uint) error

	// Query methods
	GetByBookingID(bookingID uint) (*model.Review, error)
	GetGameReviews(gameID uint, limit, offset int) ([]*model.Review, error)
	GetUserReviews(userID uint, limit, offset int) ([]*model.Review, error)
	GetPartnerReviews(partnerID uint, limit, offset int) ([]*model.Review, error)

	// Statistics
	GetGameAverageRating(gameID uint) (float64, error)
	CountGameReviews(gameID uint) (int64, error)
	GetPartnerAverageRating(partnerID uint) (float64, error)
}

type reviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) ReviewRepository {
	return &reviewRepository{db: db}
}

func (r *reviewRepository) Create(review *model.Review) error {
	return r.db.Create(review).Error
}

func (r *reviewRepository) GetByID(id uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Where("id = ?", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetByIDWithRelations(id uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Preload("User").Preload("Game").Preload("Booking").
		Where("id = ?", id).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) Update(review *model.Review) error {
	return r.db.Save(review).Error
}

func (r *reviewRepository) Delete(id uint) error {
	return r.db.Delete(&model.Review{}, id).Error
}

func (r *reviewRepository) GetByBookingID(bookingID uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Where("booking_id = ?", bookingID).First(&review).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

func (r *reviewRepository) GetGameReviews(gameID uint, limit, offset int) ([]*model.Review, error) {
	var reviews []*model.Review
	err := r.db.Preload("User").Preload("Booking").
		Where("game_id = ?", gameID).Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetUserReviews(userID uint, limit, offset int) ([]*model.Review, error) {
	var reviews []*model.Review
	err := r.db.Preload("Game").Preload("Booking").
		Where("user_id = ?", userID).Order("created_at DESC").
		Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetPartnerReviews(partnerID uint, limit, offset int) ([]*model.Review, error) {
	var reviews []*model.Review
	err := r.db.Preload("User").Preload("Game").Preload("Booking").
		Joins("JOIN games ON reviews.game_id = games.id").
		Where("games.partner_id = ?", partnerID).Order("reviews.created_at DESC").
		Limit(limit).Offset(offset).Find(&reviews).Error
	return reviews, err
}

func (r *reviewRepository) GetGameAverageRating(gameID uint) (float64, error) {
	var avgRating float64
	err := r.db.Model(&model.Review{}).Where("game_id = ?", gameID).
		Select("COALESCE(AVG(rating), 0)").Scan(&avgRating).Error
	return avgRating, err
}

func (r *reviewRepository) CountGameReviews(gameID uint) (int64, error) {
	var count int64
	err := r.db.Model(&model.Review{}).Where("game_id = ?", gameID).Count(&count).Error
	return count, err
}

func (r *reviewRepository) GetPartnerAverageRating(partnerID uint) (float64, error) {
	var avgRating float64
	err := r.db.Model(&model.Review{}).
		Joins("JOIN games ON reviews.game_id = games.id").
		Where("games.partner_id = ?", partnerID).
		Select("COALESCE(AVG(reviews.rating), 0)").Scan(&avgRating).Error
	return avgRating, err
}
