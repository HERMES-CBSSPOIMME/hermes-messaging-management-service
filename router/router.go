package router

import (
	fmt "fmt"
	http "net/http"
	models "wave-messaging-management-service/models"
	handlers "wave-messaging-management-service/router/handlers"

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
	aclV1 := v1.PathPrefix("/profiles").Subrouter()
	aclV1.Handle("", handlers.CustomHandle(env, handlers.AddVerneMQACL)).Methods("POST")
	aclV1.Handle("/mappings", handlers.CustomHandle(env, handlers.GetMappingForUsers)).Methods("POST")

	conversationsV1 := v1.PathPrefix("/conversations").Subrouter()
	conversationsV1.Handle("/group", handlers.CustomHandle(env, handlers.AddGroupConversation)).Methods("POST")

	corsHandler := cors.New(cors.Options{
		AllowedHeaders:   []string{"X-Requested-With"},
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "HEAD", "POST", "PUT", "OPTIONS"},
	})

	http.ListenAndServe(":"+fmt.Sprintf("%d", PORT), corsHandler.Handler(r))
}
