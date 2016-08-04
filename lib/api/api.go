package api

import (
	"log"
	"net/http"
)

type api interface {
	Path() string
	GET(http.ResponseWriter, *http.Request) error
	PUT(http.ResponseWriter, *http.Request) error
	POST(http.ResponseWriter, *http.Request) error
	DELETE(http.ResponseWriter, *http.Request) error
}

// Implements the default handlers for the API.
type ApiDefaults struct{}

func (a *ApiDefaults) GET(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}
func (a *ApiDefaults) PUT(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}
func (a *ApiDefaults) POST(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}
func (a *ApiDefaults) DELETE(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusMethodNotAllowed)
	return nil
}

func API(a api, logger *log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var err error
		switch req.Method {
		case "GET":
			err = a.GET(w, req)
		case "PUT":
			err = a.PUT(w, req)
		case "POST":
			err = a.POST(w, req)
		case "DELETE":
			err = a.DELETE(w, req)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
		if err != nil {
			logger.Println(err)
		}
	}
}
