package router

import (
	// Native Go Libs
	errors "errors"
	http "net/http"

	// Project Libs
	auth "hermes-messaging-service/auth"
	models "hermes-messaging-service/models"
	checkers "hermes-messaging-service/validation/checkers"

	// 3rd Party Libs
	gocustomhttpresponse "github.com/terryvogelsang/gocustomhttpresponse"
	logruswrapper "github.com/terryvogelsang/logruswrapper"
)

type (
	// Handler : Custom type to work with CustomHandle wrapper
	Handler func(env *models.Env, w http.ResponseWriter, r *http.Request) error
)

// AddVerneMQACL : Construct and store VerneMQ ACL in database
func AddVerneMQACL(env *models.Env, w http.ResponseWriter, r *http.Request) error {

	// Retrieve token from request header
	token := r.Header.Get("token")

	// Check if token has valid format (According to regex provided by environment variable)
	tokenHasValidFormat := checkers.IsTokenValid(token)

	// If token is not formatted correctly, return an error response
	if !tokenHasValidFormat {
		return errors.New(logruswrapper.CodeInvalidToken)
	}

	// Check authentication with provided endpoint
	MQTTAuthInfos, err := auth.CheckAuthentication(env, token)

	// If an error occurs, token is invalid
	if err != nil {
		return errors.New(logruswrapper.CodeInvalidToken)
	}

	// Construct MQTT User ACL with MQTT Auth Infos
	verneMQACL := models.NewVerneMQACL(MQTTAuthInfos.ClientID, MQTTAuthInfos.Username, MQTTAuthInfos.Password)

	err = env.DB.AddUserACL(verneMQACL)

	if err != nil {
		return errors.New(logruswrapper.CodeInvalidToken)
	}

	log := logruswrapper.NewEntry("MessagingService", "/helloworld", logruswrapper.CodeSuccess)

	gocustomhttpresponse.WriteResponse(nil, log, w)
	return nil
}

// CustomHandle : Custom Handlers Wrapper for API
func CustomHandle(env *models.Env, handlers ...Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range handlers {
			err := h(env, w, r)
			if err != nil {
				errorLog := logruswrapper.NewEntry("MessagingService", "/helloworld", err.Error())
				gocustomhttpresponse.WriteResponse(nil, errorLog, w)
				return
			}
		}
	})
}
