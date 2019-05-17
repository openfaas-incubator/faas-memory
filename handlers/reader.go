// Copyright (c) Edward Wilde 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/openfaas/faas/gateway/requests"
	log "github.com/sirupsen/logrus"
)

// MakeFunctionReader handler for reading functions deployed in the cluster as deployments.
func MakeFunctionReader() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		log.Info("read request")
		functions, err := readServices()
		if err != nil {
			log.Printf("Error getting service list: %s\n", err.Error())

			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		functionBytes, _ := json.Marshal(functions)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(functionBytes)
	}
}

func readServices() ([]*requests.Function, error) {
	var list []*requests.Function
	for _, v := range functions {
		list = append(list, v)
	}

	return list, nil
}

func createToRequest(request requests.CreateFunctionRequest) *requests.Function {
	return &requests.Function{
		Name:              request.Service,
		Annotations:       request.Annotations,
		EnvProcess:        request.EnvProcess,
		Image:             request.Image,
		Labels:            request.Labels,
		AvailableReplicas: 1,
		Replicas:          1,
	}
}
