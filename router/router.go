package router

import (

	// Native Go Libs
	fmt "fmt"
	models "hermes-messaging-service/models"
	handlers "hermes-messaging-service/router/handlers"
	http "net/http"

	// Project Libs

	// 3rd Party Libs
	mux "github.com/gorilla/mux"
	cors "github.com/rs/cors"
)

const (
	// PORT : Listening Port
	PORT int = 8085
)

// Listen : Defines all router routing rules and handlers.
// Serves the API at defined port constant.
func Listen(env *models.Env) {

	r := mux.NewRouter().StrictSlash(false)

	v1 := r.PathPrefix("/v1").Subrouter()

	// HelloWorld Endpoint
	helloWorldV1 := v1.PathPrefix("/helloworld").Subrouter()
	helloWorldV1.Handle("", handlers.CustomHandle(env, handlers.HelloWorld)).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedHeaders:   []string{"X-Requested-With"},
		AllowedOrigins:   []string{"http://frontend.localhost"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	})

	http.ListenAndServe(":"+fmt.Sprintf("%d", PORT), corsHandler.Handler(r))
}
