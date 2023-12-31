package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type CheckItem struct {
	ID    string
	Label string
	Done  bool
}

type Task struct {
	ID          string `bson:"_id,omitempty"`
	UserID      string
	Done        bool
	Title       string
	Description string
	Checklist   []CheckItem
}

type TaskRepository interface {
	GetAll(userID string) (results []Task, err error)
	GetByID(userID string, id string) (*Task, error)
	Add(userID string, title string) (id string, err error)
	Edit(userID string, task Task) error
	Delete(userID string, id string) error
}

type taskRepository struct {
	collection *mongo.Collection
}

func NewTaskRepository(conn *mongo.Client) TaskRepository {
	collection := conn.Database("TestBrandMonitor").Collection("task")

	return taskRepository{
		collection: collection,
	}
}

func (u taskRepository) GetAll(userID string) (results []Task, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}})

	cursor, err := u.collection.Find(ctx, bson.D{{Key: "userID", Value: userID}}, opts)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx, &results); err != nil {
		panic(err)
	}

	for _, result := range results {
		cursor.Decode(&result)
	}

	return results, nil
}

func (t taskRepository) GetByID(userID string, id string) (*Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var task Task

	err := t.collection.FindOne(ctx, bson.M{"_id": id, "userID": userID}).Decode(&task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (t taskRepository) Add(userID string, title string) (id string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := t.collection.InsertOne(
		ctx,
		bson.D{
			{Key: "userID", Value: userID},
			{Key: "title", Value: title},
			{Key: "createdAt", Value: primitive.NewDateTimeFromTime(time.Now())},
		},
	)
	if err != nil {
		return "", err
	}

	return res.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (t taskRepository) Edit(userID string, task Task) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(task.ID)
	if err != nil {
		return err
	}

	task.ID = ""

	_, err = t.collection.UpdateOne(
		ctx,
		bson.D{
			{
				Key:   "_id",
				Value: objID,
			},
			{
				Key:   "userID",
				Value: userID,
			},
		},
		bson.D{
			{
				Key:   "$set",
				Value: task,
			},
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (t taskRepository) Delete(userID string, id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = t.collection.DeleteOne(ctx, bson.M{"_id": objID, "userID": userID})
	if err != nil {
		return err
	}

	return nil
}
