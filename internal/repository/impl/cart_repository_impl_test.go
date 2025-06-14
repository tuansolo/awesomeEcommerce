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

// MockCartRepository is a mock implementation of the CartRepository interface
type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) FindByID(ctx context.Context, id uint) (*domain.Cart, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cart), args.Error(1)
}

func (m *MockCartRepository) FindByUserID(ctx context.Context, userID uint) (*domain.Cart, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Cart), args.Error(1)
}

func (m *MockCartRepository) Create(ctx context.Context, cart *domain.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

func (m *MockCartRepository) Update(ctx context.Context, cart *domain.Cart) error {
	args := m.Called(ctx, cart)
	return args.Error(0)
}

func (m *MockCartRepository) Delete(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockCartRepository) AddItem(ctx context.Context, cartItem *domain.CartItem) error {
	args := m.Called(ctx, cartItem)
	return args.Error(0)
}

func (m *MockCartRepository) UpdateItem(ctx context.Context, cartItem *domain.CartItem) error {
	args := m.Called(ctx, cartItem)
	return args.Error(0)
}

func (m *MockCartRepository) RemoveItem(ctx context.Context, cartID, itemID uint) error {
	args := m.Called(ctx, cartID, itemID)
	return args.Error(0)
}

func (m *MockCartRepository) ClearCart(ctx context.Context, cartID uint) error {
	args := m.Called(ctx, cartID)
	return args.Error(0)
}

func (m *MockCartRepository) GetCartItems(ctx context.Context, cartID uint) ([]domain.CartItem, error) {
	args := m.Called(ctx, cartID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domain.CartItem), args.Error(1)
}

func (m *MockCartRepository) GetCartTotal(ctx context.Context, cartID uint) (float64, error) {
	args := m.Called(ctx, cartID)
	return args.Get(0).(float64), args.Error(1)
}

// CartRepositoryTestSuite is a test suite for CartRepository
type CartRepositoryTestSuite struct {
	suite.Suite
	mockRepo repository.CartRepository
	ctx      context.Context
}

// SetupTest sets up the test suite
func (s *CartRepositoryTestSuite) SetupTest() {
	s.mockRepo = new(MockCartRepository)
	s.ctx = context.Background()
}

// TestFindByID tests the FindByID method
func (s *CartRepositoryTestSuite) TestFindByID() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully find a cart by ID
		cartID := uint(1)
		expectedCart := &domain.Cart{
			ID:     cartID,
			UserID: uint(1),
			Items: []domain.CartItem{
				{
					ID:        uint(1),
					CartID:    cartID,
					ProductID: uint(2),
					Quantity:  3,
					Product: domain.Product{
						ID:    uint(2),
						Name:  "Test Product",
						Price: 10.99,
					},
				},
			},
		}

		mockRepo.On("FindByID", s.ctx, cartID).Return(expectedCart, nil).Once()

		// Execute
		cart, err := s.mockRepo.FindByID(s.ctx, cartID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedCart, cart)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Cart not found
		cartID := uint(999)
		expectedError := errors.New("cart not found")

		mockRepo.On("FindByID", s.ctx, cartID).Return(nil, expectedError).Once()

		// Execute
		cart, err := s.mockRepo.FindByID(s.ctx, cartID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), cart)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestFindByUserID tests the FindByUserID method
func (s *CartRepositoryTestSuite) TestFindByUserID() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully find a cart by user ID
		userID := uint(1)
		cartID := uint(1)
		expectedCart := &domain.Cart{
			ID:     cartID,
			UserID: userID,
			Items: []domain.CartItem{
				{
					ID:        uint(1),
					CartID:    cartID,
					ProductID: uint(2),
					Quantity:  3,
					Product: domain.Product{
						ID:    uint(2),
						Name:  "Test Product",
						Price: 10.99,
					},
				},
			},
		}

		mockRepo.On("FindByUserID", s.ctx, userID).Return(expectedCart, nil).Once()

		// Execute
		cart, err := s.mockRepo.FindByUserID(s.ctx, userID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedCart, cart)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - User Has No Cart", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: User has no cart
		userID := uint(999)
		expectedError := errors.New("user has no cart")

		mockRepo.On("FindByUserID", s.ctx, userID).Return(nil, expectedError).Once()

		// Execute
		cart, err := s.mockRepo.FindByUserID(s.ctx, userID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), cart)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestCreate tests the Create method
