package repository

import (
	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"gorm.io/gorm"
)

type GameRepository interface {
	// Basic CRUD
	Create(game *model.Game) error
	GetByID(id uint) (*model.Game, error)
	GetByIDWithRelations(id uint) (*model.Game, error)
	Update(game *model.Game) error
	Delete(id uint) error

	// Query methods for public catalog
	GetApprovedGames(limit, offset int) ([]*model.Game, error)
	SearchGames(query string, limit, offset int) ([]*model.Game, error)

	// Partner methods
	GetGamesByPartner(partnerID uint, limit, offset int) ([]*model.Game, error)

	// Admin methods
	GetAllGames(limit, offset int) ([]*model.Game, error)
	UpdateApprovalStatus(gameID uint, status model.ApprovalStatus, approvedBy *uint, rejectionReason *string) error

	// Stock management
	CheckAvailability(gameID uint, quantity int) (bool, error)
	ReserveStock(gameID uint, quantity int) error
	ReleaseStock(gameID uint, quantity int) error
}

type gameRepository struct {
	db *gorm.DB
}

func NewGameRepository(db *gorm.DB) GameRepository {
	return &gameRepository{db: db}
}

func (r *gameRepository) Create(game *model.Game) error {
	return r.db.Create(game).Error
}

func (r *gameRepository) GetByID(id uint) (*model.Game, error) {
	var game model.Game
	err := r.db.Where("id = ?", id).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *gameRepository) GetByIDWithRelations(id uint) (*model.Game, error) {
	var game model.Game
	err := r.db.Preload("Partner").Preload("Category").Preload("Approver").
		Where("id = ?", id).First(&game).Error
	if err != nil {
		return nil, err
	}
	return &game, nil
}

func (r *gameRepository) Update(game *model.Game) error {
	return r.db.Save(game).Error
}

func (r *gameRepository) Delete(id uint) error {
	return r.db.Delete(&model.Game{}, id).Error
}

// Public catalog methods
func (r *gameRepository) GetApprovedGames(limit, offset int) ([]*model.Game, error) {
	var games []*model.Game
	err := r.db.Preload("Partner").Preload("Category").
		Where("approval_status = ? AND is_active = ?", model.ApprovalApproved, true).
		Limit(limit).Offset(offset).Find(&games).Error
	return games, err
}

func (r *gameRepository) SearchGames(query string, limit, offset int) ([]*model.Game, error) {
	var games []*model.Game
	searchQuery := "%" + query + "%"
	err := r.db.Preload("Partner").Preload("Category").
		Where("(name ILIKE ? OR description ILIKE ? OR platform ILIKE ?) AND approval_status = ? AND is_active = ?",
			searchQuery, searchQuery, searchQuery, model.ApprovalApproved, true).
		Limit(limit).Offset(offset).Find(&games).Error
	return games, err
}

// Partner methods
func (r *gameRepository) GetGamesByPartner(partnerID uint, limit, offset int) ([]*model.Game, error) {
	var games []*model.Game
	err := r.db.Preload("Category").Preload("Approver").
		Where("partner_id = ?", partnerID).
		Limit(limit).Offset(offset).Find(&games).Error
	return games, err
}

// Admin methods
func (r *gameRepository) GetAllGames(limit, offset int) ([]*model.Game, error) {
	var games []*model.Game
	err := r.db.Preload("Partner").Preload("Category").Preload("Approver").
		Limit(limit).Offset(offset).Find(&games).Error
	return games, err
}

func (r *gameRepository) UpdateApprovalStatus(gameID uint, status model.ApprovalStatus, approvedBy *uint, rejectionReason *string) error {
	updates := map[string]interface{}{
		"approval_status": status,
	}

	switch status {
	case model.ApprovalApproved:
		updates["approved_by"] = approvedBy
		updates["approved_at"] = gorm.Expr("CURRENT_TIMESTAMP")
		updates["is_active"] = true
		updates["rejection_reason"] = nil
	case model.ApprovalRejected:
		updates["rejection_reason"] = rejectionReason
		updates["is_active"] = false
	}

	return r.db.Model(&model.Game{}).Where("id = ?", gameID).Updates(updates).Error
}

// Stock management
func (r *gameRepository) CheckAvailability(gameID uint, quantity int) (bool, error) {
	var game model.Game
	err := r.db.Select("available_stock").Where("id = ?", gameID).First(&game).Error
	if err != nil {
		return false, err
	}
	return game.AvailableStock >= quantity, nil
}

func (r *gameRepository) ReserveStock(gameID uint, quantity int) error {
	return r.db.Model(&model.Game{}).Where("id = ?", gameID).
		Update("available_stock", gorm.Expr("available_stock - ?", quantity)).Error
}

func (r *gameRepository) ReleaseStock(gameID uint, quantity int) error {
	return r.db.Model(&model.Game{}).Where("id = ?", gameID).
		Update("available_stock", gorm.Expr("LEAST(available_stock + ?, stock)", quantity)).Error
}
