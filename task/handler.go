package task

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
)

var ErrTaskNotFound = errors.New("task not found")

func RegisterRoutes(taskRepo *MongoTaskRepository, router *gin.Engine) {
	router.GET("/tasks", func(c *gin.Context) {
		handleGetAllTasks(taskRepo, c)
	})

	router.POST("/tasks", func(c *gin.Context) {
		handleCreateTask(taskRepo, c)
	})

	router.PUT("/tasks/:id", func(c *gin.Context) {
		handleUpdateTask(taskRepo, c)
	})
}

func handleGetAllTasks(taskRepo *MongoTaskRepository, c *gin.Context) {
	tasks, err := taskRepo.GetAllTasks(c.Request.Context())
	if err != nil {
		log.Printf("Error fetching tasks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func handleCreateTask(repo *MongoTaskRepository, c *gin.Context) {
	var newTask Task
	err := json.NewDecoder(c.Request.Body).Decode(&newTask)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	newTask.ID = uuid.New().String()

	err = repo.CreateTask(c.Request.Context(), &newTask)
	if err != nil {
		log.Printf("Error creating task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusCreated, newTask)
}

func handleUpdateTask(repo *MongoTaskRepository, c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id URL parameter"})
		return
	}

	var updatedTask Task
	err := json.NewDecoder(c.Request.Body).Decode(&updatedTask)
	if err != nil {
		log.Printf("Error decoding request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	updatedTask.ID = id
	err = repo.UpdateTask(c.Request.Context(), id, &updatedTask)
	if err != nil {
		if err == ErrTaskNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
			return
		}

		log.Printf("Error updating task: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, updatedTask)
}
