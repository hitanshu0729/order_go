# API Documentation

Base URL: `/api/v1`

---

## Table of Contents

- [Health & Status](#health--status)
- [Users](#users)
- [Products](#products)
- [Orders](#orders)
- [Order Items](#order-items)

---

## Health & Status

### Hello World

```
GET /api/v1/
```

**Response:**
```json
{
  "message": "Hello World"
}
```

---

### Health Check

```
GET /api/v1/health
```

**Response:**
```json
{
  "status": "up",
  "message": "It's healthy"
}
```

---

## Users

### Get All Users

```
GET /api/v1/users
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  }
]
```

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 500 | Internal Server Error |

---

### Create User

```
POST /api/v1/users
```

**Request Body:**
```json
{
  "name": "John Doe",
  "email": "john@example.com"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | User's name |
| email | string | Yes | User's email (must be valid email format) |

**Response:**
```json
{
  "message": "User created"
}
```

| Status Code | Description |
|-------------|-------------|
| 201 | User created successfully |
| 400 | Bad Request (validation error) |
| 500 | Internal Server Error |

---

### Get User by ID

```
GET /api/v1/users/:id
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | User ID |

**Response:**
```json
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com"
}
```

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Invalid user ID |
| 404 | User not found |

---

### Update User

```
PATCH /api/v1/users/:id
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | User ID |

**Request Body:**
```json
{
  "name": "John Updated",
  "email": "john.updated@example.com"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | Yes | User's name |
| email | string | Yes | User's email (must be valid email format) |

**Response:**
```json
{
  "message": "user updated"
}
```

| Status Code | Description |
|-------------|-------------|
| 200 | User updated successfully |
| 400 | Invalid user ID or validation error |
| 500 | Internal Server Error |

---

### Delete User

```
DELETE /api/v1/users/:id
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | User ID |

**Response:**
```json
{
  "message": "user deleted"
}
```

| Status Code | Description |
|-------------|-------------|
| 200 | User deleted successfully |
| 400 | Invalid user ID |
| 500 | Internal Server Error |

---

## Products

### Get All Products

```
GET /api/v1/products
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Product Name",
    "price": 1000,
    "stock": 50
  }
]
```

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 500 | Internal Server Error |

---

### Create Product

```
POST /api/v1/products
```

**Request Body:**
```json
{
  "name": "Product Name",
  "price": 1000,
  "stock": 50
}
```

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| name | string | Yes | - | Product name |
| price | integer | Yes | > 0 | Product price (in smallest currency unit) |
| stock | integer | Yes | >= 0 | Available stock quantity |

**Response:**
```json
{
  "message": "Product created"
}
```

| Status Code | Description |
|-------------|-------------|
| 201 | Product created successfully |
| 400 | Bad Request (validation error) |
| 500 | Internal Server Error |

---

### Get Product by ID

```
GET /api/v1/products/:id
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Product ID |

**Response:**
```json
{
  "id": 1,
  "name": "Product Name",
  "price": 1000,
  "stock": 50
}
```

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Invalid product ID |
| 404 | Product not found |

---

## Orders

### Create Order

```
POST /api/v1/orders
```

**Request Body:**
```json
{
  "user_id": 1
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| user_id | integer | Yes | ID of the user creating the order |

**Response:**
```json
{
  "message": "Order created"
}
```

**Side Effects:**
- Publishes `order.created` event to Kafka

| Status Code | Description |
|-------------|-------------|
| 201 | Order created successfully |
| 400 | Bad Request (validation error) |
| 500 | Internal Server Error |

---

### Get Orders

```
GET /api/v1/orders
```

**Query Parameters:**
| Parameter | Type | Required | Format | Description |
|-----------|------|----------|--------|-------------|
| user_id | integer | No | - | Filter by user ID |
| status | string | No | pending/paid/cancelled/completed | Filter by order status |
| from | string | No | YYYY-MM-DD | Filter orders from this date |
| to | string | No | YYYY-MM-DD | Filter orders up to this date |

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "status": "pending",
    "total_amount": 2000,
    "created_at": "2025-12-31T10:00:00Z"
  }
]
```

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Invalid query parameters |
| 500 | Internal Server Error |

---

### Get Order by ID

```
GET /api/v1/orders/:id
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Order ID |

**Response:**
```json
{
  "id": 1,
  "user_id": 1,
  "status": "pending",
  "total_amount": 2000,
  "created_at": "2025-12-31T10:00:00Z"
}
```

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Invalid order ID |
| 404 | Order not found |

---

### Get Orders by Status

```
GET /api/v1/orders/status/:status
```

**Path Parameters:**
| Parameter | Type | Values | Description |
|-----------|------|--------|-------------|
| status | string | pending, paid, cancelled, completed | Order status |

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "status": "pending",
    "total_amount": 2000,
    "created_at": "2025-12-31T10:00:00Z"
  }
]
```

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 500 | Internal Server Error |

---

### Update Order Status

```
PATCH /api/v1/orders/:id/status
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Order ID |

**Request Body:**
```json
{
  "status": "paid"
}
```

| Field | Type | Required | Values | Description |
|-------|------|----------|--------|-------------|
| status | string | Yes | pending, paid, cancelled, completed | New order status |

**Response:**
```json
{
  "message": "order status updated"
}
```

| Status Code | Description |
|-------------|-------------|
| 200 | Status updated successfully |
| 400 | Invalid order ID or validation error |
| 500 | Internal Server Error |

---

### Cancel Order

```
POST /api/v1/orders/:id/cancel
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Order ID |

**Business Rules:**
- Cannot cancel an order that is already `completed` (shipped)

**Response:**
```json
{
  "message": "order status updated",
  "status": "cancelled"
}
```

| Status Code | Description |
|-------------|-------------|
| 200 | Order cancelled successfully |
| 400 | Invalid order ID or order already completed |
| 404 | Order not found |
| 500 | Internal Server Error |

---

### Pay Order

```
POST /api/v1/orders/:id/pay
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Order ID |

**Business Rules:**
- Order status must be `pending` to be paid

**Response:**
```json
{
  "message": "order status updated",
  "status": "paid"
}
```

**Side Effects:**
- Publishes `order.paid` event to Kafka

| Status Code | Description |
|-------------|-------------|
| 200 | Order paid successfully |
| 400 | Invalid order ID or order not in pending status |
| 404 | Order not found |
| 500 | Internal Server Error |

---

### Ship Order

```
POST /api/v1/orders/:id/ship
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Order ID |

**Business Rules:**
- Order status must be `paid` to be shipped
- Status changes to `completed` after shipping

**Response:**
```json
{
  "message": "order status updated",
  "status": "completed"
}
```

| Status Code | Description |
|-------------|-------------|
| 200 | Order shipped successfully |
| 400 | Invalid order ID or order not in paid status |
| 404 | Order not found |
| 500 | Internal Server Error |

---

## Order Items

### Get Order Items

```
GET /api/v1/orders/:id/items
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Order ID |

**Response:**
```json
[
  {
    "id": 1,
    "order_id": 1,
    "product_id": 1,
    "quantity": 2,
    "price": 1000
  }
]
```

| Status Code | Description |
|-------------|-------------|
| 200 | Success |
| 400 | Invalid order ID |
| 500 | Internal Server Error |

---

### Add Order Item

```
POST /api/v1/orders/:id/items
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Order ID |

**Request Body:**
```json
{
  "product_id": 1,
  "quantity": 2
}
```

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| product_id | integer | Yes | - | ID of the product to add |
| quantity | integer | Yes | > 0 | Quantity of the product |

**Response:**
```json
{
  "message": "item added"
}
```

| Status Code | Description |
|-------------|-------------|
| 201 | Item added successfully |
| 400 | Invalid order ID or validation error |
| 404 | Product not found |
| 500 | Internal Server Error |

---

### Update Order Item Quantity

```
PATCH /api/v1/orders/:id/items/:item_id
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Order ID |
| item_id | integer | Order Item ID |

**Business Rules:**
- Order status must be `pending` to update items

**Request Body:**
```json
{
  "quantity": 5
}
```

| Field | Type | Required | Validation | Description |
|-------|------|----------|------------|-------------|
| quantity | integer | Yes | > 0 | New quantity |

**Response:**
```json
{
  "message": "item quantity updated"
}
```

| Status Code | Description |
|-------------|-------------|
| 200 | Quantity updated successfully |
| 400 | Invalid ID or order not in pending status |
| 404 | Order not found |
| 500 | Internal Server Error |

---

### Remove Order Item

```
DELETE /api/v1/orders/:id/items/:item_id
```

**Path Parameters:**
| Parameter | Type | Description |
|-----------|------|-------------|
| id | integer | Order ID |
| item_id | integer | Order Item ID |

**Business Rules:**
- Order status must be `pending` to remove items

**Response:**
```json
{
  "message": "item removed"
}
```

| Status Code | Description |
|-------------|-------------|
| 200 | Item removed successfully |
| 400 | Invalid ID or order not in pending status |
| 404 | Order not found |
| 500 | Internal Server Error |

---

## Data Models

### User

| Field | Type | Description |
|-------|------|-------------|
| id | integer | Unique identifier |
| name | string | User's name |
| email | string | User's email (unique) |

### Product

| Field | Type | Description |
|-------|------|-------------|
| id | integer | Unique identifier |
| name | string | Product name |
| price | integer | Product price (in smallest currency unit) |
| stock | integer | Available stock quantity |

### Order

| Field | Type | Description |
|-------|------|-------------|
| id | integer | Unique identifier |
| user_id | integer | Reference to user |
| status | string | Order status (pending/paid/cancelled/completed) |
| total_amount | integer | Total order amount |
| created_at | datetime | Order creation timestamp |

### Order Item

| Field | Type | Description |
|-------|------|-------------|
| id | integer | Unique identifier |
| order_id | integer | Reference to order |
| product_id | integer | Reference to product |
| quantity | integer | Quantity of the product |
| price | integer | Price at time of order |

---

## Order Status Flow

```
pending → paid → completed
    ↓       
cancelled
```

- **pending**: Initial state when order is created
- **paid**: After successful payment
- **completed**: After order is shipped
- **cancelled**: Order was cancelled (only from pending or paid states)

---

## Kafka Events

The following events are published to Kafka:

| Event | Topic | Payload | Trigger |
|-------|-------|---------|---------|
| Order Created | `order.created` | `{"user_id": <int>}` | When a new order is created |
| Order Paid | `order.paid` | `{"order_id": <int>}` | When an order is paid |

---

## Error Response Format

All error responses follow this format:

```json
{
  "error": "error message description"
}
```
