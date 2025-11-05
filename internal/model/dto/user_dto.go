package dto

import (
	"github.com/Yoochan45/go-game-rental-api/internal/model"
)

type UserDTO struct {
	ID       uint           `json:"id"`
	Email    string         `json:"email"`
	FullName string         `json:"full_name"`
	Phone    *string        `json:"phone,omitempty"`
	Address  *string        `json:"address,omitempty"`
	Role     model.UserRole `json:"role"`
	IsActive bool           `json:"is_active"`
}

type UpdateProfileRequest struct {
	FullName string `json:"full_name" validate:"required,min=2"`
	Phone    string `json:"phone,omitempty" validate:"omitempty,min=10"`
	Address  string `json:"address,omitempty"`
}

type UpdateUserRoleRequest struct {
	Role model.UserRole `json:"role" validate:"required,oneof=customer partner admin"`
}

func ToUserDTO(user *model.User) *UserDTO {
	if user == nil {
		return nil
	}

	return &UserDTO{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
		Phone:    user.Phone,
		Address:  user.Address,
		Role:     user.Role,
		IsActive: user.IsActive,
	}
}

func ToUserDTOList(users []*model.User) []*UserDTO {
	result := make([]*UserDTO, len(users))
	for i, user := range users {
		result[i] = ToUserDTO(user)
	}
	return result
}
