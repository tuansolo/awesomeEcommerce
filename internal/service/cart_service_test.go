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

func TestGetCartByID(t *testing.T) {
	// Setup
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)
	ctx := context.Background()
	cartID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedCart := &domain.Cart{
			ID:     cartID,
			UserID: uint(1),
			Items:  []domain.CartItem{},
		}

		// Expectations
		mockCartRepo.On("FindByID", ctx, cartID).Return(expectedCart, nil).Once()

		// Execute
		cart, err := cartService.GetCartByID(ctx, cartID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCart, cart)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("Not Found", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Expectations
		mockCartRepo.On("FindByID", ctx, cartID).Return(nil, errors.New("cart not found")).Once()

		// Execute
		cart, err := cartService.GetCartByID(ctx, cartID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cart)
		assert.Contains(t, err.Error(), "cart not found")
		mockCartRepo.AssertExpectations(t)
	})
}

func TestGetCartByUserID(t *testing.T) {
	// Setup
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)
	ctx := context.Background()
	userID := uint(1)

	t.Run("Success - Existing Cart", func(t *testing.T) {
		// Test data
		expectedCart := &domain.Cart{
			ID:     uint(1),
			UserID: userID,
			Items:  []domain.CartItem{},
		}
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockCartRepo.On("FindByUserID", ctx, userID).Return(expectedCart, nil).Once()

		// Execute
		cart, err := cartService.GetCartByUserID(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedCart, cart)
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("Success - Create New Cart", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockCartRepo.On("FindByUserID", ctx, userID).Return(nil, errors.New("cart not found")).Once()
		mockCartRepo.On("Create", ctx, mock.AnythingOfType("*domain.Cart")).Return(nil).Once().
			Run(func(args mock.Arguments) {
				cart := args.Get(1).(*domain.Cart)
				cart.ID = uint(1) // Simulate ID assignment by database
			})

		// Execute
		cart, err := cartService.GetCartByUserID(ctx, userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, uint(1), cart.ID)
		assert.Equal(t, userID, cart.UserID)
		assert.Empty(t, cart.Items)
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		// Execute
		cart, err := cartService.GetCartByUserID(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cart)
		assert.Contains(t, err.Error(), "user not found")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertNotCalled(t, "FindByUserID")
		mockCartRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Cart Creation Error", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockCartRepo.On("FindByUserID", ctx, userID).Return(nil, errors.New("cart not found")).Once()
		mockCartRepo.On("Create", ctx, mock.AnythingOfType("*domain.Cart")).Return(errors.New("database error")).Once()

		// Execute
		cart, err := cartService.GetCartByUserID(ctx, userID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, cart)
		assert.Contains(t, err.Error(), "database error")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
	})
}

func TestCreateCart(t *testing.T) {
	// Setup
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)
	ctx := context.Background()
	userID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		cart := &domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockCartRepo.On("FindByUserID", ctx, userID).Return(nil, errors.New("cart not found")).Once()
		mockCartRepo.On("Create", ctx, cart).Return(nil).Once()

		// Execute
		err := cartService.CreateCart(ctx, cart)

		// Assert
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		cart := &domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		// Execute
		err := cartService.CreateCart(ctx, cart)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertNotCalled(t, "FindByUserID")
		mockCartRepo.AssertNotCalled(t, "Create")
	})

	t.Run("Cart Already Exists", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		cart := &domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}
		existingCart := &domain.Cart{
			ID:     uint(1),
			UserID: userID,
			Items:  []domain.CartItem{},
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockCartRepo.On("FindByUserID", ctx, userID).Return(existingCart, nil).Once()

		// Execute
		err := cartService.CreateCart(ctx, cart)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user already has a cart")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("Database Error", func(t *testing.T) {
		// Reset mocks
		mockCartRepo := new(MockCartRepository)
		mockProductRepo := new(MockProductRepository)
		mockUserRepo := new(MockUserRepository)
		cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		cart := &domain.Cart{
			UserID: userID,
			Items:  []domain.CartItem{},
		}
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockCartRepo.On("FindByUserID", ctx, userID).Return(nil, errors.New("cart not found")).Once()
		mockCartRepo.On("Create", ctx, cart).Return(errors.New("database error")).Once()

		// Execute
		err := cartService.CreateCart(ctx, cart)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
	})
}

