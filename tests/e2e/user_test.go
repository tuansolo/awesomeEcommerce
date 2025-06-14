package e2e

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRegistrationAndLogin(t *testing.T) {
	// Setup test environment
	client := SetupTest(t)
	defer TeardownTest(t)

	// Test user registration
	t.Run("Register User", func(t *testing.T) {
		// Generate a unique email to avoid conflicts
		email := "test-user-registration@example.com"

		reqBody := map[string]interface{}{
			"email":      email,
			"password":   "password123",
			"first_name": "Test",
			"last_name":  "User",
			"phone":      "1234567890",
			"address":    "123 Test St",
		}

		resp, err := client.DoRequest(http.MethodPost, "/users/register", reqBody)
		require.NoError(t, err, "Failed to register test user")

		var respBody struct {
			Message string `json:"message"`
			Token   string `json:"token"`
			User    struct {
				ID        uint   `json:"id"`
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Role      string `json:"role"`
			} `json:"user"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse registration response")

		// Verify response
		assert.Equal(t, http.StatusCreated, resp.StatusCode, "Expected status code 201")
		assert.Equal(t, "User registered successfully", respBody.Message)
		assert.NotEmpty(t, respBody.Token, "Auth token should not be empty")
		assert.Equal(t, email, respBody.User.Email)
		assert.Equal(t, "Test", respBody.User.FirstName)
		assert.Equal(t, "User", respBody.User.LastName)
		assert.Equal(t, "customer", respBody.User.Role)

		// Save token for next test
		client.SetAuthToken(respBody.Token)
	})

	// Test user login
	t.Run("Login User", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"email":    "test-user-registration@example.com",
			"password": "password123",
		}

		resp, err := client.DoRequest(http.MethodPost, "/users/login", reqBody)
		require.NoError(t, err, "Failed to login test user")

		var respBody struct {
			Token string `json:"token"`
			User  struct {
				ID        uint   `json:"id"`
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Role      string `json:"role"`
			} `json:"user"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse login response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.NotEmpty(t, respBody.Token, "Auth token should not be empty")
		assert.Equal(t, "test-user-registration@example.com", respBody.User.Email)
		assert.Equal(t, "Test", respBody.User.FirstName)
		assert.Equal(t, "User", respBody.User.LastName)
		assert.Equal(t, "customer", respBody.User.Role)

		// Save token for next test
		client.SetAuthToken(respBody.Token)
	})

	// Test get user profile
	t.Run("Get User Profile", func(t *testing.T) {
		resp, err := client.DoRequest(http.MethodGet, "/users/me", nil)
		require.NoError(t, err, "Failed to get user profile")

		var respBody struct {
			User struct {
				ID        uint   `json:"id"`
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Phone     string `json:"phone"`
				Address   string `json:"address"`
				Role      string `json:"role"`
			} `json:"user"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse profile response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, "test-user-registration@example.com", respBody.User.Email)
		assert.Equal(t, "Test", respBody.User.FirstName)
		assert.Equal(t, "User", respBody.User.LastName)
		assert.Equal(t, "1234567890", respBody.User.Phone)
		assert.Equal(t, "123 Test St", respBody.User.Address)
		assert.Equal(t, "customer", respBody.User.Role)
	})

	// Test update user profile
	t.Run("Update User Profile", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"first_name": "Updated",
			"last_name":  "User",
			"phone":      "9876543210",
			"address":    "456 Updated St",
		}

		resp, err := client.DoRequest(http.MethodPut, "/users/me", reqBody)
		require.NoError(t, err, "Failed to update user profile")

		var respBody struct {
			Message string `json:"message"`
			User    struct {
				ID        uint   `json:"id"`
				Email     string `json:"email"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Phone     string `json:"phone"`
				Address   string `json:"address"`
				Role      string `json:"role"`
			} `json:"user"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse update response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, "Profile updated successfully", respBody.Message)
		assert.Equal(t, "test-user-registration@example.com", respBody.User.Email)
		assert.Equal(t, "Updated", respBody.User.FirstName)
		assert.Equal(t, "User", respBody.User.LastName)
		assert.Equal(t, "9876543210", respBody.User.Phone)
		assert.Equal(t, "456 Updated St", respBody.User.Address)
		assert.Equal(t, "customer", respBody.User.Role)
	})

	// Test update user password
	t.Run("Update User Password", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"current_password": "password123",
			"new_password":     "newpassword123",
		}

		resp, err := client.DoRequest(http.MethodPut, "/users/me/password", reqBody)
		require.NoError(t, err, "Failed to update user password")

		var respBody struct {
			Message string `json:"message"`
		}

		err = ParseResponse(resp, &respBody)
		require.NoError(t, err, "Failed to parse password update response")

		// Verify response
		assert.Equal(t, http.StatusOK, resp.StatusCode, "Expected status code 200")
		assert.Equal(t, "Password updated successfully", respBody.Message)

		// Test login with new password
		loginReqBody := map[string]interface{}{
			"email":    "test-user-registration@example.com",
			"password": "newpassword123",
		}

		loginResp, err := client.DoRequest(http.MethodPost, "/users/login", loginReqBody)
		require.NoError(t, err, "Failed to login with new password")
		assert.Equal(t, http.StatusOK, loginResp.StatusCode, "Expected status code 200")
	})
}
