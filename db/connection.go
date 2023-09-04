package db

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/Maycon-Santos/test-brand-monitor-backend/process"
	_ "github.com/mattn/go-sqlite3"
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

func NewSQLiteConnection(env *process.Env) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", env.DBFile)
	if err != nil {
		return nil, err
	}

	return db, nil
}
