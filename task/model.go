package task

type Task struct {
	ID        string `json:"id" bson:"_id"`
	Name      string `json:"name" bson:"name"`
	Completed bool   `json:"completed" bson:"completed"`
	Link      string `json:"link" bson:"link"`
	TimeSpent int    `json:"timeSpent" bson:"timeSpent"`
	TimeGiven int    `json:"timeGiven" bson:"timeGiven"`
	Status    string `json:"status" bson:"status"`
	TaskType  string `json:"taskType" bson:"taskType"`
	CreatedAt string `json:"createdAt" bson:"createdAt"`
	UpdatedAt string `json:"updatedAt" bson:"updatedAt"`
}
