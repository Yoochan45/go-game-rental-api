package service

import (
	"errors"

	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"github.com/Yoochan45/go-game-rental-api/internal/repository"
)

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrInsufficientPermission = errors.New("insufficient permission")
	ErrCannotDeleteSuperAdmin = errors.New("cannot delete super admin")
	ErrCannotDeleteSelf       = errors.New("cannot delete yourself")
)

type UserService interface {
	// Public methods
	GetProfile(userID uint) (*model.User, error)
	UpdateProfile(userID uint, updateData interface{}) error

	// Auth methods
	Register(registerData interface{}) (*model.User, error)
	Login(loginData interface{}, jwtSecret string) (interface{}, error)
	RefreshToken(refreshData interface{}, jwtSecret string) (interface{}, error)
	ChangePassword(userID uint, currentPassword, newPassword string) error

	// Admin methods
	GetAllUsers(requestorRole model.UserRole, limit, offset int) ([]*model.User, error)
	GetUserDetail(requestorRole model.UserRole, userID uint) (*model.User, error)
	UpdateUserRole(requestorRole model.UserRole, userID uint, newRole model.UserRole) error
	ToggleUserStatus(requestorRole model.UserRole, userID uint) error

	// Super Admin methods
	DeleteUser(requestorID uint, requestorRole model.UserRole, targetUserID uint) error
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) GetProfile(userID uint) (*model.User, error) {
	return s.userRepo.GetByID(userID)
}

func (s *userService) UpdateProfile(userID uint, updateData interface{}) error {
	// TODO: Implement proper DTO to model conversion
	return errors.New("not implemented yet")
}

func (s *userService) Register(registerData interface{}) (*model.User, error) {
	// TODO: Implement proper registration logic with password hashing
	return nil, errors.New("not implemented yet")
}

func (s *userService) Login(loginData interface{}, jwtSecret string) (interface{}, error) {
	// TODO: Implement proper login logic with JWT generation
	return nil, errors.New("not implemented yet")
}

func (s *userService) RefreshToken(refreshData interface{}, jwtSecret string) (interface{}, error) {
	// TODO: Implement proper refresh token logic
	return nil, errors.New("not implemented yet")
}

func (s *userService) ChangePassword(userID uint, currentPassword, newPassword string) error {
	// Add password change logic here
	// Should verify current password and hash new password
	return nil
}

func (s *userService) GetAllUsers(requestorRole model.UserRole, limit, offset int) ([]*model.User, error) {
	if !s.canManageUsers(requestorRole) {
		return nil, ErrInsufficientPermission
	}

	return s.userRepo.GetAll(limit, offset)
}

func (s *userService) GetUserDetail(requestorRole model.UserRole, userID uint) (*model.User, error) {
	if !s.canManageUsers(requestorRole) {
		return nil, ErrInsufficientPermission
	}
	return s.userRepo.GetByID(userID)
}

func (s *userService) UpdateUserRole(requestorRole model.UserRole, userID uint, newRole model.UserRole) error {
	if !s.canManageUsers(requestorRole) {
		return ErrInsufficientPermission
	}
	return s.userRepo.UpdateRole(userID, newRole)
}

func (s *userService) ToggleUserStatus(requestorRole model.UserRole, userID uint) error {
	if !s.canManageUsers(requestorRole) {
		return ErrInsufficientPermission
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	return s.userRepo.UpdateActiveStatus(userID, !user.IsActive)
}

func (s *userService) DeleteUser(requestorID uint, requestorRole model.UserRole, targetUserID uint) error {
	if requestorRole != model.RoleSuperAdmin {
		return ErrInsufficientPermission
	}

	if requestorID == targetUserID {
		return ErrCannotDeleteSelf
	}

	targetUser, err := s.userRepo.GetByID(targetUserID)
	if err != nil {
		return ErrUserNotFound
	}

	// Super admin cannot be deleted
	if targetUser.Role == model.RoleSuperAdmin {
		return ErrCannotDeleteSuperAdmin
	}

	return s.userRepo.Delete(targetUserID)
}

// Helper methods
func (s *userService) canManageUsers(role model.UserRole) bool {
	return role == model.RoleAdmin || role == model.RoleSuperAdmin
}
