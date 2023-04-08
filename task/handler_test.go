package task

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupHandler(t *testing.T) (*mongo.Client, *MongoTaskRepository) {
	client, collection := setupTestDB(t)
	taskRepo := NewTaskRepository(collection)
	return client, taskRepo
}

func tearDownHandler(t *testing.T, client *mongo.Client, taskRepo *MongoTaskRepository) {
	db := taskRepo.collection.Database()
	err := db.Drop(context.Background())
	if err != nil {
		t.Fatalf("Failed to drop test database: %v", err)
	}
	client.Disconnect(context.Background())
}

func performRequest(handler gin.HandlerFunc, method, url string, requestBody io.Reader) *httptest.ResponseRecorder {
	router := gin.Default()
	router.Handle(method, "tasks/:id", handler)

	req := httptest.NewRequest(method, url, requestBody)
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	if method == "PUT" {
		parts := strings.Split(url, "/")
		id := parts[len(parts)-1]
		c.Params = append(c.Params, gin.Param{
			Key:   "id",
			Value: id,
		})
	}

	handler(c)
	return resp
}

func TestHandleGetAllTasks(t *testing.T) {
	client, taskRepo := setupHandler(t)
	defer tearDownHandler(t, client, taskRepo)

	handler := func(c *gin.Context) {
		handleGetAllTasks(taskRepo, c)
	}
	rr := performRequest(handler, "GET", "/tasks", nil)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	sampleTasks := []Task{
		{
			Name: "Task 4",
		},
	}

	var tasks []Task
	err := json.Unmarshal(rr.Body.Bytes(), &tasks)
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

	handler := func(c *gin.Context) {
		handleCreateTask(taskRepo, c)
	}
	rr := performRequest(handler, "POST", "/tasks", bytes.NewReader(jsonTask))

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

func TestHandleUpdateTask(t *testing.T) {
	client, taskRepo := setupHandler(t)
	defer tearDownHandler(t, client, taskRepo)

	// Create a task first
	newTask := Task{
		ID:        uuid.New().String(),
		Name:      "Task to be updated",
		Completed: false,
	}

	err := taskRepo.CreateTask(context.Background(), &newTask)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Prepare the updated task
	updatedTask := Task{
		ID:        newTask.ID,
		Name:      "Updated Task",
		Completed: true,
	}

	jsonTask, err := json.Marshal(updatedTask)
	if err != nil {
		t.Fatal(err)
	}

	handler := func(c *gin.Context) {
		handleUpdateTask(taskRepo, c) // Change this line
	}
	rr := performRequest(handler, "PUT", "/tasks/"+newTask.ID, bytes.NewReader(jsonTask))

	if status := rr.Code; status != http.StatusOK {
		t.Fatalf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedUpdatedTask Task
	err = json.Unmarshal(rr.Body.Bytes(), &returnedUpdatedTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	assert.Equal(t, newTask.ID, returnedUpdatedTask.ID)
	assert.Equal(t, updatedTask.Name, returnedUpdatedTask.Name)
	assert.Equal(t, updatedTask.Completed, returnedUpdatedTask.Completed)

	var foundTask bson.M
	err = taskRepo.collection.FindOne(
		context.Background(),
		bson.M{"_id": newTask.ID},
	).Decode(&foundTask)
	if err != nil {
		t.Fatalf("Failed to find updated task: %v", err)
	}

	assert.Equal(t, newTask.ID, foundTask["_id"].(string))
	assert.Equal(t, updatedTask.Name, foundTask["name"].(string))
	assert.Equal(t, updatedTask.Completed, foundTask["completed"].(bool))
}
