package service

import (
	"errors"
	"time"

	"github.com/Yoochan45/go-game-rental-api/internal/dto"
	"github.com/Yoochan45/go-game-rental-api/internal/model"
	"github.com/Yoochan45/go-game-rental-api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
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

	// Admin methods
	GetAllUsers(requestorRole model.UserRole, limit, offset int) ([]*model.User, int64, error)
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
	req := updateData.(*dto.UpdateProfileRequest)

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	user.FullName = req.FullName
	user.Phone = &req.Phone
	user.Address = &req.Address

	return s.userRepo.Update(user)
}

func (s *userService) Register(registerData interface{}) (*model.User, error) {
	req := registerData.(*dto.RegisterRequest)

	// Check if user exists
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:    req.Email,
		Password: string(hashed),
		FullName: req.FullName,
		Phone:    &req.Phone,
		Address:  &req.Address,
		Role:     model.RoleCustomer,
	}

	err = s.userRepo.Create(user)
	return user, err
}

func (s *userService) Login(loginData interface{}, jwtSecret string) (interface{}, error) {
	req := loginData.(*dto.LoginRequest)

	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	accessToken, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken: accessToken,
		User:        user,
		ExpiresAt:   time.Now().Add(24 * time.Hour),
	}, nil
}

func (s *userService) RefreshToken(refreshData interface{}, jwtSecret string) (interface{}, error) {
	// For now, just return error - implement later if needed
	return nil, errors.New("refresh token not implemented")
}

func (s *userService) GetAllUsers(requestorRole model.UserRole, limit, offset int) ([]*model.User, int64, error) {
	if !s.canManageUsers(requestorRole) {
		return nil, 0, ErrInsufficientPermission
	}

	users, err := s.userRepo.GetAll(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.userRepo.Count()
	return users, count, err
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
