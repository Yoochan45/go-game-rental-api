package repository

import (
	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"gorm.io/gorm"
)

type DisputeRepository interface {
	// Basic CRUD
	Create(dispute *model.Dispute) error

	// Admin methods
	GetAllDisputes(limit, offset int) ([]*model.Dispute, error)
}

type disputeRepository struct {
	db *gorm.DB
}

func NewDisputeRepository(db *gorm.DB) DisputeRepository {
	return &disputeRepository{db: db}
}

func (r *disputeRepository) Create(dispute *model.Dispute) error {
	return r.db.Create(dispute).Error
}

func (r *disputeRepository) GetAllDisputes(limit, offset int) ([]*model.Dispute, error) {
	var disputes []*model.Dispute
	err := r.db.Preload("Booking").Preload("Reporter").Preload("Resolver").
		Order("created_at DESC").Limit(limit).Offset(offset).Find(&disputes).Error
	return disputes, err
}