func TestAddItemToCart(t *testing.T) {
	// Setup
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)
	ctx := context.Background()
	userID := uint(1)
	productID := uint(2)
	quantity := 3

	t.Run("Success", func(t *testing.T) {
		// Test data
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}
		cart := &domain.Cart{
			ID:     uint(1),
			UserID: userID,
			Items:  []domain.CartItem{},
		}
		product := &domain.Product{
			ID:    productID,
			Name:  "Test Product",
			Price: 10.99,
			Stock: 10,
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockCartRepo.On("FindByUserID", ctx, userID).Return(cart, nil).Once()
		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Once()
		mockCartRepo.On("AddItem", ctx, mock.AnythingOfType("*domain.CartItem")).Return(nil).Once()

		// Execute
		err := cartService.AddItemToCart(ctx, userID, productID, quantity)

		// Assert
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("Invalid Quantity", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Execute
		err := cartService.AddItemToCart(ctx, userID, productID, 0)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "quantity must be greater than zero")
		mockUserRepo.AssertNotCalled(t, "FindByID")
		mockCartRepo.AssertNotCalled(t, "FindByUserID")
		mockProductRepo.AssertNotCalled(t, "FindByID")
		mockCartRepo.AssertNotCalled(t, "AddItem")
	})

	t.Run("User Not Found", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, errors.New("user not found")).Once()

		// Execute
		err := cartService.AddItemToCart(ctx, userID, productID, quantity)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertNotCalled(t, "FindByUserID")
		mockProductRepo.AssertNotCalled(t, "FindByID")
		mockCartRepo.AssertNotCalled(t, "AddItem")
	})

	t.Run("Product Not Found", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}
		cart := &domain.Cart{
			ID:     uint(1),
			UserID: userID,
			Items:  []domain.CartItem{},
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockCartRepo.On("FindByUserID", ctx, userID).Return(cart, nil).Once()
		mockProductRepo.On("FindByID", ctx, productID).Return(nil, errors.New("product not found")).Once()

		// Execute
		err := cartService.AddItemToCart(ctx, userID, productID, quantity)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product not found")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
		mockCartRepo.AssertNotCalled(t, "AddItem")
	})

	t.Run("Insufficient Stock", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}
		cart := &domain.Cart{
			ID:     uint(1),
			UserID: userID,
			Items:  []domain.CartItem{},
		}
		product := &domain.Product{
			ID:    productID,
			Name:  "Test Product",
			Price: 10.99,
			Stock: 2, // Less than requested quantity
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockCartRepo.On("FindByUserID", ctx, userID).Return(cart, nil).Once()
		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Once()

		// Execute
		err := cartService.AddItemToCart(ctx, userID, productID, quantity)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient stock")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
		mockCartRepo.AssertNotCalled(t, "AddItem")
	})

	t.Run("Database Error", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		user := &domain.User{
			ID:    userID,
			Email: "test@example.com",
		}
		cart := &domain.Cart{
			ID:     uint(1),
			UserID: userID,
			Items:  []domain.CartItem{},
		}
		product := &domain.Product{
			ID:    productID,
			Name:  "Test Product",
			Price: 10.99,
			Stock: 10,
		}

		// Expectations
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockCartRepo.On("FindByUserID", ctx, userID).Return(cart, nil).Once()
		mockProductRepo.On("FindByID", ctx, productID).Return(product, nil).Once()
		mockCartRepo.On("AddItem", ctx, mock.AnythingOfType("*domain.CartItem")).Return(errors.New("database error")).Once()

		// Execute
		err := cartService.AddItemToCart(ctx, userID, productID, quantity)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockUserRepo.AssertExpectations(t)
		mockCartRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestUpdateCartItem(t *testing.T) {
	// Setup
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)
	ctx := context.Background()
	cartID := uint(1)
	itemID := uint(2)
	quantity := 5

	t.Run("Success", func(t *testing.T) {
		// Test data
		cartItems := []domain.CartItem{
			{
				ID:        itemID,
				CartID:    cartID,
				ProductID: uint(3),
				Quantity:  2,
			},
		}
		product := &domain.Product{
			ID:    uint(3),
			Name:  "Test Product",
			Price: 10.99,
			Stock: 10,
		}

		// Expectations
		mockCartRepo.On("GetCartItems", ctx, cartID).Return(cartItems, nil).Once()
		mockProductRepo.On("FindByID", ctx, uint(3)).Return(product, nil).Once()
		mockCartRepo.On("UpdateItem", ctx, mock.AnythingOfType("*domain.CartItem")).Return(nil).Once()

		// Execute
		err := cartService.UpdateCartItem(ctx, cartID, itemID, quantity)

		// Assert
		assert.NoError(t, err)
		mockCartRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})

	t.Run("Remove Item When Quantity <= 0", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Expectations
		mockCartRepo.On("RemoveItem", ctx, cartID, itemID).Return(nil).Once()

		// Execute
		err := cartService.UpdateCartItem(ctx, cartID, itemID, 0)

		// Assert
		assert.NoError(t, err)
		mockCartRepo.AssertExpectations(t)
		mockCartRepo.AssertNotCalled(t, "GetCartItems")
		mockProductRepo.AssertNotCalled(t, "FindByID")
		mockCartRepo.AssertNotCalled(t, "UpdateItem")
	})

	t.Run("Item Not Found", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		cartItems := []domain.CartItem{
			{
				ID:        uint(999), // Different ID
				CartID:    cartID,
				ProductID: uint(3),
				Quantity:  2,
			},
		}

		// Expectations
		mockCartRepo.On("GetCartItems", ctx, cartID).Return(cartItems, nil).Once()

		// Execute
		err := cartService.UpdateCartItem(ctx, cartID, itemID, quantity)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cart item not found")
		mockCartRepo.AssertExpectations(t)
		mockProductRepo.AssertNotCalled(t, "FindByID")
		mockCartRepo.AssertNotCalled(t, "UpdateItem")
	})

	t.Run("Product Not Found", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		cartItems := []domain.CartItem{
			{
				ID:        itemID,
				CartID:    cartID,
				ProductID: uint(3),
				Quantity:  2,
			},
		}

		// Expectations
		mockCartRepo.On("GetCartItems", ctx, cartID).Return(cartItems, nil).Once()
		mockProductRepo.On("FindByID", ctx, uint(3)).Return(nil, errors.New("product not found")).Once()

		// Execute
		err := cartService.UpdateCartItem(ctx, cartID, itemID, quantity)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product not found")
		mockCartRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
		mockCartRepo.AssertNotCalled(t, "UpdateItem")
	})

	t.Run("Insufficient Stock", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		cartItems := []domain.CartItem{
			{
				ID:        itemID,
				CartID:    cartID,
				ProductID: uint(3),
				Quantity:  2,
			},
		}
		product := &domain.Product{
			ID:    uint(3),
			Name:  "Test Product",
			Price: 10.99,
			Stock: 3, // Less than requested quantity
		}

		// Expectations
		mockCartRepo.On("GetCartItems", ctx, cartID).Return(cartItems, nil).Once()
		mockProductRepo.On("FindByID", ctx, uint(3)).Return(product, nil).Once()

		// Execute
		err := cartService.UpdateCartItem(ctx, cartID, itemID, quantity)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient stock")
		mockCartRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
		mockCartRepo.AssertNotCalled(t, "UpdateItem")
	})

	t.Run("Database Error - GetCartItems", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Expectations
		mockCartRepo.On("GetCartItems", ctx, cartID).Return(nil, errors.New("database error")).Once()

		// Execute
		err := cartService.UpdateCartItem(ctx, cartID, itemID, quantity)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockCartRepo.AssertExpectations(t)
		mockProductRepo.AssertNotCalled(t, "FindByID")
		mockCartRepo.AssertNotCalled(t, "UpdateItem")
	})

	t.Run("Database Error - UpdateItem", func(t *testing.T) {
		// Reset mocks
		mockCartRepo = new(MockCartRepository)
		mockProductRepo = new(MockProductRepository)
		mockUserRepo = new(MockUserRepository)
		cartService = service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Test data
		cartItems := []domain.CartItem{
			{
				ID:        itemID,
				CartID:    cartID,
				ProductID: uint(3),
				Quantity:  2,
			},
		}
		product := &domain.Product{
			ID:    uint(3),
			Name:  "Test Product",
			Price: 10.99,
			Stock: 10,
		}

		// Expectations
		mockCartRepo.On("GetCartItems", ctx, cartID).Return(cartItems, nil).Once()
		mockProductRepo.On("FindByID", ctx, uint(3)).Return(product, nil).Once()
		mockCartRepo.On("UpdateItem", ctx, mock.AnythingOfType("*domain.CartItem")).Return(errors.New("database error")).Once()

		// Execute
		err := cartService.UpdateCartItem(ctx, cartID, itemID, quantity)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockCartRepo.AssertExpectations(t)
		mockProductRepo.AssertExpectations(t)
	})
}

