package endpoint

import (
	"fmt"
	"net/http"
)

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

const httpErrorResponsePage string = `
ERROR %d
	%s: %s`

func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ae := fn(w, r)
	if ae != nil {
		message := fmt.Sprintf(httpErrorResponsePage, ae.Code, ae.Message, ae.Error.Error())
		http.Error(w, message, ae.Code)
	}
}
