package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type response struct {
	Function     string
	ResponseBody string
	HostName     string
}

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
			log.Errorf("%s not found", name)
			return
		}

		v.InvocationCount = v.InvocationCount + 1

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		hostName, _ := os.Hostname()
		d := &response{
			Function:     name,
			ResponseBody: string(body),
			HostName:     hostName,
		}

		responseBody, err := json.Marshal(d)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Errorf("error invoking %s. %v", name, err)
			return
		}

		w.Write(responseBody)

		log.Infof("!!!!!proxy request: %s completed.", name)
	}
}
