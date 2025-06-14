package service_test

import (
	"context"
	"errors"
	"testing"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUserByID(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()
	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedUser := &domain.User{
			ID:        userID,
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(expectedUser, nil).Once()

		// Execute
		user, err := userService.GetUserByID(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		// Execute
		user, err := userService.GetUserByID(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGetUserByEmail(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()
	email := "test@example.com"

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedUser := &domain.User{
			ID:        1,
			FirstName: "Test",
			LastName:  "User",
			Email:     email,
			Password:  "hashedpassword",
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByEmail", ctx, email).Return(expectedUser, nil).Once()

		// Execute
		user, err := userService.GetUserByEmail(ctx, email)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Expectations
		mockUserRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("user not found")).Once()

		// Execute
		user, err := userService.GetUserByEmail(ctx, email)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "user not found")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestCreateUser(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Test data
		user := &domain.User{
			FirstName: "New",
			LastName:  "User",
			Email:     "new@example.com",
			Password:  "password123",
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByEmail", ctx, user.Email).Return(nil, errors.New("not found")).Once()
		mockUserRepo.On("Create", ctx, user).Return(nil).Once()
		mockCartRepo.On("Create", ctx, &domain.Cart{UserID: user.ID}).Return(nil).Once()

		// Execute
		err := userService.CreateUser(ctx, user)

		// Assert
		assert.NoError(t, err)
		assert.NotEqual(t, "password123", user.Password, "Password should be hashed")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("Invalid Email", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		user := &domain.User{
			FirstName: "New",
			LastName:  "User",
			Email:     "invalid-email",
			Password:  "password123",
			Role:      "customer",
		}

		// Execute
		err := userService.CreateUser(ctx, user)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
		mockUserRepo.AssertNotCalled(t, "Create")
		mockCartRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Email Already Exists", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		existingUser := &domain.User{
			ID:        1,
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}

		user := &domain.User{
			FirstName: "New",
			LastName:  "User",
			Email:     "existing@example.com",
			Password:  "password123",
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByEmail", ctx, user.Email).Return(existingUser, nil).Once()

		// Execute
		err := userService.CreateUser(ctx, user)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email already in use")
		mockUserRepo.AssertNotCalled(t, "Create")
		mockCartRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Password Too Short", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		user := &domain.User{
			FirstName: "New",
			LastName:  "User",
			Email:     "new@example.com",
			Password:  "short",
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByEmail", ctx, user.Email).Return(nil, errors.New("not found")).Once()

		// Execute
		err := userService.CreateUser(ctx, user)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password must be at least 8 characters long")
		mockUserRepo.AssertNotCalled(t, "Create")
		mockCartRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Cart Creation Error", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		user := &domain.User{
			FirstName: "New",
			LastName:  "User",
			Email:     "new@example.com",
			Password:  "password123",
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByEmail", ctx, user.Email).Return(nil, errors.New("not found")).Once()
		mockUserRepo.On("Create", ctx, user).Return(nil).Once()
		mockCartRepo.On("Create", ctx, &domain.Cart{UserID: user.ID}).Return(errors.New("cart creation error")).Once()

		// Execute
		err := userService.CreateUser(ctx, user)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create cart for user")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
	})
}

func TestUpdateUser(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		// Test data
		existingUser := &domain.User{
			ID:        1,
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}

		updatedUser := &domain.User{
			ID:        1,
			FirstName: "Updated",
			LastName:  "User",
			Email:     "existing@example.com",
			Password:  "newpassword", // Should be ignored
			Role:      "admin",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, updatedUser.ID).Return(existingUser, nil).Once()
		mockUserRepo.On("Update", ctx, updatedUser).Return(nil).Once()

		// Execute
		err := userService.UpdateUser(ctx, updatedUser)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "hashedpassword", updatedUser.Password, "Password should not be updated")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		user := &domain.User{
			ID:        999,
			FirstName: "Non-existent",
			LastName:  "User",
			Email:     "nonexistent@example.com",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, user.ID).Return(nil, errors.New("user not found")).Once()

		// Execute
		err := userService.UpdateUser(ctx, user)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		mockUserRepo.AssertNotCalled(t, "Update")
	})

	t.Run("Invalid Email", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		existingUser := &domain.User{
			ID:        1,
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}

		updatedUser := &domain.User{
			ID:        1,
			FirstName: "Updated",
			LastName:  "User",
			Email:     "invalid-email",
			Role:      "admin",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, updatedUser.ID).Return(existingUser, nil).Once()

		// Execute
		err := userService.UpdateUser(ctx, updatedUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid email format")
		mockUserRepo.AssertNotCalled(t, "Update")
	})

	t.Run("Email Already In Use", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		existingUser := &domain.User{
			ID:        1,
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}

		anotherUser := &domain.User{
			ID:        2,
			FirstName: "Another",
			LastName:  "User",
			Email:     "another@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}

		updatedUser := &domain.User{
			ID:        1,
			FirstName: "Updated",
			LastName:  "User",
			Email:     "another@example.com", // Trying to use another user's email
			Role:      "admin",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, updatedUser.ID).Return(existingUser, nil).Once()
		mockUserRepo.On("FindByEmail", ctx, updatedUser.Email).Return(anotherUser, nil).Once()

		// Execute
		err := userService.UpdateUser(ctx, updatedUser)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email already in use")
		mockUserRepo.AssertNotCalled(t, "Update")
	})
}

func TestDeleteUser(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()
	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("Delete", ctx, userID).Return(nil).Once()

		// Execute
		err := userService.DeleteUser(ctx, userID)

		// Assert
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		// Execute
		err := userService.DeleteUser(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		mockUserRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("Delete Error", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Existing",
			LastName:  "User",
			Email:     "existing@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("Delete", ctx, userID).Return(errors.New("delete error")).Once()

		// Execute
		err := userService.DeleteUser(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete error")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGetAllUsers(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()
	page := 1
	pageSize := 10

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedUsers := []domain.User{
			{
				ID:        1,
				FirstName: "User",
				LastName:  "One",
				Email:     "user1@example.com",
				Password:  "hashedpassword1",
				Role:      "customer",
			},
			{
				ID:        2,
				FirstName: "User",
				LastName:  "Two",
				Email:     "user2@example.com",
				Password:  "hashedpassword2",
				Role:      "admin",
			},
		}
		expectedTotal := int64(2)

		// Expectations
		mockUserRepo.On("FindAll", ctx, page, pageSize).Return(expectedUsers, expectedTotal, nil).Once()

		// Execute
		users, total, err := userService.GetAllUsers(ctx, page, pageSize)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, expectedTotal, total)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Expectations
		mockUserRepo.On("FindAll", ctx, page, pageSize).Return([]domain.User{}, int64(0), errors.New("database error")).Once()

		// Execute
		users, total, err := userService.GetAllUsers(ctx, page, pageSize)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, users)
		assert.Equal(t, int64(0), total)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestAuthenticateUser(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()
	email := "test@example.com"
	password := "password123"

	t.Run("Success", func(t *testing.T) {
		// Test data
		// Note: The password hash is for "password123" using SHA-256
		hashedPassword := "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f"
		expectedUser := &domain.User{
			ID:        1,
			FirstName: "Test",
			LastName:  "User",
			Email:     email,
			Password:  hashedPassword,
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByEmail", ctx, email).Return(expectedUser, nil).Once()

		// Execute
		user, err := userService.AuthenticateUser(ctx, email, password)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		assert.Empty(t, user.Password, "Password should be cleared for security")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Expectations
		mockUserRepo.On("FindByEmail", ctx, email).Return(nil, errors.New("user not found")).Once()

		// Execute
		user, err := userService.AuthenticateUser(ctx, email, password)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid email or password")
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Incorrect Password", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		// Note: This is a hash for a different password
		hashedPassword := "wronghash"
		existingUser := &domain.User{
			ID:        1,
			FirstName: "Test",
			LastName:  "User",
			Email:     email,
			Password:  hashedPassword,
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByEmail", ctx, email).Return(existingUser, nil).Once()

		// Execute
		user, err := userService.AuthenticateUser(ctx, email, password)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "invalid email or password")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUpdatePassword(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()
	userID := uint(1)
	currentPassword := "password123"
	newPassword := "newpassword123"

	t.Run("Success", func(t *testing.T) {
		// Test data
		// Note: The password hash is for "password123" using SHA-256
		hashedPassword := "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f"
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			Password:  hashedPassword,
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("UpdatePassword", ctx, userID, mock.AnythingOfType("string")).Return(nil).Once()

		// Execute
		err := userService.UpdatePassword(ctx, userID, currentPassword, newPassword)

		// Assert
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		// Execute
		err := userService.UpdatePassword(ctx, userID, currentPassword, newPassword)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		mockUserRepo.AssertNotCalled(t, "UpdatePassword")
	})

	t.Run("Incorrect Current Password", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		// Note: This is a hash for a different password
		hashedPassword := "wronghash"
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			Password:  hashedPassword,
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()

		// Execute
		err := userService.UpdatePassword(ctx, userID, currentPassword, newPassword)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "current password is incorrect")
		mockUserRepo.AssertNotCalled(t, "UpdatePassword")
	})

	t.Run("New Password Too Short", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		// Note: The password hash is for "password123" using SHA-256
		hashedPassword := "ef92b778bafe771e89245b89ecbc08a44a4e166c06659911881f383d4473e94f"
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			Password:  hashedPassword,
			Role:      "customer",
		}
		shortNewPassword := "short"

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()

		// Execute
		err := userService.UpdatePassword(ctx, userID, currentPassword, shortNewPassword)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "new password must be at least 8 characters long")
		mockUserRepo.AssertNotCalled(t, "UpdatePassword")
	})
}

func TestGetUsersByRole(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()
	role := "admin"
	page := 1
	pageSize := 10

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedUsers := []domain.User{
			{
				ID:        1,
				FirstName: "Admin",
				LastName:  "One",
				Email:     "admin1@example.com",
				Password:  "hashedpassword1",
				Role:      role,
			},
			{
				ID:        2,
				FirstName: "Admin",
				LastName:  "Two",
				Email:     "admin2@example.com",
				Password:  "hashedpassword2",
				Role:      role,
			},
		}
		expectedTotal := int64(2)

		// Expectations
		mockUserRepo.On("FindByRole", ctx, role, page, pageSize).Return(expectedUsers, expectedTotal, nil).Once()

		// Execute
		users, total, err := userService.GetUsersByRole(ctx, role, page, pageSize)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, expectedTotal, total)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Expectations
		mockUserRepo.On("FindByRole", ctx, role, page, pageSize).Return([]domain.User{}, int64(0), errors.New("database error")).Once()

		// Execute
		users, total, err := userService.GetUsersByRole(ctx, role, page, pageSize)

		// Assert
		assert.Error(t, err)
		assert.Empty(t, users)
		assert.Equal(t, int64(0), total)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGetUserOrders(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()
	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}
		expectedOrders := []domain.Order{
			{
				ID:          1,
				UserID:      userID,
				TotalAmount: 99.99,
				Status:      domain.OrderStatusPending,
			},
			{
				ID:          2,
				UserID:      userID,
				TotalAmount: 149.99,
				Status:      domain.OrderStatusDelivered,
			},
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("GetUserOrders", ctx, userID).Return(expectedOrders, nil).Once()

		// Execute
		orders, err := userService.GetUserOrders(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedOrders, orders)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		// Execute
		orders, err := userService.GetUserOrders(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Contains(t, err.Error(), "user not found")
		mockUserRepo.AssertNotCalled(t, "GetUserOrders")
	})

	t.Run("Database Error", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("GetUserOrders", ctx, userID).Return(nil, errors.New("database error")).Once()

		// Execute
		orders, err := userService.GetUserOrders(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, orders)
		assert.Contains(t, err.Error(), "database error")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestGetUserCart(t *testing.T) {
	// Setup
	mockUserRepo := new(MockUserRepository)
	mockCartRepo := new(MockCartRepository)
	mockOrderRepo := new(MockOrderRepository)
	userService := service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)
	ctx := context.Background()
	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}
		expectedCart := &domain.Cart{
			ID:     1,
			UserID: userID,
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("GetUserCart", ctx, userID).Return(expectedCart, nil).Once()

		// Execute
		cart, err := userService.GetUserCart(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCart, cart)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		// Execute
		cart, err := userService.GetUserCart(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cart)
		assert.Contains(t, err.Error(), "user not found")
		mockUserRepo.AssertNotCalled(t, "GetUserCart")
	})

	t.Run("Database Error", func(t *testing.T) {
		// Reset mocks
		mockUserRepo = new(MockUserRepository)
		mockCartRepo = new(MockCartRepository)
		mockOrderRepo = new(MockOrderRepository)
		userService = service.NewUserService(mockUserRepo, mockCartRepo, mockOrderRepo)

		// Test data
		existingUser := &domain.User{
			ID:        userID,
			FirstName: "Test",
			LastName:  "User",
			Email:     "test@example.com",
			Password:  "hashedpassword",
			Role:      "customer",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(existingUser, nil).Once()
		mockUserRepo.On("GetUserCart", ctx, userID).Return(nil, errors.New("database error")).Once()

		// Execute
		cart, err := userService.GetUserCart(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cart)
		assert.Contains(t, err.Error(), "database error")
		mockUserRepo.AssertExpectations(t)
	})
}
