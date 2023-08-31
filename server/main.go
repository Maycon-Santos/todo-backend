package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/process"
	"github.com/julienschmidt/httprouter"
)

func Listen(container container.Container) error {
	var env process.Env

	err := container.Retrieve(&env)
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()

	corsMiddleware := CORSMiddleware(container)

	router.GlobalOPTIONS = http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get("Access-Control-Request-Method") != "" {
			header := writer.Header()
			header.Set("Access-Control-Allow-Origin", env.AccessControlAllowOrigin)
			header.Set("Access-Control-Allow-Headers", env.AccessControlAllowHeaders)
			header.Set("Access-Control-Allow-Methods", header.Get("Allow"))
		}

		writer.WriteHeader(http.StatusNoContent)
	})

	router.GET("/", corsMiddleware(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		w.Write([]byte("Hello World"))
	}))

	return http.ListenAndServe(fmt.Sprintf(":%d", env.ServerPort), router)
}
