package db

import (
	"context"
	"log"
	"time"

	"github.com/Maycon-Santos/test-brand-monitor-backend/process"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoConnection(env *process.Env) *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(env.MongoDbConnUri))
	if err != nil {
		log.Fatal(err)
	}

	return client
}
