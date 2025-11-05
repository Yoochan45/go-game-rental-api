package dto

import (
	"time"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
)

type ReviewDTO struct {
	ID        uint      `json:"id"`
	BookingID uint      `json:"booking_id"`
	UserID    uint      `json:"user_id"`
	User      *UserDTO  `json:"user,omitempty"`
	GameID    uint      `json:"game_id"`
	Game      *GameDTO  `json:"game,omitempty"`
	Rating    int       `json:"rating"`
	Comment   *string   `json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment,omitempty"`
}

type UpdateReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment,omitempty"`
}


func ToReviewDTO(review *model.Review) *ReviewDTO {
	if review == nil {
		return nil
	}

	return &ReviewDTO{
		ID:        review.ID,
		BookingID: review.BookingID,
		UserID:    review.UserID,
		User:      ToUserDTO(&review.User),
		GameID:    review.GameID,
		Game:      ToGameDTO(&review.Game),
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: review.CreatedAt,
		UpdatedAt: review.UpdatedAt,
	}
}

func ToReviewDTOList(reviews []*model.Review) []*ReviewDTO {
	result := make([]*ReviewDTO, len(reviews))
	for i, review := range reviews {
		result[i] = ToReviewDTO(review)
	}
	return result
}

func FromCreateReviewRequest(req *CreateReviewRequest) *model.Review {
	return &model.Review{
		Rating:  req.Rating,
		Comment: &req.Comment,
	}
}
