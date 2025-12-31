package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hitanshu0729/order_go/internal/kafka"
	_ "github.com/joho/godotenv/autoload"

	"github.com/hitanshu0729/order_go/internal/database"
)

type Server struct {
	port int

	db database.Service

	KafkaProducer *kafka.Producer

	DLQProducer *kafka.DLQProducer
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	log.Println("Creating kafka producer")
	producer := kafka.NewProducer([]string{"localhost:9092"})
	// defer producer.Close()
	log.Println("Kafka producer created successfully.")

	dlqproducer := kafka.NewDLQProducer([]string{"localhost:9092"})
	// defer dlqproducer.Close()
	log.Println("Kafka DLQ producer created successfully.")

	log.Println("Started kafka consumer")
	go func() {
		time.Sleep(3 * time.Second)
		err := producer.Publish(
			context.Background(),
			"test.message",
			map[string]string{"msg": "hello from startup"},
		)
		if err != nil {
			log.Println("startup publish failed:", err)
		} else {
			log.Println("startup publish success")
		}
	}()

	NewServer := &Server{
		port: port,

		db: database.New(),

		KafkaProducer: producer,

		DLQProducer: dlqproducer,
	}

	log.Println("Database connected successfully.")

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Printf("Server is running on port %d", NewServer.port)

	return server
}
