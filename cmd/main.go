package main

import (
	"log"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/db"
	"github.com/Maycon-Santos/test-brand-monitor-backend/process"
	"github.com/Maycon-Santos/test-brand-monitor-backend/server"
)

func main() {
	env, err := process.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	mongoConn := db.NewMongoConnection(env)
	sqliteConn, err := db.NewSQLiteConnection(env)
	if err != nil {
		log.Fatal(err)
	}

	taskRepository := db.NewTaskRepository(mongoConn)
	userRepository := db.NewUserRepository(sqliteConn)

	dependenciesContainer := container.New()

	err = dependenciesContainer.Inject(env, &taskRepository, &userRepository)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Listen(dependenciesContainer)
	if err != nil {
		log.Fatal(err)
	}
}
