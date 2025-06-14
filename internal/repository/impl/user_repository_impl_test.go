package impl_test

import (
	"context"
	"errors"
	"testing"

	"awesomeEcommerce/internal/domain"
	"awesomeEcommerce/internal/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockUserRepository is a mock implementation of the UserRepository interface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Create(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Update(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) FindAll(ctx context.Context, page, pageSize int) ([]domain.User, int64, error) {
	args := m.Called(ctx, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) UpdatePassword(ctx context.Context, id uint, hashedPassword string) error {
	args := m.Called(ctx, id, hashedPassword)
	return args.Error(0)
}

func (m *MockUserRepository) FindByRole(ctx context.Context, role string, page, pageSize int) ([]domain.User, int64, error) {
	args := m.Called(ctx, role, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]domain.User), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserRepository) GetUserOrders(ctx context.Context, userID uint) ([]domain.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.Order), args.Error(1)
}

func (m *MockUserRepository) GetUserCart(ctx context.Context, userID uint) (*domain.Cart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cart), args.Error(1)
}

// UserRepositoryTestSuite is a test suite for UserRepository
type UserRepositoryTestSuite struct {
	suite.Suite
	mockRepo repository.UserRepository
	ctx      context.Context
}

// SetupTest sets up the test suite
func (s *UserRepositoryTestSuite) SetupTest() {
	s.mockRepo = new(MockUserRepository)
	s.ctx = context.Background()
}

