package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/db"
	"github.com/julienschmidt/httprouter"
)

func EditHandler(container container.Container) httprouter.Handle {
	var (
		taskRepository db.TaskRepository
	)

	if err := container.Retrieve(&taskRepository); err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		userID := p.ByName("user_id")

		var requestBody Task

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("playload is invalid"))

			return
		}

		task := db.Task{
			ID:          requestBody.ID,
			Done:        requestBody.Done,
			Title:       requestBody.Title,
			Description: requestBody.Description,
		}

		for _, checkItem := range requestBody.Checklist {
			task.Checklist = append(task.Checklist, db.CheckItem{
				ID:    checkItem.ID,
				Label: checkItem.Label,
				Done:  checkItem.Done,
			})
		}

		err := taskRepository.Edit(userID, task)
		if err != nil {
			log.Fatal(err)
		}

		responseBytes, err := json.Marshal(requestBody)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(responseBytes)
	}
}
