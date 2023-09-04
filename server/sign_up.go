package server

import (
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

type signUpRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type signUpResponse struct {
	Token string `json:"token"`
}

type signUpResponseError struct {
	Error string `json:"error"`
}

func SignUpHandler(container container.Container) httprouter.Handle {
	var (
		env            process.Env
		userRepository db.UserRepository
	)

	if err := container.Retrieve(&userRepository, &env); err != nil {
		log.Fatal(err)
	}

	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		requestBody := signUpRequestBody{}

		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			w.Write([]byte("playload is invalid"))
			return
		}

		if len(requestBody.Username) < 4 {
			responseBytes, err := json.Marshal(signUpResponseError{Error: "O nome de usuário deve conter, no mínimo, 4 caracteres"})
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Write(responseBytes)

			return
		}

		if len(requestBody.Password) < 4 {
			responseBytes, err := json.Marshal(signUpResponseError{Error: "A senha deve conter, no mínimo, 4 caracteres"})
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Write(responseBytes)

			return
		}

		password, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), env.PasswordsEncryptCost)
		if err != nil {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Println(err)
			return
		}

		userID, err := userRepository.Add(requestBody.Username, string(password))
		if err != nil {
			responseBytes, err := json.Marshal(signUpResponseError{Error: "Usuário já existe"})
			if err != nil {
				w.WriteHeader(http.StatusUnprocessableEntity)
				return
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Write(responseBytes)

			fmt.Println(err)
			return
		}

		token, err := auth.CreateToken(env.SecretKey, auth.Claims{
			ID: userID,
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			fmt.Println(err)
			return
		}

		if responseBytes, err := json.Marshal(signUpResponse{Token: token}); err == nil {
			w.Write(responseBytes)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}
