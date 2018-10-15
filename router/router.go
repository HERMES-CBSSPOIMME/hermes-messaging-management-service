package router

import (

	// Native Go Libs
	fmt "fmt"
	http "net/http"

	// Project Libs
	models "hermes-messaging-service/models"
	handlers "hermes-messaging-service/router/handlers"

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
	aclV1 := v1.PathPrefix("/profile").Subrouter()
	aclV1.Handle("", handlers.CustomHandle(env, handlers.AddVerneMQACL)).Methods("POST")

	conversationsV1 := v1.PathPrefix("/conversations").Subrouter()
	conversationsV1.Handle("/group", handlers.CustomHandle(env, handlers.AddGroupConversation)).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedHeaders:   []string{"X-Requested-With"},
		AllowedOrigins:   []string{"http://frontend.localhost"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	})

	http.ListenAndServe(":"+fmt.Sprintf("%d", PORT), corsHandler.Handler(r))
}
