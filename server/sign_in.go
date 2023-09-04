package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/db"
	"github.com/Maycon-Santos/test-brand-monitor-backend/process"
	"github.com/Maycon-Santos/test-brand-monitor-backend/server/auth"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/bcrypt"
)

type signInRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signInResponse struct {
	Token string `json:"token"`
}

type signInResponseError struct {
	Error string `json:"error"`
}

func SignInHandler(container container.Container) httprouter.Handle {
	var (
		env            process.Env
		userRepository db.UserRepository
	)

	if err := container.Retrieve(&userRepository, &env); err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		requestBody := signInRequestBody{}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("playload is invalid"))
			return
		}

		user, err := userRepository.Get(requestBody.Username)
		if err != nil {
			if err == sql.ErrNoRows {
				responseBytes, err := json.Marshal(signInResponseError{Error: "Usuário não existe"})
				if err != nil {
					w.WriteHeader(http.StatusUnprocessableEntity)
					return
				}

				w.WriteHeader(http.StatusUnauthorized)
				w.Write(responseBytes)
				return
			}

			w.WriteHeader(http.StatusUnprocessableEntity)

			fmt.Println(err)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
		if err != nil {
			responseBytes, err := json.Marshal(signInResponseError{Error: "Senha incorreta"})
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Write(responseBytes)

			return
		}

		token, err := auth.CreateToken(env.SecretKey, auth.Claims{
			ID: user.ID,
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
			return
		}

		if responseBytes, err := json.Marshal(signInResponse{Token: token}); err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write(responseBytes)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}
