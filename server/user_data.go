package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/process"
	"github.com/julienschmidt/httprouter"
)

type userDataResponse struct {
	Username string `json:"username"`
}

func UserDataHandler(container container.Container) httprouter.Handle {
	var (
		env process.Env
	)

	if err := container.Retrieve(&env); err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		username := p.ByName("username")

		if responseBytes, err := json.Marshal(userDataResponse{Username: username}); err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(responseBytes)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}
