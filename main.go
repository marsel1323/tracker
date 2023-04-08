package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"tracker/task"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	client, collection := connectToDB()
	defer client.Disconnect(context.Background())

	taskRepo := task.NewTaskRepository(collection)

	router := gin.Default()

	router.GET("/tasks", func(c *gin.Context) {
		task.Handler(taskRepo, c)
	})

	router.POST("/tasks", func(c *gin.Context) {
		task.Handler(taskRepo, c)
	})

	log.Println("Listening on :8080...")
	err = router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectToDB() (*mongo.Client, *mongo.Collection) {
	uri := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	collection := client.Database("task").Collection("tasks")

	return client, collection
}
