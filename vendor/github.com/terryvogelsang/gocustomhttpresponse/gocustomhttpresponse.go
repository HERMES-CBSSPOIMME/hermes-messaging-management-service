package gocustomhttpresponse

import (
	json "encoding/json"
	http "net/http"
	reflect "reflect"

	logruswrapper "github.com/terryvogelsang/logruswrapper"
)

var (
	defaultHeaders = map[string]string{
		"Content-Type": "application/json",
	}

	mapping = map[string]CustomHTTPResponse{

		logruswrapper.CodeSuccess: CustomHTTPResponse{
			StatusCode: http.StatusOK,
			Headers:    defaultHeaders,
		},

		logruswrapper.CodeBadLogin: CustomHTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Headers:    defaultHeaders,
		},

		logruswrapper.CodeInvalidToken: CustomHTTPResponse{
			StatusCode: http.StatusUnauthorized,
			Headers:    defaultHeaders,
		},
	}
)

// ProxyRequest : Received Request from Client
type CustomHTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       CustomHTTPResponseBody
}

type CustomHTTPResponseBody struct {
	LogInfos logruswrapper.LogEntryInfos `json:"logEntryInfos"`
	Content  interface{}                 `json:"content"`
}

func WriteResponse(content interface{}, logEntryInfos *logruswrapper.LogEntryInfos, w http.ResponseWriter) {

	var responseBody CustomHTTPResponseBody

	responseBody.LogInfos = *logEntryInfos
	responseBody.Content = content

	// Retrieve response settings
	settings := mapping[logEntryInfos.Code]

	// Write Headers
	for _, key := range reflect.ValueOf(settings.Headers).MapKeys() {

		w.Header().Set(key.String(), settings.Headers[key.String()])
	}

	// Write Status Code
	w.WriteHeader(settings.StatusCode)
	res, _ := json.Marshal(responseBody)
	w.Write(res)
}
