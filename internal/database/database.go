package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	// The keys and values in the map are service-specific.
	Health() map[string]string

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error

	GetSqlDB() (*sql.DB, error)
}

type service struct {
	db *gorm.DB
}

var (
	dburl      = os.Getenv("BLUEPRINT_DB_URL")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	db, err := gorm.Open(sqlite.Open("app3.db"), &gorm.Config{})
	db.Exec("PRAGMA foreign_keys = ON")

	if err != nil {
		log.Fatal(err)
	}

	// Run migrations here
	// db.AutoMigrate()

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) GetSqlDB() (*sql.DB, error) {
	return s.db.DB()
}

// Health checks the health of the database connection by pinging the database.
// It returns a map with keys indicating various health statistics.
func (s *service) Health() map[string]string {
	_, err := s.db.DB()
	if err != nil {
		return map[string]string{
			"status": "unhealthy",
			"error":  fmt.Sprintf("failed to get sql.DB: %v", err),
		}
	}
	return map[string]string{
		"status": "healthy",
	}
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", dburl)
	sqldb, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqldb.Close()
}
