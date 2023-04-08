package task

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type TaskRepository interface {
	GetAllTasks(ctx context.Context) ([]Task, error)
	CreateTask(ctx context.Context, task *Task) error
	UpdateTask(ctx context.Context, id string, task *Task) error
}

type MongoTaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) *MongoTaskRepository {
	return &MongoTaskRepository{
		collection: collection,
	}
}

func (repo *MongoTaskRepository) GetAllTasks(ctx context.Context) ([]Task, error) {
	cursor, err := repo.collection.Find(ctx, bson.D{})
	if err != nil {
		log.Printf("Error fetching tasks from the database: %v", err)
		return nil, err
	}

	var tasks []Task
	if err = cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *MongoTaskRepository) CreateTask(ctx context.Context, task *Task) error {
	_, err := repo.collection.InsertOne(ctx, task)
	if err != nil {
		log.Printf("Error inserting task into the database: %v", err)
		return err
	}

	return nil
}

func (repo *MongoTaskRepository) UpdateTask(ctx context.Context, id string, task *Task) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"name": task.Name, "completed": task.Completed}}

	result, err := repo.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating task in the database: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		return ErrTaskNotFound
	}

	return nil
}
