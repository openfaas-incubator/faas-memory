package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

// MakeProxy creates a proxy for HTTP web requests which can be routed to a function.
func MakeProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]

		v, okay := functions[name]
		if !okay {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{ \"status\" : \"Not found\"}"))
		}

		v.InvocationCount = v.InvocationCount + 1
		responseBody := "{ \"status\" : \"Okay\"}"
		w.Write([]byte(responseBody))
	}
}