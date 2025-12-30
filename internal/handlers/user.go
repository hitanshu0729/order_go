package handlers

import (
	"github.com/hitanshu0729/order_go/internal/storage/sqlite"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandlers holds user-related handlers.
type UserHandler struct {
	users *sqlite.Repo
}

func NewUserHandler(users *sqlite.Repo) *UserHandler {
	return &UserHandler{users: users}
}

// RegisterUserRoutes registers user routes under the given router group.
func (h *UserHandler) RegisterUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	users.GET("", h.GetUsers)
	users.POST("", h.CreateUser)
	users.GET("/:id", h.GetUserByID)
	users.PATCH("/:id", h.UpdateUser)
	users.DELETE("/:id", h.DeleteUser)
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.users.GetUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("Creating user: %+v", req)

	err := h.users.CreateUser(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Printf("User created successfully: %+v", req)

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id") // always string

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	user, err := h.users.GetUserByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id") // always string
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}

	err = h.users.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id") // always string
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id",
		})
		return
	}
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.users.UpdateUser(c.Request.Context(), id, req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to update user",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user updated"})
}
