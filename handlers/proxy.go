package handlers

import (
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// MakeProxy creates a proxy for HTTP web requests which can be routed to a function.
func MakeProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		log.Info("proxy request: " + name)

		v, okay := functions[name]
		if !okay {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{ \"status\" : \"Not found\"}"))
			return
		}

		v.InvocationCount = v.InvocationCount + 1
		responseBody := "{ \"status\" : \"Okay\"}"
		w.Write([]byte(responseBody))
	}
}