func TestRemoveItemFromCart(t *testing.T) {
	// Setup
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)
	ctx := context.Background()
	cartID := uint(1)
	itemID := uint(2)

	t.Run("Success", func(t *testing.T) {
		// Expectations
		mockCartRepo.On("RemoveItem", ctx, cartID, itemID).Return(nil).Once()

		// Execute
		err := cartService.RemoveItemFromCart(ctx, cartID, itemID)

		// Assert
		assert.NoError(t, err)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("Database Error", func(t *testing.T) {
		// Reset mocks
		mockCartRepo := new(MockCartRepository)
		mockProductRepo := new(MockProductRepository)
		mockUserRepo := new(MockUserRepository)
		cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Expectations
		mockCartRepo.On("RemoveItem", ctx, cartID, itemID).Return(errors.New("database error")).Once()

		// Execute
		err := cartService.RemoveItemFromCart(ctx, cartID, itemID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockCartRepo.AssertExpectations(t)
	})
}

func TestClearCart(t *testing.T) {
	// Setup
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)
	ctx := context.Background()
	cartID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Expectations
		mockCartRepo.On("ClearCart", ctx, cartID).Return(nil).Once()

		// Execute
		err := cartService.ClearCart(ctx, cartID)

		// Assert
		assert.NoError(t, err)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("Database Error", func(t *testing.T) {
		// Reset mocks
		mockCartRepo := new(MockCartRepository)
		mockProductRepo := new(MockProductRepository)
		mockUserRepo := new(MockUserRepository)
		cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Expectations
		mockCartRepo.On("ClearCart", ctx, cartID).Return(errors.New("database error")).Once()

		// Execute
		err := cartService.ClearCart(ctx, cartID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database error")
		mockCartRepo.AssertExpectations(t)
	})
}

func TestGetCartItems(t *testing.T) {
	// Setup
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)
	ctx := context.Background()
	cartID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedItems := []domain.CartItem{
			{
				ID:        uint(1),
				CartID:    cartID,
				ProductID: uint(2),
				Quantity:  3,
			},
			{
				ID:        uint(2),
				CartID:    cartID,
				ProductID: uint(3),
				Quantity:  1,
			},
		}

		// Expectations
		mockCartRepo.On("GetCartItems", ctx, cartID).Return(expectedItems, nil).Once()

		// Execute
		items, err := cartService.GetCartItems(ctx, cartID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedItems, items)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("Database Error", func(t *testing.T) {
		// Reset mocks
		mockCartRepo := new(MockCartRepository)
		mockProductRepo := new(MockProductRepository)
		mockUserRepo := new(MockUserRepository)
		cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Expectations
		mockCartRepo.On("GetCartItems", ctx, cartID).Return(nil, errors.New("database error")).Once()

		// Execute
		items, err := cartService.GetCartItems(ctx, cartID)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, items)
		assert.Contains(t, err.Error(), "database error")
		mockCartRepo.AssertExpectations(t)
	})
}

