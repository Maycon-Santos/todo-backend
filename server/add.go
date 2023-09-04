package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/db"
	"github.com/julienschmidt/httprouter"
)

type addRequestBody struct {
	Title string `json:"title"`
}

func AddHandler(container container.Container) httprouter.Handle {
	var (
		taskRepository db.TaskRepository
	)

	if err := container.Retrieve(&taskRepository); err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		userID := p.ByName("user_id")

		var requestBody addRequestBody

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("playload is invalid"))

			return
		}

		id, err := taskRepository.Add(userID, requestBody.Title)
		if err != nil {
			log.Fatal(err)
		}

		response := Task{
			ID:    id,
			Title: requestBody.Title,
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(responseBytes)
	}
}
