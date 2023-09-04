package auth

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/db"
	"github.com/Maycon-Santos/test-brand-monitor-backend/process"
	"github.com/golang-jwt/jwt/v5"
	"github.com/julienschmidt/httprouter"
)

func GetDataMiddleware(container container.Container) func(next httprouter.Handle) httprouter.Handle {
	var (
		env            process.Env
		userRepository db.UserRepository
	)

	err := container.Retrieve(&env, &userRepository)
	if err != nil {
		log.Fatal(err)
	}

	return func(next httprouter.Handle) httprouter.Handle {
		return func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
			tokenStr := request.Header.Get("Token")
			token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
				return []byte(env.SecretKey), nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					writer.WriteHeader(http.StatusUnauthorized)
					return
				}
				writer.WriteHeader(http.StatusBadRequest)
				return
			}
			if !token.Valid {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				user, err := userRepository.GetByID(int64(claims["id"].(float64)))
				if err != nil {
					panic(err)
				}

				params = append(params, httprouter.Param{
					Key:   "user_id",
					Value: fmt.Sprintf("%d", user.ID),
				})

				params = append(params, httprouter.Param{
					Key:   "username",
					Value: string(user.Username),
				})
			}

			next(writer, request, params)
		}
	}
}
