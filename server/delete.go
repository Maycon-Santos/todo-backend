package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/db"
	"github.com/julienschmidt/httprouter"
)

type deleteRequestBody struct {
	ID string `json:"id"`
}

func DeleteHandler(container container.Container) httprouter.Handle {
	var (
		taskRepository db.TaskRepository
	)

	if err := container.Retrieve(&taskRepository); err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		userID := p.ByName("user_id")

		var requestBody deleteRequestBody

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("playload is invalid"))

			return
		}

		err := taskRepository.Delete(userID, requestBody.ID)
		if err != nil {
			log.Fatal(err)
		}
	}
}
