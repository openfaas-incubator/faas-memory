// Copyright (c) Edward Wilde 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/openfaas/faas/gateway/requests"
	log "github.com/sirupsen/logrus"
)

var functions = map[string]*requests.Function{}

// MakeDeployHandler creates a handler to create new functions in the cluster
func MakeDeployHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info("deployment request")
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)

		request := requests.CreateFunctionRequest{}
		if err := json.Unmarshal(body, &request); err != nil {
			log.Errorln("error during unmarshal of create function request. ", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		functions[request.Service] = createToRequest(request)

		log.Infof("deployment request for function %s", request.Service)

		w.WriteHeader(http.StatusOK)
	}
}
