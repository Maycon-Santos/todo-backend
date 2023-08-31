package main

import (
	"log"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/process"
	"github.com/Maycon-Santos/test-brand-monitor-backend/server"
)

func main() {
	env, err := process.NewEnv()
	if err != nil {
		log.Fatal(err)
	}

	dependenciesContainer := container.New()

	err = dependenciesContainer.Inject(env)
	if err != nil {
		log.Fatal(err)
	}

	err = server.Listen(dependenciesContainer)
	if err != nil {
		log.Fatal(err)
	}
}
