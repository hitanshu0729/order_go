package handlers

import (
	"github.com/hitanshu0729/order_go/internal/storage/sqlite"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	products *sqlite.Repo
}

func NewProductHandler(products *sqlite.Repo) *ProductHandler {
	return &ProductHandler{products: products}
}

func (h *ProductHandler) RegisterProductRoutes(rg *gin.RouterGroup) {
	products := rg.Group("/products")
	products.GET("", h.GetProducts)
	products.POST("", h.CreateProduct)
	products.GET(":id", h.GetProductByID)
}

type CreateProductRequest struct {
	Name  string `json:"name" binding:"required"`
	Price int64  `json:"price" binding:"required,gt=0"`
	Stock int64  `json:"stock" binding:"required,gte=0"`
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.products.GetProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Creating product: %+v", req)
	err := h.products.CreateProduct(c.Request.Context(), req.Name, req.Price, req.Stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	log.Printf("Product created successfully: %+v", req)
	c.JSON(http.StatusCreated, gin.H{"message": "Product created"})
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}
	product, err := h.products.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, product)
}