// TestFindByID tests the FindByID method
func (s *UserRepositoryTestSuite) TestFindByID() {
	mockRepo := s.mockRepo.(*MockUserRepository)

	s.Run("Success", func() {
		// Test case: Successfully find a user by ID
		userID := uint(1)
		expectedUser := &domain.User{
			ID:        userID,
			Email:     "test@example.com",
			FirstName: "Test",
			LastName:  "User",
			Role:      "customer",
		}

		mockRepo.On("FindByID", s.ctx, userID).Return(expectedUser, nil).Once()

		// Execute
		user, err := s.mockRepo.FindByID(s.ctx, userID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedUser, user)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: User not found
		userID := uint(999)
		expectedError := errors.New("user not found")

		mockRepo.On("FindByID", s.ctx, userID).Return(nil, expectedError).Once()

		// Execute
		user, err := s.mockRepo.FindByID(s.ctx, userID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), user)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByEmail tests the FindByEmail method
func (s *UserRepositoryTestSuite) TestFindByEmail() {
	mockRepo := s.mockRepo.(*MockUserRepository)

	s.Run("Success", func() {
		// Test case: Successfully find a user by email
		email := "test@example.com"
		expectedUser := &domain.User{
			ID:        1,
			Email:     email,
			FirstName: "Test",
			LastName:  "User",
			Role:      "customer",
		}

		mockRepo.On("FindByEmail", s.ctx, email).Return(expectedUser, nil).Once()

		// Execute
		user, err := s.mockRepo.FindByEmail(s.ctx, email)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedUser, user)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: User not found
		email := "nonexistent@example.com"
		expectedError := errors.New("user not found")

		mockRepo.On("FindByEmail", s.ctx, email).Return(nil, expectedError).Once()

		// Execute
		user, err := s.mockRepo.FindByEmail(s.ctx, email)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), user)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestCreate tests the Create method
func (s *UserRepositoryTestSuite) TestCreate() {
	mockRepo := s.mockRepo.(*MockUserRepository)

	s.Run("Success", func() {
		// Test case: Successfully create a user
		user := &domain.User{
			Email:     "new@example.com",
			Password:  "hashedpassword",
			FirstName: "New",
			LastName:  "User",
			Role:      "customer",
		}

		mockRepo.On("Create", s.ctx, user).Return(nil).Once()

		// Execute
		err := s.mockRepo.Create(s.ctx, user)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: Database error
		user := &domain.User{
			Email:     "new@example.com",
			Password:  "hashedpassword",
			FirstName: "New",
			LastName:  "User",
			Role:      "customer",
		}
		expectedError := errors.New("database error")

		mockRepo.On("Create", s.ctx, user).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Create(s.ctx, user)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdate tests the Update method
func (s *UserRepositoryTestSuite) TestUpdate() {
	mockRepo := s.mockRepo.(*MockUserRepository)

	s.Run("Success", func() {
		// Test case: Successfully update a user
		user := &domain.User{
			ID:        1,
			Email:     "updated@example.com",
			FirstName: "Updated",
			LastName:  "User",
			Role:      "customer",
		}

		mockRepo.On("Update", s.ctx, user).Return(nil).Once()

		// Execute
		err := s.mockRepo.Update(s.ctx, user)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - User Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: User not found
		user := &domain.User{
			ID:        999,
			Email:     "nonexistent@example.com",
			FirstName: "Nonexistent",
			LastName:  "User",
			Role:      "customer",
		}
		expectedError := errors.New("user not found")

		mockRepo.On("Update", s.ctx, user).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Update(s.ctx, user)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestDelete tests the Delete method
func (s *UserRepositoryTestSuite) TestDelete() {
	mockRepo := s.mockRepo.(*MockUserRepository)

	s.Run("Success", func() {
		// Test case: Successfully delete a user
		userID := uint(1)

		mockRepo.On("Delete", s.ctx, userID).Return(nil).Once()

		// Execute
		err := s.mockRepo.Delete(s.ctx, userID)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - User Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: User not found
		userID := uint(999)
		expectedError := errors.New("user not found")

		mockRepo.On("Delete", s.ctx, userID).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Delete(s.ctx, userID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindAll tests the FindAll method
func (s *UserRepositoryTestSuite) TestFindAll() {
	mockRepo := s.mockRepo.(*MockUserRepository)

	s.Run("Success", func() {
		// Test case: Successfully find all users
		page := 1
		pageSize := 10
		var total int64 = 2
		expectedUsers := []domain.User{
			{
				ID:        1,
				Email:     "user1@example.com",
				FirstName: "User",
				LastName:  "One",
				Role:      "customer",
			},
			{
				ID:        2,
				Email:     "user2@example.com",
				FirstName: "User",
				LastName:  "Two",
				Role:      "admin",
			},
		}

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(expectedUsers, total, nil).Once()

		// Execute
		users, count, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedUsers, users)
		assert.Equal(s.T(), total, count)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Empty Result", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: No users found
		page := 1
		pageSize := 10
		var total int64 = 0
		var expectedUsers []domain.User

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(expectedUsers, total, nil).Once()

		// Execute
		users, count, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), users)
		assert.Equal(s.T(), total, count)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: Database error
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindAll", s.ctx, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		users, count, err := s.mockRepo.FindAll(s.ctx, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), users)
		assert.Equal(s.T(), int64(0), count)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdatePassword tests the UpdatePassword method
func (s *UserRepositoryTestSuite) TestUpdatePassword() {
	mockRepo := s.mockRepo.(*MockUserRepository)

	s.Run("Success", func() {
		// Test case: Successfully update a user's password
		userID := uint(1)
		hashedPassword := "newhashpassword"

		mockRepo.On("UpdatePassword", s.ctx, userID, hashedPassword).Return(nil).Once()

		// Execute
		err := s.mockRepo.UpdatePassword(s.ctx, userID, hashedPassword)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - User Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: User not found
		userID := uint(999)
		hashedPassword := "newhashpassword"
		expectedError := errors.New("user not found")

		mockRepo.On("UpdatePassword", s.ctx, userID, hashedPassword).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.UpdatePassword(s.ctx, userID, hashedPassword)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByRole tests the FindByRole method
func (s *UserRepositoryTestSuite) TestFindByRole() {
	mockRepo := s.mockRepo.(*MockUserRepository)

	s.Run("Success", func() {
		// Test case: Successfully find users by role
		role := "admin"
		page := 1
		pageSize := 10
		var total int64 = 1
		expectedUsers := []domain.User{
			{
				ID:        2,
				Email:     "admin@example.com",
				FirstName: "Admin",
				LastName:  "User",
				Role:      "admin",
			},
		}

		mockRepo.On("FindByRole", s.ctx, role, page, pageSize).Return(expectedUsers, total, nil).Once()

		// Execute
		users, count, err := s.mockRepo.FindByRole(s.ctx, role, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedUsers, users)
		assert.Equal(s.T(), total, count)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Empty Result", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: No users found with the specified role
		role := "manager"
		page := 1
		pageSize := 10
		var total int64 = 0
		var expectedUsers []domain.User

		mockRepo.On("FindByRole", s.ctx, role, page, pageSize).Return(expectedUsers, total, nil).Once()

		// Execute
		users, count, err := s.mockRepo.FindByRole(s.ctx, role, page, pageSize)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), users)
		assert.Equal(s.T(), total, count)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: Database error
		role := "admin"
		page := 1
		pageSize := 10
		expectedError := errors.New("database error")

		mockRepo.On("FindByRole", s.ctx, role, page, pageSize).Return(nil, int64(0), expectedError).Once()

		// Execute
		users, count, err := s.mockRepo.FindByRole(s.ctx, role, page, pageSize)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), users)
		assert.Equal(s.T(), int64(0), count)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestGetUserOrders tests the GetUserOrders method
func (s *UserRepositoryTestSuite) TestGetUserOrders() {
	mockRepo := s.mockRepo.(*MockUserRepository)

	s.Run("Success", func() {
		// Test case: Successfully get user orders
		userID := uint(1)
		expectedOrders := []domain.Order{
			{
				ID:     1,
				UserID: userID,
				Status: "completed",
				Items: []domain.OrderItem{
					{
						ID:        1,
						OrderID:   1,
						ProductID: 2,
						Quantity:  3,
						Price:     10.99,
					},
				},
			},
			{
				ID:     2,
				UserID: userID,
				Status: "processing",
				Items: []domain.OrderItem{
					{
						ID:        2,
						OrderID:   2,
						ProductID: 3,
						Quantity:  1,
						Price:     5.99,
					},
				},
			},
		}

		mockRepo.On("GetUserOrders", s.ctx, userID).Return(expectedOrders, nil).Once()

		// Execute
		orders, err := s.mockRepo.GetUserOrders(s.ctx, userID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedOrders, orders)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - No Orders", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: User has no orders
		userID := uint(1)
		var expectedOrders []domain.Order

		mockRepo.On("GetUserOrders", s.ctx, userID).Return(expectedOrders, nil).Once()

		// Execute
		orders, err := s.mockRepo.GetUserOrders(s.ctx, userID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), orders)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - User Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: User not found
		userID := uint(999)
		expectedError := errors.New("user not found")

		mockRepo.On("GetUserOrders", s.ctx, userID).Return(nil, expectedError).Once()

		// Execute
		orders, err := s.mockRepo.GetUserOrders(s.ctx, userID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), orders)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestGetUserCart tests the GetUserCart method
func (s *UserRepositoryTestSuite) TestGetUserCart() {
	mockRepo := s.mockRepo.(*MockUserRepository)

	s.Run("Success", func() {
		// Test case: Successfully get user cart
		userID := uint(1)
		expectedCart := &domain.Cart{
			ID:     1,
			UserID: userID,
			Items: []domain.CartItem{
				{
					ID:        1,
					CartID:    1,
					ProductID: 2,
					Quantity:  3,
					Product: domain.Product{
						ID:    2,
						Name:  "Test Product",
						Price: 10.99,
					},
				},
			},
		}

		mockRepo.On("GetUserCart", s.ctx, userID).Return(expectedCart, nil).Once()

		// Execute
		cart, err := s.mockRepo.GetUserCart(s.ctx, userID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedCart, cart)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - User Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockUserRepository)

		// Test case: User not found
		userID := uint(999)
		expectedError := errors.New("user not found")

		mockRepo.On("GetUserCart", s.ctx, userID).Return(nil, expectedError).Once()

		// Execute
		cart, err := s.mockRepo.GetUserCart(s.ctx, userID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), cart)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUserRepositorySuite runs the test suite
func TestUserRepositorySuite(t *testing.T) {
	suite.Run(t, new(UserRepositoryTestSuite))
}
