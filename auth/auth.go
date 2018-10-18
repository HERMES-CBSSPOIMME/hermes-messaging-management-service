package auth

import (

	// Native Go Libs
	json "encoding/json"
	errors "errors"
	fmt "fmt"
	http "net/http"

	bcrypt "golang.org/x/crypto/bcrypt"

	// Project Libs
	models "hermes-messaging-service/models"
	utils "hermes-messaging-service/utils"

	// 3rd Party Libs
	uuid "github.com/satori/go.uuid"
	logruswrapper "github.com/terryvogelsang/logruswrapper"
)

// CheckAuthentication : Return MQTT Auth Infos if provided auth token is valid,
// an error if present, a boolean flag indicating if user was already cached and a boolean flag indicating wether token was updated in Redis
func CheckAuthentication(env *models.Env, token string) (*models.MQTTAuthInfos, bool, bool, error) {

	// If no token, return an error
	if token == "" {
		return nil, false, false, errors.New("No Token Provided")
	}

	hashedToken, err := HashPassword(token)
	if err != nil {
		return nil, false, false, err
	}

	// Check if token is cached in Redis, Get UserID if it is
	cachedInternalUserID, _ := CheckIfTokenIsCached(env, token)

	if cachedInternalUserID != "" {

		// If yes : Return the cached infos
		return models.NewMQTTAuthInfos(cachedInternalUserID, hashedToken), true, false, nil

	} else {

		// If no : Verify with external endpoint
		err = env.RefreshConfig()

		if err != nil {
			return nil, false, false, err
		}

		MQTTAuthInfos, wasCached, wasTokenUpdated, err := VerifyTokenWithExternalEndpoint(env, token, hashedToken)

		if err != nil {
			return nil, false, false, err
		}

		return MQTTAuthInfos, wasCached, wasTokenUpdated, nil
	}
}

// HashPassword : Hash password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CheckIfTokenIsCached : Check if token is cached in Redis
func CheckIfTokenIsCached(env *models.Env, token string) (string, error) {

	// Check if token is cached in redis
	cachedInternalUserID, err := env.Redis.Get(fmt.Sprintf("session:%s", token))

	if err != nil {
		return "", err
	}

	if cachedInternalUserID != nil {
		return string(cachedInternalUserID), nil
	}

	return "", nil
}

// VerifyTokenWithExternalEndpoint : Verify token with provided external auth endpoint
func VerifyTokenWithExternalEndpoint(env *models.Env, token string, hashedToken string) (*models.MQTTAuthInfos, bool, bool, error) {

	// Create HTTP Client
	client := &http.Client{}

	// Init request
	req, err := http.NewRequest("GET", env.Config.AuthenticationCheckEndpoint, nil)

	if err != nil {
		return nil, false, false, err
	}

	// Add token header
	req.Header.Add("token", token)

	// Execute request
	res, err := client.Do(req)
	if err != nil {
		return nil, false, false, err
	}

	//=============================================================================
	// Authentication endpoint should return the following if token is invalid :
	//
	// HTTP Status Code : 400 (Bad Request)
	// Empty Body
	//=============================================================================
	// Authentication endpoint should return the following if token is valid :
	//
	// HTTP Status Code : 200 (OK)
	// Header(s) : content-type:application/json
	// {"UserID": "userID"}
	//=============================================================================
	if res.StatusCode == 200 {

		// Parse response to get Original User ID
		authCheckerBody := utils.AuthCheckerBody{}
		json.NewDecoder(res.Body).Decode(&authCheckerBody)

		// Check if user already has a cached token
		cachedInternalHermesUserID, cachedOldToken, _ := CheckIfUserAlreadyHasToken(env, authCheckerBody.OriginalUserID)

		if cachedOldToken != "" {

			// If yes : Update Redis with new token and revoke the older token
			UpdateRedisAndMongoDBWithNewToken(env, authCheckerBody.OriginalUserID, cachedInternalHermesUserID, cachedOldToken, token, hashedToken)

			// Return MQTTAuthInfos
			return models.NewMQTTAuthInfos(cachedInternalHermesUserID, hashedToken), true, true, nil

		} else {

			// If no : Create new token store + mapping in Redis
			newInternalHermesUserID := uuid.NewV4().String()

			// Store token in Redis
			env.Redis.Set(fmt.Sprintf("session:%s", token), []byte(newInternalHermesUserID))

			// Store mapping in Redis
			// TODO: Change Hset method to be able to set multiple field at a time
			env.Redis.HSet(fmt.Sprintf("mapping:%s", authCheckerBody.OriginalUserID), "token", []byte(token))
			env.Redis.HSet(fmt.Sprintf("mapping:%s", authCheckerBody.OriginalUserID), "internalHermesUserID", []byte(newInternalHermesUserID))

			// Return MQTTAuthInfos
			return models.NewMQTTAuthInfos(newInternalHermesUserID, hashedToken), false, false, nil
		}
	}

	return nil, false, false, errors.New(logruswrapper.CodeInvalidToken)
}

// CheckIfUserAlreadyHasToken : Check if originalUserID is already matched with one token in redis
func CheckIfUserAlreadyHasToken(env *models.Env, originalUserID string) (string, string, error) {

	cachedOldToken, err := env.Redis.HGet(fmt.Sprintf("mapping:%s", originalUserID), "token")
	cachedInternalUserID, err := env.Redis.HGet(fmt.Sprintf("mapping:%s", originalUserID), "internalHermesUserID")

	if err != nil {
		return "", "", err
	}

	return string(cachedInternalUserID), string(cachedOldToken), nil
}

// UpdateRedisAndMongoDBWithNewToken : Update old token store and mapping with new token
func UpdateRedisAndMongoDBWithNewToken(env *models.Env, originalUserID string, internalHermesUserID string, oldToken string, newToken string, newHashedToken string) error {

	// Update MongoDB Profile
	err := env.MongoDB.UpdatePassHash(internalHermesUserID, newHashedToken)

	if err != nil {
		return err
	}

	// Update Redis Token Store Key :
	// session:{oldToken} -> session:{newToken}
	err = env.Redis.Rename(fmt.Sprintf("session:%s", oldToken), fmt.Sprintf("session:%s", newToken))

	if err != nil {
		return err
	}

	// Update Redis Mapping Values :
	// mapping:{originalUserID} token {oldToken} ... --> mapping:{originalUserID} token {newToken} ...
	err = env.Redis.HSet(fmt.Sprintf("mapping:%s", originalUserID), "token", []byte(newToken))

	if err != nil {
		return err
	}

	return nil
}