func TestGetCartTotal(t *testing.T) {
	// Setup
	mockCartRepo := new(MockCartRepository)
	mockProductRepo := new(MockProductRepository)
	mockUserRepo := new(MockUserRepository)
	cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)
	ctx := context.Background()
	cartID := uint(1)

	t.Run("Success", func(t *testing.T) {
		// Test data
		expectedTotal := 99.99

		// Expectations
		mockCartRepo.On("GetCartTotal", ctx, cartID).Return(expectedTotal, nil).Once()

		// Execute
		total, err := cartService.GetCartTotal(ctx, cartID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedTotal, total)
		mockCartRepo.AssertExpectations(t)
	})

	t.Run("Database Error", func(t *testing.T) {
		// Reset mocks
		mockCartRepo := new(MockCartRepository)
		mockProductRepo := new(MockProductRepository)
		mockUserRepo := new(MockUserRepository)
		cartService := service.NewCartService(mockCartRepo, mockProductRepo, mockUserRepo)

		// Expectations
		mockCartRepo.On("GetCartTotal", ctx, cartID).Return(0.0, errors.New("database error")).Once()

		// Execute
		total, err := cartService.GetCartTotal(ctx, cartID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, 0.0, total)
		assert.Contains(t, err.Error(), "database error")
		mockCartRepo.AssertExpectations(t)
	})
}
