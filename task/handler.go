package task

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func Handler(taskRepo *TaskRepository, c *gin.Context) {
	switch c.Request.Method {
	case http.MethodGet:
		handleGetAllTasks(taskRepo, c)
	case http.MethodPost:
		handleCreateTask(taskRepo, c)
	default:
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
	}
}

func handleGetAllTasks(taskRepo *TaskRepository, c *gin.Context) {
	tasks, err := taskRepo.GetAllTasks()
	if err != nil {
		//http.Error(w, "Internal server error", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	//json.NewEncoder(w).Encode(tasks)

	c.JSON(http.StatusOK, tasks)
}

func handleCreateTask(repo *TaskRepository, c *gin.Context) {
	var newTask Task
	err := json.NewDecoder(c.Request.Body).Decode(&newTask)
	if err != nil {
		//http.Error(w, "Invalid request body", http.StatusBadRequest)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	newTask.ID = uuid.New().String()

	err = repo.CreateTask(&newTask)
	if err != nil {
		//http.Error(w, "Internal server error", http.StatusInternalServerError)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	//w.Header().Set("Content-Type", "application/json")
	//w.WriteHeader(http.StatusCreated)
	//json.NewEncoder(w).Encode(newTask)
	c.JSON(http.StatusCreated, newTask)
}
