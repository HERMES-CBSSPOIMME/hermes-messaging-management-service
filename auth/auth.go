package auth

import (
	"github.com/satori/go.uuid"
	// Native Go Libs

	errors "errors"
	http "net/http"

	// Project Libs
	models "hermes-messaging-service/models"
)

// CheckAuthentication : When a token is received :
// (1) Check if token is cached in redis (GET session:{token})
// -> If not it might be a new token for an already checked-in user
//
// (2) Check if originalUserID returned by a GET HTTP Request on provided auth endpoint with token as header value is already matched with one token (HGET mapping:originalUserID token)
// -> If yes, replace session:{oldToken}:{internalHermesUserID} by new token (RENAME session:{oldToken} session:{newToken}) in redis
// ---> Update mapping in redis (HSET mapping:originalUserID token {newToken} )
// -> If no, generate an internal hermesUserID
// ---> Cache it in redis (SET session{newToken}:{internalHermesUserID})
// ---> Set mapping in redis (HSET mapping:originalUserID token {newToken} internalHermesUserID: {internalHermesUserID} )
//
// Return MQTT Auth Infos if provided auth token is valid, an error otherwise
func CheckAuthentication(env *models.Env, token string) (*models.MQTTAuthInfos, error) {

	// If no token, return an error
	if token == "" {
		return nil, errors.New("No Token Provided")
	}

	client := &http.Client{}

	// Refresh config to get actual environment variables
	env.RefreshConfig()

	// Init request
	req, err := http.NewRequest("GET", env.Config.AuthenticationCheckEndpoint, nil)

	if err != nil {
		return nil, err
	}

	// Add token header
	req.Header.Add("token", token)

	// Execute request
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	// Authentication endpoint should return the following if token is invalid :
	//
	// HTTP Status Code : 400 (Bad Request)
	// Empty Body
	if res.StatusCode == 400 {
		return nil, errors.New("Invalid Token")
	}

	// Authentication endpoint should return the following if token is valid :
	//
	// HTTP Status Code : 200 (OK)
	// Header(s) : content-type:application/json
	// Empty Body
	if res.StatusCode == 200 {

		MQTTAuthInfos := models.MQTTAuthInfos{}

		clientID := uuid.NewV4().String()

		// Set new UUID V4
		MQTTAuthInfos.ClientID = clientID

		// Copy ClientID as Username
		MQTTAuthInfos.Username = clientID

		// Set token as password
		MQTTAuthInfos.Password = token

		// TODO: Store session:{token} -> {hermesUserID} in redis
		// - Store mapping:originalUserID -> [token, hermesUserID]
		return &MQTTAuthInfos, nil
	}

	return nil, errors.New("Invalid Token")
}
