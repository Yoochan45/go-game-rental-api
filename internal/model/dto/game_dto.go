package dto

import (
	"time"

	"github.com/lib/pq"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
)

type GameDTO struct {
	ID                uint                 `json:"id"`
	PartnerID         uint                 `json:"partner_id"`
	Partner           *UserDTO             `json:"partner,omitempty"`
	CategoryID        uint                 `json:"category_id"`
	Category          *CategoryDTO         `json:"category,omitempty"`
	Name              string               `json:"name"`
	Description       *string              `json:"description,omitempty"`
	Platform          *string              `json:"platform,omitempty"`
	Stock             int                  `json:"stock"`
	AvailableStock    int                  `json:"available_stock"`
	RentalPricePerDay float64              `json:"rental_price_per_day"`
	SecurityDeposit   float64              `json:"security_deposit"`
	Condition         string               `json:"condition"`
	Images            []string             `json:"images,omitempty"`
	IsActive          bool                 `json:"is_active"`
	ApprovalStatus    model.ApprovalStatus `json:"approval_status"`
	ApprovedBy        *uint                `json:"approved_by,omitempty"`
	Approver          *UserDTO             `json:"approver,omitempty"`
	ApprovedAt        *time.Time           `json:"approved_at,omitempty"`
	RejectionReason   *string              `json:"rejection_reason,omitempty"`
	AverageRating     float64              `json:"average_rating,omitempty"`
	ReviewCount       int64                `json:"review_count,omitempty"`
	CreatedAt         time.Time            `json:"created_at"`
	UpdatedAt         time.Time            `json:"updated_at"`
}

type CreateGameRequest struct {
	CategoryID        uint     `json:"category_id" validate:"required"`
	Name              string   `json:"name" validate:"required,min=2"`
	Description       string   `json:"description,omitempty"`
	Platform          string   `json:"platform,omitempty"`
	Stock             int      `json:"stock" validate:"required,min=1"`
	RentalPricePerDay float64  `json:"rental_price_per_day" validate:"required,gt=0"`
	SecurityDeposit   float64  `json:"security_deposit" validate:"gte=0"`
	Condition         string   `json:"condition" validate:"omitempty,oneof=excellent good fair"`
	Images            []string `json:"images,omitempty"`
}

type UpdateGameRequest struct {
	CategoryID        uint     `json:"category_id" validate:"required"`
	Name              string   `json:"name" validate:"required,min=2"`
	Description       string   `json:"description,omitempty"`
	Platform          string   `json:"platform,omitempty"`
	RentalPricePerDay float64  `json:"rental_price_per_day" validate:"required,gt=0"`
	SecurityDeposit   float64  `json:"security_deposit" validate:"gte=0"`
	Condition         string   `json:"condition" validate:"omitempty,oneof=excellent good fair"`
	Images            []string `json:"images,omitempty"`
}

type UpdateStockRequest struct {
	Stock int `json:"stock" validate:"required,min=0"`
}

type ApproveGameRequest struct {
	Action string `json:"action" validate:"required,oneof=approve reject"`
	Reason string `json:"reason,omitempty" validate:"required_if=Action reject"`
}

type GameListResponse struct {
	Games      []*GameDTO `json:"games"`
	TotalCount int64      `json:"total_count"`
	Page       int        `json:"page"`
	Limit      int        `json:"limit"`
}

type GameSearchRequest struct {
	Query      string  `form:"q"`
	CategoryID uint    `form:"category_id"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	Platform   string  `form:"platform"`
	Condition  string  `form:"condition"`
	Available  bool    `form:"available"`
	Page       int     `form:"page" validate:"min=1"`
	Limit      int     `form:"limit" validate:"min=1,max=100"`
}

func ToGameDTO(game *model.Game) *GameDTO {
	if game == nil {
		return nil
	}

	return &GameDTO{
		ID:                game.ID,
		PartnerID:         game.PartnerID,
		Partner:           ToUserDTO(&game.Partner),
		CategoryID:        game.CategoryID,
		Category:          ToCategoryDTO(&game.Category),
		Name:              game.Name,
		Description:       game.Description,
		Platform:          game.Platform,
		Stock:             game.Stock,
		AvailableStock:    game.AvailableStock,
		RentalPricePerDay: game.RentalPricePerDay,
		SecurityDeposit:   game.SecurityDeposit,
		Condition:         game.Condition,
		Images:            []string(game.Images),
		IsActive:          game.IsActive,
		ApprovalStatus:    game.ApprovalStatus,
		ApprovedBy:        game.ApprovedBy,
		Approver:          ToUserDTO(game.Approver),
		ApprovedAt:        game.ApprovedAt,
		RejectionReason:   game.RejectionReason,
		CreatedAt:         game.CreatedAt,
		UpdatedAt:         game.UpdatedAt,
	}
}

func ToGameDTOList(games []*model.Game) []*GameDTO {
	result := make([]*GameDTO, len(games))
	for i, game := range games {
		result[i] = ToGameDTO(game)
	}
	return result
}

func FromCreateGameRequest(req *CreateGameRequest) *model.Game {
	return &model.Game{
		CategoryID:        req.CategoryID,
		Name:              req.Name,
		Description:       &req.Description,
		Platform:          &req.Platform,
		Stock:             req.Stock,
		AvailableStock:    req.Stock,
		RentalPricePerDay: req.RentalPricePerDay,
		SecurityDeposit:   req.SecurityDeposit,
		Condition:         req.Condition,
		Images:            pq.StringArray(req.Images),
	}
}
