package task

import (
	"context"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) (*mongo.Client, *mongo.Collection) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		t.Fatalf("Failed to ping MongoDB: %v", err)
	}

	db := client.Database("task_tracker_test")
	collection := db.Collection("tasks")

	return client, collection
}

func TestFindAll(t *testing.T) {
	client, collection := setupTestDB(t)
	defer client.Disconnect(context.Background())

	// Clean up the test collection
	collection.Drop(context.Background())

	// Insert some sample tests
	sampleTasks := []Task{
		{
			ID:   uuid.New().String(),
			Name: "Task 1",
		},
		{
			ID:   uuid.New().String(),
			Name: "Task 2",
		},
	}

	// Convert []Task to []interface{}
	interfaceTasks := make([]interface{}, len(sampleTasks))
	for i, task := range sampleTasks {
		interfaceTasks[i] = task
	}

	_, err := collection.InsertMany(context.Background(), interfaceTasks)
	if err != nil {
		t.Fatalf("Failed to insert sample tasks: %v", err)
	}

	repo := NewTaskRepository(collection)

	tasks, err := repo.GetAllTasks(context.Background())
	if err != nil {
		t.Fatalf("Failed to call FindAll: %v", err)
	}

	if len(tasks) != 2 {
		t.Fatalf("Expected 2 tasks, got %d", len(tasks))
	}
}

func TestCreateTask(t *testing.T) {
	client, collection := setupTestDB(t)
	defer client.Disconnect(context.Background())

	// Clean up the test collection
	collection.Drop(context.Background())

	repo := NewTaskRepository(collection)

	newTask := Task{
		ID:   uuid.New().String(),
		Name: "Task 4",
	}

	err := repo.CreateTask(context.Background(), &newTask)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	foundTask := Task{}
	err = collection.FindOne(context.Background(), bson.M{"_id": newTask.ID}).Decode(&foundTask)
	if err != nil {
		t.Fatalf("Failed to find created task: %v", err)
	}

	if foundTask.ID != newTask.ID || foundTask.Name != newTask.Name {
		t.Errorf("Task not created correclty: expected %v, got %v", newTask, foundTask)
	}
}

func TestUpdateTask(t *testing.T) {
	client, collection := setupTestDB(t)
	defer client.Disconnect(context.Background())

	// Clean up the test collection
	collection.Drop(context.Background())

	repo := NewTaskRepository(collection)

	// Create a task first
	newTask := Task{
		ID:   uuid.New().String(),
		Name: "Task to be updated",
	}
	err := repo.CreateTask(context.Background(), &newTask)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Prepare the updated task
	updatedTask := Task{
		Name:      "Updated task",
		Completed: true,
	}

	// Update the task
	err = repo.UpdateTask(context.Background(), newTask.ID, &updatedTask)
	if err != nil {
		t.Fatalf("Failed to update task: %v", err)
	}

	// Find the updated task
	var foundTask bson.M
	err = collection.FindOne(context.Background(), bson.M{"_id": newTask.ID}).Decode(&foundTask)
	if err != nil {
		t.Fatalf("Failed to find updated task: %v", err)
	}

	// Check the updated fields
	if newTask.ID != foundTask["_id"].(string) {
		t.Errorf("Task ID not updated correctly: expected %v, fot %v", newTask.ID, foundTask["_id"].(string))
	}
	if updatedTask.Name != foundTask["name"].(string) {
		t.Errorf("Task name not updated correctly: expected %v, got %v", updatedTask.Name, foundTask["name"].(string))
	}
	if updatedTask.Completed != foundTask["completed"].(bool) {
		t.Errorf("Task completion status not updated correctly: expected %v, got %v", updatedTask.Completed, foundTask["completed"].(bool))
	}
}
