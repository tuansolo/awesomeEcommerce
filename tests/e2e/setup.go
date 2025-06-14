package e2e

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/require"
)

const (
	// API endpoint for tests
	BaseURL = "http://localhost:8081/api/v1"

	// Timeout for HTTP requests
	RequestTimeout = 10 * time.Second

	// Kafka settings
	KafkaBroker = "localhost:9093"
)

// TestClient is a helper for making HTTP requests in tests
type TestClient struct {
	BaseURL    string
	HTTPClient *http.Client
	AuthToken  string
}

// NewTestClient creates a new test client
func NewTestClient() *TestClient {
	return &TestClient{
		BaseURL: BaseURL,
		HTTPClient: &http.Client{
			Timeout: RequestTimeout,
		},
	}
}

// SetAuthToken sets the authentication token for subsequent requests
func (c *TestClient) SetAuthToken(token string) {
	c.AuthToken = token
}

// DoRequest makes an HTTP request and returns the response
func (c *TestClient) DoRequest(method, path string, body interface{}) (*http.Response, error) {
	var reqBody []byte
	var err error

	if body != nil {
		reqBody, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req, err := http.NewRequest(method, c.BaseURL+path, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if c.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthToken)
	}

	return c.HTTPClient.Do(req)
}

// ParseResponse parses the response body into the provided struct
func ParseResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("failed to parse error response: %w", err)
		}
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, errResp.Error)
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}
	}

	return nil
}

// KafkaConsumerTest is a helper for testing Kafka message consumption
type KafkaConsumerTest struct {
	Reader *kafka.Reader
	Topic  string
}

// NewKafkaConsumerTest creates a new Kafka consumer for testing
func NewKafkaConsumerTest(topic string) *KafkaConsumerTest {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{KafkaBroker},
		Topic:       topic,
		GroupID:     "e2e-test-consumer-" + topic,
		MinBytes:    10e3, // 10KB
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
	})

	return &KafkaConsumerTest{
		Reader: reader,
		Topic:  topic,
	}
}

// ConsumeMessage consumes a message from Kafka with a timeout
func (k *KafkaConsumerTest) ConsumeMessage(ctx context.Context, timeout time.Duration) (kafka.Message, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	return k.Reader.ReadMessage(ctx)
}

// Close closes the Kafka reader
func (k *KafkaConsumerTest) Close() error {
	return k.Reader.Close()
}

// SetupTest sets up the test environment
func SetupTest(t *testing.T) *TestClient {
	// Check if the API is available
	client := NewTestClient()

	// Try to connect to the API with retries
	var resp *http.Response
	var err error

	for i := 0; i < 5; i++ {
		resp, err = client.DoRequest(http.MethodGet, "/health", nil)
		if err == nil && resp.StatusCode == http.StatusOK {
			break
		}

		t.Logf("API not ready, retrying in 2 seconds... (attempt %d/5)", i+1)
		time.Sleep(2 * time.Second)
	}

	require.NoError(t, err, "Failed to connect to API")
	require.Equal(t, http.StatusOK, resp.StatusCode, "API health check failed")

	return client
}

// TeardownTest tears down the test environment
func TeardownTest(t *testing.T) {
	// Nothing to do for now, but we can add cleanup logic here if needed
}

// RegisterTestUser registers a test user and returns the auth token
func RegisterTestUser(t *testing.T, client *TestClient) string {
	// Generate a unique email to avoid conflicts
	email := fmt.Sprintf("test-user-%d@example.com", time.Now().UnixNano())

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
		Token string `json:"token"`
	}

	err = ParseResponse(resp, &respBody)
	require.NoError(t, err, "Failed to parse registration response")
	require.NotEmpty(t, respBody.Token, "Auth token is empty")

	return respBody.Token
}

// CreateTestAdmin creates a test admin user directly in the database
// This is a placeholder - in a real implementation, you would need to
// connect to the database and create the admin user
func CreateTestAdmin(t *testing.T) string {
	// This is a placeholder - in a real implementation, you would need to
	// connect to the database and create the admin user
	// For now, we'll just return a dummy token
	return "admin-token"
}
