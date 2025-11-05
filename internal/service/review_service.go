package service

import (
	"errors"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"github.com/Yoochan45/go-game-rental-api/internal/repository"
)

var (
	ErrReviewNotFound            = errors.New("review not found")
	ErrReviewAlreadyExists       = errors.New("review already exists for this booking")
	ErrReviewNotOwned            = errors.New("you don't own this review")
	ErrReviewBookingNotCompleted = errors.New("can only review completed bookings")
)

type ReviewService interface {
	// Customer methods
	CreateReview(userID uint, bookingID uint, reviewData *model.Review) error
	UpdateReview(userID uint, reviewID uint, updateData *model.Review) error
	DeleteReview(userID uint, reviewID uint) error
	GetUserReviews(userID uint, limit, offset int) ([]*model.Review, error)
	GetMyReviewForBooking(userID uint, bookingID uint) (*model.Review, error)

	// Public methods
	GetGameReviews(gameID uint, limit, offset int) ([]*model.Review, error)
	GetGameRating(gameID uint) (float64, int64, error)

	// Partner methods
	GetPartnerReviews(partnerID uint, limit, offset int) ([]*model.Review, error)
	GetPartnerRating(partnerID uint) (float64, error)

	// Admin methods
	DeleteReviewByAdmin(requestorRole model.UserRole, reviewID uint) error
}

type reviewService struct {
	reviewRepo  repository.ReviewRepository
	bookingRepo repository.BookingRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, bookingRepo repository.BookingRepository) ReviewService {
	return &reviewService{
		reviewRepo:  reviewRepo,
		bookingRepo: bookingRepo,
	}
}

func (s *reviewService) CreateReview(userID uint, bookingID uint, reviewData *model.Review) error {
	// Validate booking exists and belongs to user
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return ErrBookingNotFound
	}

	if booking.UserID != userID {
		return ErrBookingNotOwned
	}

	// Can only review completed bookings
	if booking.Status != model.BookingCompleted {
		return ErrReviewBookingNotCompleted
	}

	// Check if review already exists
	existingReview, _ := s.reviewRepo.GetByBookingID(bookingID)
	if existingReview != nil {
		return ErrReviewAlreadyExists
	}

	// Set review details
	reviewData.BookingID = bookingID
	reviewData.UserID = userID
	reviewData.GameID = booking.GameID

	return s.reviewRepo.Create(reviewData)
}

func (s *reviewService) UpdateReview(userID uint, reviewID uint, updateData *model.Review) error {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return ErrReviewNotFound
	}

	if review.UserID != userID {
		return ErrReviewNotOwned
	}

	// Update allowed fields
	review.Rating = updateData.Rating
	review.Comment = updateData.Comment

	return s.reviewRepo.Update(review)
}

func (s *reviewService) DeleteReview(userID uint, reviewID uint) error {
	review, err := s.reviewRepo.GetByID(reviewID)
	if err != nil {
		return ErrReviewNotFound
	}

	if review.UserID != userID {
		return ErrReviewNotOwned
	}

	return s.reviewRepo.Delete(reviewID)
}

func (s *reviewService) GetUserReviews(userID uint, limit, offset int) ([]*model.Review, error) {
	return s.reviewRepo.GetUserReviews(userID, limit, offset)
}

func (s *reviewService) GetMyReviewForBooking(userID uint, bookingID uint) (*model.Review, error) {
	// Validate booking belongs to user
	booking, err := s.bookingRepo.GetByID(bookingID)
	if err != nil {
		return nil, ErrBookingNotFound
	}

	if booking.UserID != userID {
		return nil, ErrBookingNotOwned
	}

	return s.reviewRepo.GetByBookingID(bookingID)
}

func (s *reviewService) GetGameReviews(gameID uint, limit, offset int) ([]*model.Review, error) {
	return s.reviewRepo.GetGameReviews(gameID, limit, offset)
}

func (s *reviewService) GetGameRating(gameID uint) (float64, int64, error) {
	avgRating, err := s.reviewRepo.GetGameAverageRating(gameID)
	if err != nil {
		return 0, 0, err
	}

	count, err := s.reviewRepo.CountGameReviews(gameID)
	if err != nil {
		return 0, 0, err
	}

	return avgRating, count, nil
}

func (s *reviewService) GetPartnerReviews(partnerID uint, limit, offset int) ([]*model.Review, error) {
	return s.reviewRepo.GetPartnerReviews(partnerID, limit, offset)
}

func (s *reviewService) GetPartnerRating(partnerID uint) (float64, error) {
	return s.reviewRepo.GetPartnerAverageRating(partnerID)
}

func (s *reviewService) DeleteReviewByAdmin(requestorRole model.UserRole, reviewID uint) error {
	if !s.canManageReviews(requestorRole) {
		return ErrInsufficientPermission
	}

	return s.reviewRepo.Delete(reviewID)
}

func (s *reviewService) canManageReviews(role model.UserRole) bool {
	return role == model.RoleAdmin || role == model.RoleSuperAdmin
}
