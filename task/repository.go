package task

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(collection *mongo.Collection) *TaskRepository {
	return &TaskRepository{
		collection: collection,
	}
}

func (repo *TaskRepository) GetAllTasks() ([]Task, error) {
	cursor, err := repo.collection.Find(context.Background(), bson.D{})
	if err != nil {
		return nil, err
	}

	var tasks []Task
	if err = cursor.All(context.Background(), &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (repo *TaskRepository) CreateTask(task *Task) error {
	_, err := repo.collection.InsertOne(context.Background(), task)
	return err
}
