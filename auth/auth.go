package auth

import (
	"encoding/json"
	"fmt"

	"github.com/satori/go.uuid"

	// Native Go Libs

	errors "errors"
	http "net/http"

	// Project Libs
	models "hermes-messaging-service/models"
	"hermes-messaging-service/utils"
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

	// Check if token is cached in redis (GET session:{token})
	cachedUserID, err := env.Redis.Get(fmt.Sprintf("session:%s", token))

	if cachedUserID != nil {
		return models.NewMQTTAuthInfos(string(cachedUserID), token), nil
	}

	// Refresh config to get actual environment variables (Auth Endpoint in our case)
	env.RefreshConfig()

	// Create HTTP Client
	client := &http.Client{}

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
	// {"UserID": "userID"}
	if res.StatusCode == 200 {

		// Parse response to get Original User ID
		authCheckerBody := utils.AuthCheckerBody{}
		json.NewDecoder(res.Body).Decode(&authCheckerBody)

		// Check if originalUserID returned is already matched with one token (HGET mapping:originalUserID token)
		cachedOldToken, _ := env.Redis.HGet(fmt.Sprintf("mapping:%s", authCheckerBody.OriginalUserID), "token")

		// -> If yes,
		if cachedOldToken != nil {

			// Replace session:{oldToken}:{internalHermesUserID} by new token (RENAME session:{oldToken} session:{newToken}) in redis
			err = env.Redis.Rename(fmt.Sprintf("session:%s", cachedOldToken), fmt.Sprintf("session:%s", token))

			if err != nil {
				fmt.Println(err)
				return nil, err
			}

			// Update mapping in redis (HSET mapping:originalUserID token {newToken} )
			err = env.Redis.HSet(fmt.Sprintf("mapping:%s", authCheckerBody.OriginalUserID), "token", []byte(token))

			if err != nil {
				fmt.Println(err)
				return nil, err
			}
		}
		// -> If no, generate an internal hermesUserID
		internalHermesUserID := uuid.NewV4().String()

		// ---> Cache it in redis (SET session{newToken}:{internalHermesUserID})
		err = env.Redis.Set(fmt.Sprintf("session:%s", token), []byte(internalHermesUserID))

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// ---> Set mapping in redis (HSET mapping:originalUserID token {newToken} internalHermesUserID: {internalHermesUserID} )
		// TODO: Change Hset method to be able to set multiple field at a time
		err = env.Redis.HSet(fmt.Sprintf("mapping:%s", authCheckerBody.OriginalUserID), "token", []byte(token))

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		err = env.Redis.HSet(fmt.Sprintf("mapping:%s", authCheckerBody.OriginalUserID), "internalHermesUserID", []byte(internalHermesUserID))

		if err != nil {
			fmt.Println(err)
			return nil, err
		}

		// Construct MQTTAuthInfos struct
		MQTTAuthInfos := models.NewMQTTAuthInfos(internalHermesUserID, token)

		return MQTTAuthInfos, nil
	}

	return nil, errors.New("Invalid Token")
}
