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
	log.Printf("Serving %s %s", r.Method, r.RequestURI)
	ae := fn(w, r)
	if ae != nil {
		errorString := ""
		if ae.Error != nil {
			errorString = ": " + ae.Error.Error()
		}

		log.Printf("ERROR %d - Request %s %s - %s%s", ae.Code, r.Method, r.RequestURI, ae.Message, errorString)

		message := fmt.Sprintf(httpErrorResponsePage, ae.Code, ae.Message, errorString)
		http.Error(w, message, ae.Code)
	}
}
