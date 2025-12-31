package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/hitanshu0729/order_go/internal/kafka"
	"github.com/hitanshu0729/order_go/internal/storage/sqlite"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orders        *sqlite.Repo
	kafkaProducer *kafka.Producer
}

func NewOrderHandler(orders *sqlite.Repo, kafkaProducer *kafka.Producer) *OrderHandler {
	return &OrderHandler{orders: orders, kafkaProducer: kafkaProducer}
}

func (h *OrderHandler) RegisterOrderRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders")
	orders.POST("", h.CreateOrder)
	orders.GET("", h.GetOrders)

	orders.GET("/:id", h.GetOrderByID)
	orders.GET("/status/:status", h.GetOrdersByStatus)

	orders.PATCH("/:id/status", h.UpdateOrderStatus)
	orders.POST("/:id/cancel", h.CancelOrder)
	orders.POST("/:id/pay", h.PayOrder)
	orders.POST("/:id/ship", h.ShipOrder)

	// ✅ Order Items — properly nested
	orders.GET("/:id/items", h.GetOrderItems)
	orders.POST("/:id/items", h.AddOrderItem)
	orders.PATCH("/:id/items/:item_id", h.UpdateOrderItemQuantity)
	orders.DELETE("/:id/items/:item_id", h.RemoveOrderItem)
}

type CreateOrderRequest struct {
	UserID int64 `json:"user_id" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending paid cancelled completed"`
}

type AddOrderItemRequest struct {
	ProductID int64 `json:"product_id" binding:"required"`
	Quantity  int64 `json:"quantity" binding:"required,gt=0"`
}

type UpdateOrderItemQuantityRequest struct {
	Quantity int64 `json:"quantity" binding:"required,gt=0"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Creating order: %+v", req)
	err := h.orders.CreateOrder(c.Request.Context(), req.UserID, "pending", 0)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Order created successfully: %+v", req)
	// ✅ AFTER DB commit
	log.Println("🚀 Publishing order.created event to Kafka")
	err = h.kafkaProducer.Publish(
		context.Background(),
		"order.created",
		map[string]any{
			"user_id": req.UserID,
		},
	)
	if err != nil {
		log.Println("❌ Kafka publish failed:", err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Order created"})
}

func (h *OrderHandler) GetOrders(c *gin.Context) {
	var filter sqlite.OrderFilter

	if userIDStr := c.Query("user_id"); userIDStr != "" {
		if userID, err := strconv.ParseInt(userIDStr, 10, 64); err == nil {
			filter.UserID = &userID
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
			return
		}
	}
	if status := c.Query("status"); status != "" {
		filter.Status = &status
	}
	if from := c.Query("from"); from != "" {
		if t, err := time.Parse("2006-01-02", from); err == nil {
			filter.From = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
			return
		}
	}
	if to := c.Query("to"); to != "" {
		if t, err := time.Parse("2006-01-02", to); err == nil {
			filter.To = &t
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date"})
			return
		}
	}

	orders, err := h.orders.GetOrdersFiltered(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) GetOrderByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	order, err := h.orders.GetOrderByID(c.Request.Context(), id)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetOrdersByStatus(c *gin.Context) {
	status := c.Param("status")
	orders, err := h.orders.GetOrdersByStatus(c.Request.Context(), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	var req UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.orders.UpdateOrderStatus(c.Request.Context(), id, req.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}

func (h *OrderHandler) CancelOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	order, err := h.orders.GetOrderByID(c.Request.Context(), id)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if order.Status == "completed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cannot cancel an order that is already shipped (completed)"})
		return
	}
	err = h.orders.UpdateOrderStatus(c.Request.Context(), id, "cancelled")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order status updated", "status": "cancelled"})
}

func (h *OrderHandler) PayOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	order, err := h.orders.GetOrderByID(c.Request.Context(), id)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order can only be paid if status is 'pending'"})
		return
	}
	err = h.orders.UpdateOrderStatus(c.Request.Context(), id, "paid")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}
	go h.kafkaProducer.Publish(
		context.Background(),
		"order.paid",
		map[string]any{
			"order_id": id,
		},
	)
	c.JSON(http.StatusOK, gin.H{"message": "order status updated", "status": "paid"})
}

func (h *OrderHandler) ShipOrder(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	order, err := h.orders.GetOrderByID(c.Request.Context(), id)
	if err != nil || order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if order.Status != "paid" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order can only be shipped if status is 'paid'"})
		return
	}
	err = h.orders.UpdateOrderStatus(c.Request.Context(), id, "completed")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update order status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order status updated", "status": "completed"})
}

func (h *OrderHandler) GetOrderItems(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64) // <-- use "id"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	items, err := h.orders.GetOrderItems(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

func (h *OrderHandler) AddOrderItem(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64) // <-- use "id"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	var req AddOrderItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	product, err := h.orders.GetProductByID(c.Request.Context(), req.ProductID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if product == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	err = h.orders.AddOrderItem(c.Request.Context(), orderID, req.ProductID, req.Quantity, product.Price)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "item added"})
}

func (h *OrderHandler) UpdateOrderItemQuantity(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64) // <-- use "id"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}
	itemID, err := strconv.ParseInt(c.Param("item_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	var req UpdateOrderItemQuantityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.orders.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error fetching order": err.Error()})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can only update items for orders with status 'pending'"})
		return
	}
	err = h.orders.UpdateOrderItemQuantity(c.Request.Context(), orderID, itemID, req.Quantity)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error updating item quantity": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item quantity updated"})
}

func (h *OrderHandler) RemoveOrderItem(c *gin.Context) {
	orderID, err := strconv.ParseInt(c.Param("id"), 10, 64) // <-- use "id"
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order id"})
		return
	}

	order, err := h.orders.GetOrderByID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error fetching order": err.Error()})
		return
	}

	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	if order.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "can only update items for orders with status 'pending'"})
		return
	}

	itemID, err := strconv.ParseInt(c.Param("item_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item id"})
		return
	}
	err = h.orders.RemoveOrderItem(c.Request.Context(), orderID, itemID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error removing item": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "item removed"})
}

