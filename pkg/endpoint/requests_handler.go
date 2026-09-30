package endpoint

import (
	"fmt"
	"net/http"

	"log"
)

// custom type for handlers to return their error
type appError struct {
	Error   error
	Message string
	Code    int
}

// custom wrapper for http handler
type appHandler func(http.ResponseWriter, *http.Request) *appError

// format for error response page with message string
const httpErrorResponsePage string = `
ERROR %d
	%s%s`

func (serveHttp appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	handledRequest := serveHttp(w, r)

	if handledRequest != nil {
		errorString := ""
		if handledRequest.Error != nil {
			errorString = ": " + handledRequest.Error.Error()
		}

		log.Printf("<- ERROR %d | %s %s | %s%s | User-Agent: %s",
			handledRequest.Code,
			r.Method,
			r.RequestURI,
			handledRequest.Message,
			errorString,
			r.UserAgent())

		responseErrorPage := fmt.Sprintf(httpErrorResponsePage,
			handledRequest.Code,
			handledRequest.Message,
			errorString)

		http.Error(w, responseErrorPage, handledRequest.Code)
	}
}

func logRequest(r *http.Request) {
	// Extract the username from basic authentication
	username := ""
	username, _, ok := r.BasicAuth()
	if ok {
		username = " | User: " + username
	}

	log.Printf(
		"-> %s %s%s | %s | User-Agent: %s",
		r.Method,
		r.RequestURI,
		username,
		r.Proto,
		r.UserAgent(),
	)
}

// custom wrapper to capture the status code for logging
type responseWriter struct {
	http.ResponseWriter
	Code int
}

// WriteHeader captures the status code to log later
func (rw *responseWriter) WriteHeader(code int) {
	rw.Code = code
	rw.ResponseWriter.WriteHeader(code)
}

// Logger middleware to log each request and capture errors
func loggedHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap the original ResponseWriter
		rw := &responseWriter{ResponseWriter: w, Code: http.StatusOK} // Default to 200

		logRequest(r)

		handler.ServeHTTP(rw, r)

		// Log the status code after serving the request
		if rw.Code >= 400 {
			log.Printf("<- ERROR %d | %s %s | User-Agent: %s",
				rw.Code,
				r.Method,
				r.RequestURI,
				r.UserAgent())
		}
	})
}
