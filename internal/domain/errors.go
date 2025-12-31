package domain

import "errors"

// Poison errors - non-retryable, should go to DLQ
var (
	// ErrInsufficientStock indicates the product does not have enough stock
	ErrInsufficientStock = errors.New("insufficient stock")

	// ErrInvalidPayload indicates the message payload could not be parsed
	ErrInvalidPayload = errors.New("invalid payload")

	// ErrOrderNotFound indicates the order does not exist
	ErrOrderNotFound = errors.New("order not found")

	// ErrProductNotFound indicates the product does not exist
	ErrProductNotFound = errors.New("product not found")

	// ErrUserNotFound indicates the user does not exist
	ErrUserNotFound = errors.New("user not found")
)

// Transient errors - retryable
var (
	// ErrDatabaseConnection indicates a temporary database connection issue
	ErrDatabaseConnection = errors.New("database connection error")

	// ErrKafkaConnection indicates a temporary Kafka connection issue
	ErrKafkaConnection = errors.New("kafka connection error")
)

// Business logic errors
var (
	// ErrDuplicateEmail indicates the email already exists
	ErrDuplicateEmail = errors.New("email already exists")

	// ErrInvalidOrderStatus indicates an invalid order status transition
	ErrInvalidOrderStatus = errors.New("invalid order status")

	// ErrOrderAlreadyProcessed indicates the order event was already processed
	ErrOrderAlreadyProcessed = errors.New("order already processed")
)
