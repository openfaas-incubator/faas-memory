package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/openfaas/faas/gateway/requests"
)

// MakeUpdateHandler update specified function
func MakeUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("update request")

		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.CreateFunctionRequest{}
		if err := json.Unmarshal(body, &request); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, okay := functions[request.Service]
		if !okay {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{ \"status\" : \"Not found\"}"))
			return
		}

		functions[request.Service] = createToRequest(request)
	}
}