func (s *CartRepositoryTestSuite) TestCreate() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully create a cart
		cart := &domain.Cart{
			UserID: uint(1),
		}

		mockRepo.On("Create", s.ctx, cart).Return(nil).Once()

		// Execute
		err := s.mockRepo.Create(s.ctx, cart)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Database error
		cart := &domain.Cart{
			UserID: uint(1),
		}
		expectedError := errors.New("database error")

		mockRepo.On("Create", s.ctx, cart).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Create(s.ctx, cart)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdate tests the Update method
func (s *CartRepositoryTestSuite) TestUpdate() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully update a cart
		cart := &domain.Cart{
			ID:     uint(1),
			UserID: uint(1),
		}

		mockRepo.On("Update", s.ctx, cart).Return(nil).Once()

		// Execute
		err := s.mockRepo.Update(s.ctx, cart)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Database error
		cart := &domain.Cart{
			ID:     uint(1),
			UserID: uint(1),
		}
		expectedError := errors.New("database error")

		mockRepo.On("Update", s.ctx, cart).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Update(s.ctx, cart)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestDelete tests the Delete method
func (s *CartRepositoryTestSuite) TestDelete() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully delete a cart
		cartID := uint(1)

		mockRepo.On("Delete", s.ctx, cartID).Return(nil).Once()

		// Execute
		err := s.mockRepo.Delete(s.ctx, cartID)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Cart Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Cart not found
		cartID := uint(999)
		expectedError := errors.New("cart not found")

		mockRepo.On("Delete", s.ctx, cartID).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.Delete(s.ctx, cartID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestAddItem tests the AddItem method
func (s *CartRepositoryTestSuite) TestAddItem() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully add an item to a cart
		cartItem := &domain.CartItem{
			CartID:    uint(1),
			ProductID: uint(2),
			Quantity:  3,
		}

		mockRepo.On("AddItem", s.ctx, cartItem).Return(nil).Once()

		// Execute
		err := s.mockRepo.AddItem(s.ctx, cartItem)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Database Error", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Database error
		cartItem := &domain.CartItem{
			CartID:    uint(1),
			ProductID: uint(2),
			Quantity:  3,
		}
		expectedError := errors.New("database error")

		mockRepo.On("AddItem", s.ctx, cartItem).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.AddItem(s.ctx, cartItem)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestUpdateItem tests the UpdateItem method
func (s *CartRepositoryTestSuite) TestUpdateItem() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully update an item in a cart
		cartItem := &domain.CartItem{
			ID:        uint(1),
			CartID:    uint(1),
			ProductID: uint(2),
			Quantity:  5,
		}

		mockRepo.On("UpdateItem", s.ctx, cartItem).Return(nil).Once()

		// Execute
		err := s.mockRepo.UpdateItem(s.ctx, cartItem)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Item Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Item not found
		cartItem := &domain.CartItem{
			ID:        uint(999),
			CartID:    uint(1),
			ProductID: uint(2),
			Quantity:  5,
		}
		expectedError := errors.New("item not found")

		mockRepo.On("UpdateItem", s.ctx, cartItem).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.UpdateItem(s.ctx, cartItem)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestRemoveItem tests the RemoveItem method
func (s *CartRepositoryTestSuite) TestRemoveItem() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully remove an item from a cart
		cartID := uint(1)
		itemID := uint(2)

		mockRepo.On("RemoveItem", s.ctx, cartID, itemID).Return(nil).Once()

		// Execute
		err := s.mockRepo.RemoveItem(s.ctx, cartID, itemID)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Item Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Item not found
		cartID := uint(1)
		itemID := uint(999)
		expectedError := errors.New("item not found")

		mockRepo.On("RemoveItem", s.ctx, cartID, itemID).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.RemoveItem(s.ctx, cartID, itemID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestClearCart tests the ClearCart method
func (s *CartRepositoryTestSuite) TestClearCart() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully clear a cart
		cartID := uint(1)

		mockRepo.On("ClearCart", s.ctx, cartID).Return(nil).Once()

		// Execute
		err := s.mockRepo.ClearCart(s.ctx, cartID)

		// Assert
		assert.NoError(s.T(), err)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Cart Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Cart not found
		cartID := uint(999)
		expectedError := errors.New("cart not found")

		mockRepo.On("ClearCart", s.ctx, cartID).Return(expectedError).Once()

		// Execute
		err := s.mockRepo.ClearCart(s.ctx, cartID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestGetCartItems tests the GetCartItems method
func (s *CartRepositoryTestSuite) TestGetCartItems() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully get items from a cart
		cartID := uint(1)
		expectedItems := []domain.CartItem{
			{
				ID:        uint(1),
				CartID:    cartID,
				ProductID: uint(2),
				Quantity:  3,
				Product: domain.Product{
					ID:    uint(2),
					Name:  "Product 1",
					Price: 10.99,
				},
			},
			{
				ID:        uint(2),
				CartID:    cartID,
				ProductID: uint(3),
				Quantity:  1,
				Product: domain.Product{
					ID:    uint(3),
					Name:  "Product 2",
					Price: 5.99,
				},
			},
		}

		mockRepo.On("GetCartItems", s.ctx, cartID).Return(expectedItems, nil).Once()

		// Execute
		items, err := s.mockRepo.GetCartItems(s.ctx, cartID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedItems, items)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Empty Cart", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Cart is empty
		cartID := uint(1)
		var expectedItems []domain.CartItem

		mockRepo.On("GetCartItems", s.ctx, cartID).Return(expectedItems, nil).Once()

		// Execute
		items, err := s.mockRepo.GetCartItems(s.ctx, cartID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Empty(s.T(), items)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Cart Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Cart not found
		cartID := uint(999)
		expectedError := errors.New("cart not found")

		mockRepo.On("GetCartItems", s.ctx, cartID).Return(nil, expectedError).Once()

		// Execute
		items, err := s.mockRepo.GetCartItems(s.ctx, cartID)

		// Assert
		assert.Error(s.T(), err)
		assert.Nil(s.T(), items)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestGetCartTotal tests the GetCartTotal method
func (s *CartRepositoryTestSuite) TestGetCartTotal() {
	mockRepo := s.mockRepo.(*MockCartRepository)

	s.Run("Success", func() {
		// Test case: Successfully calculate the total price of a cart
		cartID := uint(1)
		expectedTotal := 38.96 // (10.99 * 3) + (5.99 * 1) = 38.96

		mockRepo.On("GetCartTotal", s.ctx, cartID).Return(expectedTotal, nil).Once()

		// Execute
		total, err := s.mockRepo.GetCartTotal(s.ctx, cartID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Success - Empty Cart", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Cart is empty
		cartID := uint(1)
		expectedTotal := 0.0

		mockRepo.On("GetCartTotal", s.ctx, cartID).Return(expectedTotal, nil).Once()

		// Execute
		total, err := s.mockRepo.GetCartTotal(s.ctx, cartID)

		// Assert
		assert.NoError(s.T(), err)
		assert.Equal(s.T(), expectedTotal, total)
		mockRepo.AssertExpectations(s.T())
	})

	s.Run("Error - Cart Not Found", func() {
		// Reset mock
		s.SetupTest()
		mockRepo = s.mockRepo.(*MockCartRepository)

		// Test case: Cart not found
		cartID := uint(999)
		expectedError := errors.New("cart not found")

		mockRepo.On("GetCartTotal", s.ctx, cartID).Return(0.0, expectedError).Once()

		// Execute
		total, err := s.mockRepo.GetCartTotal(s.ctx, cartID)

		// Assert
		assert.Error(s.T(), err)
		assert.Equal(s.T(), 0.0, total)
		assert.Equal(s.T(), expectedError, err)
		mockRepo.AssertExpectations(s.T())
	})
}

// TestCartRepositorySuite runs the test suite
func TestCartRepositorySuite(t *testing.T) {
	suite.Run(t, new(CartRepositoryTestSuite))
}
