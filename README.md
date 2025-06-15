# Awesome E-commerce

A modern e-commerce platform built with Go, featuring a microservices architecture with an API service for business logic and a worker service for asynchronous tasks.

## Technologies Used

- **Backend**: Go (Golang)
- **Web Framework**: Gin
- **ORM**: GORM
- **Database**: MySQL
- **Caching**: Redis
- **Message Broker**: Kafka
- **Dependency Injection**: Uber FX
- **API Documentation**: Swagger
- **Testing**: Testify
- **Containerization**: Docker & Docker Compose

## Project Structure

```
awesomeEcommerce/
├── cmd/                    # Application entry points
│   ├── api/                # API service
│   └── worker/             # Worker service for async tasks
├── internal/               # Private application code
│   ├── api/                # API handlers
│   ├── config/             # Configuration
│   ├── domain/             # Domain models
│   ├── messaging/          # Kafka messaging
│   ├── middleware/         # HTTP middleware
│   ├── repository/         # Data access layer
│   ├── service/            # Business logic
│   └── worker/             # Worker implementations
├── pkg/                    # Public libraries
├── test/                   # Test utilities and integration tests
├── docker-compose.yml      # Docker Compose configuration
├── Dockerfile              # Docker image definition
└── README.md               # This file
```

## Prerequisites

- Docker and Docker Compose installed on your machine
- Git (to clone the repository)

## Getting Started

### Clone the Repository

```bash
git clone https://github.com/yourusername/awesomeEcommerce.git
cd awesomeEcommerce
```

### Run with Docker Compose

1. Start all services:

```bash
docker-compose up -d
```

This will start the following services:
- MySQL database
- Redis cache
- Zookeeper
- Kafka
- Kafka UI (for monitoring Kafka)
- API service
- Worker service

2. Check the status of the services:

```bash
docker-compose ps
```

3. View logs:

```bash
# View logs of all services
docker-compose logs

# View logs of a specific service
docker-compose logs api
docker-compose logs worker
```

### Accessing the Services

- **API Service**: http://localhost:8080
- **Kafka UI**: http://localhost:8090

### API Endpoints

#### Products
- `GET /api/v1/products` - List all products
- `GET /api/v1/products/:id` - Get product details
- `POST /api/v1/products` - Create a new product (admin only)
- `PUT /api/v1/products/:id` - Update a product (admin only)
- `DELETE /api/v1/products/:id` - Delete a product (admin only)

#### Cart
- `GET /api/v1/cart` - View cart
- `POST /api/v1/cart/items` - Add item to cart
- `PUT /api/v1/cart/items/:id` - Update cart item
- `DELETE /api/v1/cart/items/:id` - Remove item from cart
- `DELETE /api/v1/cart` - Clear cart

#### Orders
- `GET /api/v1/orders` - List user orders
- `GET /api/v1/orders/:id` - Get order details
- `POST /api/v1/orders` - Create a new order from cart
- `PUT /api/v1/orders/:id/cancel` - Cancel an order

#### Payments
- `POST /api/v1/payments` - Process payment for an order
- `GET /api/v1/payments/:id` - Get payment details

#### Users
- `POST /api/v1/users/register` - Register a new user
- `POST /api/v1/users/login` - Login
- `GET /api/v1/users/profile` - Get user profile
- `PUT /api/v1/users/profile` - Update user profile

### Stopping the Services

```bash
docker-compose down
```

To remove all data volumes as well:

```bash
docker-compose down -v
```

## Development

### Running Locally (Without Docker)

1. Install Go 1.20 or later (project uses Go 1.24 features, which may require using the development version)
2. Install MySQL, Redis, and Kafka locally
3. Set up environment variables (see docker-compose.yml for reference)
4. Run the API service:

```bash
go run cmd/api/main.go
```

5. Run the Worker service:

```bash
go run cmd/worker/main.go
```

### API Documentation

The API is documented using Swagger. Once the application is running, you can access the Swagger UI at:

```
http://localhost:8080/swagger/index.html
```

This provides an interactive documentation where you can:
- View all available endpoints
- See request/response schemas
- Test API endpoints directly from the browser

### Running Tests

#### Unit Tests

Run all unit tests with:

```bash
go test ./...
```

#### End-to-End Tests

The project includes comprehensive end-to-end tests that verify the functionality of the entire system. These tests make HTTP requests to the API and check the side effects on the database, Redis cache, and Kafka message broker.

To run the e2e tests:

1. Start the test environment:
```bash
docker-compose -f docker-compose-test.yml up -d
```

2. Run the tests:
```bash
cd tests/e2e
go test -v
```

3. Shut down the test environment when done:
```bash
docker-compose -f docker-compose-test.yml down
```

For more details on e2e tests, see [tests/e2e/README.md](tests/e2e/README.md).

## Contributing

Contributions to the Awesome E-commerce project are welcome! Here's how you can contribute:

1. **Fork the repository**
2. **Create a feature branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```
3. **Make your changes**:
   - Follow the existing code style
   - Add tests for new functionality
   - Update documentation as needed
4. **Run tests**:
   ```bash
   go test ./...
   ```
5. **Commit your changes**:
   ```bash
   git commit -m "Add your meaningful commit message"
   ```
6. **Push to your branch**:
   ```bash
   git push origin feature/your-feature-name
   ```
7. **Create a pull request**

### Code Style Guidelines

- Follow standard Go code conventions
- Use meaningful variable and function names
- Write comments for complex logic
- Include unit tests for new functionality

## License

[MIT License](LICENSE)
