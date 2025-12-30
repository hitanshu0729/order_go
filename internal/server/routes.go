package server

import (
	"github.com/hitanshu0729/order_go/internal/handlers"
	"github.com/hitanshu0729/order_go/internal/storage/sqlite"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	api := r.Group("/api/v1")
	api.GET("/", s.HelloWorldHandler)
	api.GET("/health", s.healthHandler)

	// User Routes
	sqlDB, err := s.db.GetSqlDB()
	if err != nil {
		log.Fatal("Failed to get SQL DB:", err)
	}

	Repo := sqlite.NewRepo(sqlDB)
	userHandler := handlers.NewUserHandler(Repo)
	userHandler.RegisterUserRoutes(api)

	// Product Routes
	productHandler := handlers.NewProductHandler(Repo)
	productHandler.RegisterProductRoutes(api)

	// Order Routes
	orderHandler := handlers.NewOrderHandler(Repo)
	orderHandler.RegisterOrderRoutes(api)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
