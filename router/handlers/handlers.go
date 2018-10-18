package router

import (

	// Native Go Libs
	json "encoding/json"
	errors "errors"
	"log"
	http "net/http"

	// Project Libs
	auth "hermes-messaging-service/auth"
	models "hermes-messaging-service/models"
	utils "hermes-messaging-service/utils"
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
	tokenHasValidFormat, err := checkers.IsTokenValid(env, token)

	if err != nil {
		return err
	}

	// If token is not formatted correctly, return an error response
	if !tokenHasValidFormat {
		log.Println("Invalid token format")
		return errors.New(logruswrapper.CodeInvalidToken)
	}

	// Check authentication with provided endpoint
	MQTTAuthInfos, wasCached, wasTokenUpdated, err := auth.CheckAuthentication(env, token)

	// If an error occurs, token is invalid
	if err != nil {
		log.Println(err)
		return errors.New(logruswrapper.CodeInvalidToken)
	}

	if wasTokenUpdated {
		log.Println("Token Updated")
		return errors.New(logruswrapper.CodeUpdated)
	}

	if wasCached {
		log.Println("Already cached")
		return errors.New(logruswrapper.CodeAlreadyExists)
	}

	// Construct MQTT User ACL with MQTT Auth Infos + default ACLs
	verneMQACL := models.NewVerneMQACL(MQTTAuthInfos.ClientID, MQTTAuthInfos.Username, MQTTAuthInfos.Password)

	err = env.MongoDB.AddProfileACL(verneMQACL)

	if err != nil {
		log.Println(err)
		return errors.New(logruswrapper.CodeInvalidToken)
	}

	log := logruswrapper.NewEntry("MessagingService", "/helloworld", logruswrapper.CodeSuccess)

	gocustomhttpresponse.WriteResponse(MQTTAuthInfos.ClientID, log, w)
	return nil
}

// AddGroupConversation : Add group conversation ACLs in database
func AddGroupConversation(env *models.Env, w http.ResponseWriter, r *http.Request) error {

	// Retrieve token from request header
	token := r.Header.Get("token")

	// Check if token has valid format (According to regex provided by environment variable)
	tokenHasValidFormat, err := checkers.IsTokenValid(env, token)

	if err != nil {
		return err
	}

	// If token is not formatted correctly, return an error response
	if !tokenHasValidFormat {
		return errors.New(logruswrapper.CodeInvalidToken)
	}

	// Check authentication with provided endpoint
	MQTTAuthInfos, _, _, err := auth.CheckAuthentication(env, token)

	// If an error occurs, token is invalid
	if err != nil {
		return errors.New(logruswrapper.CodeInvalidToken)
	}

	reqBody := utils.GroupConversationBody{}
	err = json.NewDecoder(r.Body).Decode(&reqBody)

	if err != nil {
		// TODO: Change this to invalid JSON
		return errors.New(logruswrapper.CodeInvalidJSON)
	}

	// Create new group conversation struct
	groupConv := models.NewGroupConversation(reqBody.Name, append(reqBody.Members, MQTTAuthInfos.ClientID))

	// TODO: Add check if provided users exist

	// Store conversation infos in DB
	err = env.MongoDB.AddGroupConversation(groupConv)

	// Update ACL in DB (Request maker get publish rights on recipient private topic)
	err = env.MongoDB.UpdateProfilesWithGroupACL(groupConv)

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
