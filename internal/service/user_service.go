package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"regexp"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/repository"
)

// UserService defines the interface for user-related business logic
type UserService interface {
	// GetUserByID retrieves a user by its ID
	GetUserByID(ctx context.Context, id uint) (*domain.User, error)

	// GetUserByEmail retrieves a user by email
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)

	// CreateUser creates a new user
	CreateUser(ctx context.Context, user *domain.User) error

	// UpdateUser updates an existing user
	UpdateUser(ctx context.Context, user *domain.User) error

	// DeleteUser deletes a user by its ID
	DeleteUser(ctx context.Context, id uint) error

	// GetAllUsers retrieves all users with optional pagination
	GetAllUsers(ctx context.Context, page, pageSize int) ([]domain.User, int64, error)

	// UpdatePassword updates a user's password
	UpdatePassword(ctx context.Context, id uint, currentPassword, newPassword string) error

	// GetUsersByRole retrieves users by role with optional pagination
	GetUsersByRole(ctx context.Context, role string, page, pageSize int) ([]domain.User, int64, error)

	// GetUserOrders retrieves all orders for a user
	GetUserOrders(ctx context.Context, userID uint) ([]domain.Order, error)

	// GetUserCart retrieves the cart for a user
	GetUserCart(ctx context.Context, userID uint) (*domain.Cart, error)

	// AuthenticateUser authenticates a user with email and password
	AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error)
}

// UserServiceImpl implements the UserService interface
type UserServiceImpl struct {
	userRepo  repository.UserRepository
	cartRepo  repository.CartRepository
	orderRepo repository.OrderRepository
}

// NewUserService creates a new UserServiceImpl
func NewUserService(
	userRepo repository.UserRepository,
	cartRepo repository.CartRepository,
	orderRepo repository.OrderRepository,
) UserService {
	return &UserServiceImpl{
		userRepo:  userRepo,
		cartRepo:  cartRepo,
		orderRepo: orderRepo,
	}
}

// GetUserByID retrieves a user by its ID
func (s *UserServiceImpl) GetUserByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

// GetUserByEmail retrieves a user by email
func (s *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	return s.userRepo.FindByEmail(ctx, email)
}

// CreateUser creates a new user
func (s *UserServiceImpl) CreateUser(ctx context.Context, user *domain.User) error {
	// Validate email format
	if !isValidEmail(user.Email) {
		return errors.New("invalid email format")
	}

	// Check if email already exists
	existingUser, err := s.userRepo.FindByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return errors.New("email already in use")
	}

	// Validate password
	if len(user.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	// Hash the password
	hasher := sha256.New()
	hasher.Write([]byte(user.Password))
	user.Password = hex.EncodeToString(hasher.Sum(nil))

	// Set default role if not provided
	if user.Role == "" {
		user.Role = "customer"
	}

	// Create the user
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}

	// Create a cart for the user
	cart := &domain.Cart{
		UserID: user.ID,
	}
	err = s.cartRepo.Create(ctx, cart)
	if err != nil {
		return errors.New("failed to create cart for user")
	}

	return nil
}

// UpdateUser updates an existing user
func (s *UserServiceImpl) UpdateUser(ctx context.Context, user *domain.User) error {
	// Check if user exists
	existingUser, err := s.userRepo.FindByID(ctx, user.ID)
	if err != nil {
		return errors.New("user not found")
	}

	// Validate email format if changed
	if user.Email != existingUser.Email {
		if !isValidEmail(user.Email) {
			return errors.New("invalid email format")
		}

		// Check if new email already exists
		emailUser, err := s.userRepo.FindByEmail(ctx, user.Email)
		if err == nil && emailUser != nil && emailUser.ID != user.ID {
			return errors.New("email already in use")
		}
	}

	// Don't update password through this method
	user.Password = existingUser.Password

	return s.userRepo.Update(ctx, user)
}

// DeleteUser deletes a user by its ID
func (s *UserServiceImpl) DeleteUser(ctx context.Context, id uint) error {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(ctx, id)
}

// GetAllUsers retrieves all users with optional pagination
func (s *UserServiceImpl) GetAllUsers(ctx context.Context, page, pageSize int) ([]domain.User, int64, error) {
	return s.userRepo.FindAll(ctx, page, pageSize)
}

// UpdatePassword updates a user's password
func (s *UserServiceImpl) UpdatePassword(ctx context.Context, id uint, currentPassword, newPassword string) error {
	// Check if user exists
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify current password
	hasher := sha256.New()
	hasher.Write([]byte(currentPassword))
	hashedCurrentPassword := hex.EncodeToString(hasher.Sum(nil))
	if user.Password != hashedCurrentPassword {
		return errors.New("current password is incorrect")
	}

	// Validate new password
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters long")
	}

	// Hash the new password
	hasher = sha256.New()
	hasher.Write([]byte(newPassword))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))

	return s.userRepo.UpdatePassword(ctx, id, hashedPassword)
}

// GetUsersByRole retrieves users by role with optional pagination
func (s *UserServiceImpl) GetUsersByRole(ctx context.Context, role string, page, pageSize int) ([]domain.User, int64, error) {
	return s.userRepo.FindByRole(ctx, role, page, pageSize)
}

// GetUserOrders retrieves all orders for a user
func (s *UserServiceImpl) GetUserOrders(ctx context.Context, userID uint) ([]domain.Order, error) {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.userRepo.GetUserOrders(ctx, userID)
}

// GetUserCart retrieves the cart for a user
func (s *UserServiceImpl) GetUserCart(ctx context.Context, userID uint) (*domain.Cart, error) {
	// Check if user exists
	_, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return s.userRepo.GetUserCart(ctx, userID)
}

// AuthenticateUser authenticates a user with email and password
func (s *UserServiceImpl) AuthenticateUser(ctx context.Context, email, password string) (*domain.User, error) {
	// Get user by email
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	hasher := sha256.New()
	hasher.Write([]byte(password))
	hashedPassword := hex.EncodeToString(hasher.Sum(nil))
	if user.Password != hashedPassword {
		return nil, errors.New("invalid email or password")
	}

	// Don't return the password
	user.Password = ""

	return user, nil
}

// isValidEmail validates email format
func isValidEmail(email string) bool {
	pattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(email)
}
