package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/db"
	"github.com/julienschmidt/httprouter"
)

type listResponse struct {
	Results []Task `json:"results"`
}

func ListHandler(container container.Container) httprouter.Handle {
	var (
		taskRepository db.TaskRepository
	)

	if err := container.Retrieve(&taskRepository); err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		taskList, err := taskRepository.GetAll()
		if err != nil {
			log.Fatal(err)
		}

		response := listResponse{}

		for _, item := range taskList {
			task := Task{
				ID:          item.ID,
				Done:        item.Done,
				Title:       item.Title,
				Description: item.Description,
			}

			for _, checkItem := range item.Checklist {
				task.Checklist = append(task.Checklist, CheckItem{
					ID:    checkItem.ID,
					Label: checkItem.Label,
					Done:  checkItem.Done,
				})
			}

			response.Results = append(response.Results, task)
		}

		responseBytes, err := json.Marshal(response)
		if err != nil {
			log.Fatal(err)
		}

		w.Write(responseBytes)
	}
}
