# End-to-End (E2E) Tests for Awesome E-Commerce

This directory contains end-to-end tests for the Awesome E-Commerce application. These tests verify the functionality of the entire system by making HTTP requests to the API and checking the side effects on the database, Redis cache, and Kafka message broker.

## Test Structure

The tests are organized into the following files:

- `setup.go`: Contains common setup and teardown functions, as well as helper functions for making HTTP requests and verifying responses.
- `user_test.go`: Tests for user registration, login, and profile management.
- `product_test.go`: Tests for product and category management.
- `cart_order_test.go`: Tests for cart operations and order creation/management.
- `payment_test.go`: Tests for payment processing.
- `kafka_test.go`: Tests for Kafka message publishing and consumption.

## Prerequisites

Before running the tests, you need to have the following installed:

- Docker and Docker Compose
- Go 1.20 or later

## Running the Tests

### 1. Start the Test Environment

First, start the test environment using Docker Compose:

```bash
docker-compose -f docker-compose-test.yml up -d
```

This will start the following services:
- MySQL database (port 3307)
- Redis cache (port 6380)
- Zookeeper (port 2182)
- Kafka (port 9093)
- API service (port 8081)

### 2. Wait for Services to be Ready

Wait for all services to be healthy before running the tests. You can check the status with:

```bash
docker-compose -f docker-compose-test.yml ps
```

All services should show as "Up (healthy)".

### 3. Run the Tests

Run the tests using the Go test command:

```bash
cd tests/e2e
go test -v
```

To run a specific test file:

```bash
go test -v user_test.go setup.go
go test -v product_test.go setup.go
go test -v cart_order_test.go setup.go
go test -v payment_test.go setup.go
go test -v kafka_test.go setup.go
```

### 4. Shut Down the Test Environment

After running the tests, shut down the test environment:

```bash
docker-compose -f docker-compose-test.yml down
```

To remove all data volumes as well:

```bash
docker-compose -f docker-compose-test.yml down -v
```

## Test Coverage

The e2e tests cover the following functionality:

1. **User Management**
   - User registration
   - User login
   - Getting user profile
   - Updating user profile
   - Updating user password

2. **Product Management**
   - Creating categories
   - Creating products
   - Getting products
   - Getting products by category
   - Updating products
   - Deleting products

3. **Cart Operations**
   - Adding items to cart
   - Updating cart items
   - Getting cart contents
   - Getting cart total
   - Removing items from cart
   - Clearing cart

4. **Order Management**
   - Creating orders
   - Getting user orders
   - Getting specific order
   - Cancelling orders

5. **Payment Processing**
   - Creating payments
   - Getting payment by order ID
   - Refunding payments

6. **Kafka Messaging**
   - Verifying order creation messages
   - Verifying order status update messages
   - Verifying payment status messages

## Troubleshooting

If you encounter issues running the tests, try the following:

1. **API not reachable**: Make sure the API service is running and healthy. Check the logs with `docker-compose -f docker-compose-test.yml logs api`.

2. **Database connection issues**: Check the database logs with `docker-compose -f docker-compose-test.yml logs mysql`.

3. **Kafka connection issues**: Check the Kafka logs with `docker-compose -f docker-compose-test.yml logs kafka`.

4. **Tests timing out**: The tests include retries and timeouts, but in some environments, you might need to increase the timeout values in `setup.go`.

5. **Clean state**: If tests are failing due to data conflicts, try shutting down the environment with `docker-compose -f docker-compose-test.yml down -v` to remove all data volumes, then start it again.