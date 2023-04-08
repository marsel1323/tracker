package task

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupHandler(t *testing.T) (*mongo.Client, *TaskRepository) {
	client, collection := setupTestDB(t)
	taskRepo := NewTaskRepository(collection)
	return client, taskRepo
}

func tearDownHandler(t *testing.T, client *mongo.Client, taskRepo *TaskRepository) {
	db := taskRepo.collection.Database()
	err := db.Drop(context.Background())
	if err != nil {
		t.Fatalf("Failed to drop test database: %v", err)
	}
	client.Disconnect(context.Background())
}

type testHandler struct {
	taskRepo *TaskRepository
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	Handler(h.taskRepo, w, r)
}

func TestHandleGetAllTasks(t *testing.T) {
	client, taskRepo := setupHandler(t)
	defer tearDownHandler(t, client, taskRepo)

	req, err := http.NewRequest("GET", "/tasks", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := &testHandler{taskRepo: taskRepo}
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	sampleTasks := []Task{
		{
			Name: "Task 4",
		},
	}

	var tasks []Task
	err = json.Unmarshal(rr.Body.Bytes(), &tasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if len(tasks) != len(sampleTasks) {
		t.Errorf("Unexpected number of tasks: got %v want %v", len(tasks), len(sampleTasks))
	}
}

func TestHandleCreateTask(t *testing.T) {
	client, taskRepo := setupHandler(t)
	defer tearDownHandler(t, client, taskRepo)

	newTask := Task{
		Name: "New Task",
	}

	jsonTask, err := json.Marshal(newTask)
	if err != nil {
		t.Fatal(err)
	}

	req, err := http.NewRequest("POST", "/tasks", bytes.NewBuffer(jsonTask))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := &testHandler{taskRepo: taskRepo}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var createdTask Task
	err = json.Unmarshal(rr.Body.Bytes(), &createdTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.NotEmpty(t, createdTask.ID)
	assert.Equal(t, newTask.Name, createdTask.Name)

	var foundTask bson.M
	err = taskRepo.collection.FindOne(context.Background(), bson.M{"_id": createdTask.ID}).Decode(&foundTask)
	if err != nil {
		t.Fatalf("Failed to find created task: %v", err)
	}
}
