package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	typesv1 "github.com/openfaas/faas-provider/types"
	log "github.com/sirupsen/logrus"
)

// MakeUpdateHandler update specified function
func MakeUpdateHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("update request")

		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := typesv1.FunctionDeployment{}
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

		functions[request.Service] = requestToStatus(request)
	}
}
