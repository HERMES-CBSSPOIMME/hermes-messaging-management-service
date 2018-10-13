package router

import (
	"fmt"

	// Native Go Libs
	"hermes-messaging-service/models"
	http "net/http"

	// 3rd Party Libs
	gocustomhttpresponse "github.com/terryvogelsang/gocustomhttpresponse"
	logruswrapper "github.com/terryvogelsang/logruswrapper"
)

type (
	Handler func(env *models.Env, w http.ResponseWriter, r *http.Request) error
)

type Greeter struct {
	Message string
}

// HelloWorld : A Simple HelloWorld Endpoint
func HelloWorld(env *models.Env, w http.ResponseWriter, r *http.Request) error {

	// Logging demo
	log := logruswrapper.NewEntry("UsersService", "/helloworld", logruswrapper.CodeSuccess)
	logruswrapper.Info(log)

	err := env.DB.AddUserACL("testclientid", "testusername", "testpassword")

	if err != nil {
		fmt.Println(err)
		return err
	}

	greeter := Greeter{Message: "Hello World"}

	gocustomhttpresponse.WriteResponse(greeter, log, w)
	return nil
}

// CustomHandle : Custom Handlers Wrapper for API
func CustomHandle(env *models.Env, handlers ...Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range handlers {
			err := h(env, w, r)
			if err != nil {
				// w.Write(getResponseOfError(err))
				return
			}
		}
	})
}
