package endpoint

import (
	"fmt"
	"net/http"

	"log"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

const httpErrorResponsePage string = `
ERROR %d
	%s%s`

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract the username from basic authentication
	username := ""
	username, _, ok := r.BasicAuth()
	if ok {
		username = " | User: " + username
	}

	log.Printf(
		"%s %s%s | %s | User-Agent: %s",
		r.Method,      // HTTP method (GET, POST, etc.)
		r.RequestURI,  // Full request URI
		username,      // Client's user agent string
		r.Proto,       // Protocol (HTTP/1.1, etc.)
		r.UserAgent(), // Username from basic authentication or "anonymous"
	)

	ae := fn(w, r)

	if ae != nil {
		errorString := ""
		if ae.Error != nil {
			errorString = ": " + ae.Error.Error()
		}

		log.Printf("ERROR %d | %s %s | %s%s | User-Agent: %s",
			ae.Code,
			r.Method,
			r.RequestURI,
			ae.Message,
			errorString,
			r.UserAgent())

		message := fmt.Sprintf(httpErrorResponsePage, ae.Code, ae.Message, errorString)
		http.Error(w, message, ae.Code)
	}
}
