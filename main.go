package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"tracker/task"
)

func main() {
	router := gin.Default()

	client, collection := connectToDB()
	defer client.Disconnect(context.Background())

	taskRepo := task.NewTaskRepository(collection)

	router.GET("/tasks", func(c *gin.Context) {
		task.Handler(taskRepo, c.Writer, c.Request)
	})

	router.POST("/tasks", func(c *gin.Context) {
		task.Handler(taskRepo, c.Writer, c.Request)
	})

	log.Println("Listening on :8080...")
	err := router.Run(":8080")
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func connectToDB() (*mongo.Client, *mongo.Collection) {
	uri := "mongodb://localhost:27017"
	clientOptions := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("task").Collection("tasks")

	return client, collection
}