package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Maycon-Santos/test-brand-monitor-backend/container"
	"github.com/Maycon-Santos/test-brand-monitor-backend/process"
	"github.com/Maycon-Santos/test-brand-monitor-backend/server/auth"
	"github.com/julienschmidt/httprouter"
)

func Listen(container container.Container) error {
	var env process.Env

	err := container.Retrieve(&env)
	if err != nil {
		log.Fatal(err)
	}

	router := httprouter.New()

	authGetDataMiddleware := auth.GetDataMiddleware(container)
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

	router.POST("/sign_in", corsMiddleware(SignInHandler(container)))
	router.POST("/sign_up", corsMiddleware(SignUpHandler(container)))
	router.GET("/user_data", corsMiddleware(authGetDataMiddleware(UserDataHandler(container))))
	router.GET("/is_authenticated", corsMiddleware(isAuthenticatedHandler(container)))

	router.POST("/add", corsMiddleware(authGetDataMiddleware(AddHandler(container))))
	router.GET("/list", corsMiddleware(authGetDataMiddleware(ListHandler(container))))
	router.POST("/edit", corsMiddleware(authGetDataMiddleware(EditHandler(container))))
	router.DELETE("/delete", corsMiddleware(authGetDataMiddleware(DeleteHandler(container))))

	fmt.Println(fmt.Sprintf("Server listening on port: %d", env.ServerPort))

	return http.ListenAndServe(fmt.Sprintf(":%d", env.ServerPort), router)
}
