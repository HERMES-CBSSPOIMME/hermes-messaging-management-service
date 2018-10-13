package auth

import (
	// Native Go Libs
	json "encoding/json"
	errors "errors"
	http "net/http"

	// Project Libs
	models "hermes-messaging-service/models"
)

// CheckAuthentication : Check authentication by executing a GET HTTP Request with token as "token" header value
// Return MQTT Auth Infos if provided auth token is valid, an error otherwise
func CheckAuthentication(env *models.Env, token string) (*models.MQTTAuthInfos, error) {

	// If no token, return an error
	if token == "" {
		return nil, errors.New("No Token Provided")
	}

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
	// Body : {"clientID": "userID"}
	if res.StatusCode == 200 {

		MQTTAuthInfos := models.MQTTAuthInfos{}

		// Set ClientID from response
		err := json.NewDecoder(res.Body).Decode(&MQTTAuthInfos)

		if err != nil {
			return nil, err
		}

		// Copy ClientID as Username
		MQTTAuthInfos.Username = MQTTAuthInfos.ClientID

		// Set token as password
		MQTTAuthInfos.Password = token

		return &MQTTAuthInfos, nil
	}

	return nil, errors.New("Invalid Token")
}
