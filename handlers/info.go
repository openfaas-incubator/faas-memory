package handlers

import (
	"encoding/json"
	"net/http"

	typesv1 "github.com/openfaas/faas-provider/types"
	log "github.com/sirupsen/logrus"
)

const (
	//OrchestrationIdentifier identifier string for provider orchestration
	OrchestrationIdentifier = "memory"
	//ProviderName name of the provider
	ProviderName = "faas-memory"
)

//MakeInfoHandler creates handler for /system/info endpoint
func MakeInfoHandler(version, sha string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}

		log.Info("info request")
		
		temp := typesv1.VersionInfo{
				Release: version,
				SHA:     sha,
			}
		
		infoResponse := typesv1.ProviderInfo{
			Orchestration: OrchestrationIdentifier,
			Name:      ProviderName,
			Version: &temp,
		}

		jsonOut, marshalErr := json.Marshal(infoResponse)
		if marshalErr != nil {
			log.Error("Error during unmarshal of info request ", marshalErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonOut)
	}
}
