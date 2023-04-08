package task

import (
	"encoding/json"
	"github.com/google/uuid"
	"net/http"
)

func Handler(taskRepo *TaskRepository, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGetAllTasks(taskRepo, w, r)
	case http.MethodPost:
		handleCreateTask(taskRepo, w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleGetAllTasks(taskRepo *TaskRepository, w http.ResponseWriter, r *http.Request) {
	tasks, err := taskRepo.GetAllTasks()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func handleCreateTask(repo *TaskRepository, w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newTask.ID = uuid.New().String()

	err = repo.CreateTask(&newTask)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}
