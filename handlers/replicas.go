// Copyright (c) Edward Wilde 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/openfaas/faas/gateway/requests"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// MakeReplicaUpdater updates desired count of replicas
func MakeReplicaUpdater() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("update replicas, nothing to do here")

	}
}

// MakeReplicaReader reads the amount of replicas for a deployment
func MakeReplicaReader() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Info("read replicas")

		vars := mux.Vars(r)
		functionName := vars["name"]

		found := &requests.Function{}
		found.Name = functionName
		found.AvailableReplicas = 1

		functionBytes, _ := json.Marshal(found)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(functionBytes)
	}
}
