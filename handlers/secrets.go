package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	typesv1 "github.com/openfaas/faas-provider/types"
)

var secrets = map[string]typesv1.Secret{}

func MakeSecretsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			defer r.Body.Close()
		}

		body, readBodyErr := ioutil.ReadAll(r.Body)
		if readBodyErr != nil {
			log.Printf("couldn't read body of a request: %s", readBodyErr)

			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		var (
			responseStatus int
			responseBody   []byte
			responseErr    error
		)

		switch r.Method {
		case http.MethodGet:
			responseStatus, responseBody, responseErr = getSecrets()
		case http.MethodPost:
			responseStatus, responseBody, responseErr = createNewSecret(body)
		case http.MethodPut:
			responseStatus, responseBody, responseErr = updateSecret(body)
		case http.MethodDelete:
			responseStatus, responseBody, responseErr = deleteSecret(body)
		}

		if responseErr != nil {
			log.Println(responseErr)
			w.WriteHeader(responseStatus)
			return
		}

		w.WriteHeader(responseStatus)
		if responseBody != nil {
			_, writeErr := w.Write(responseBody)

			if writeErr != nil {
				log.Println("cannot write body of a response")
				return
			}
		}
	}
}

func getSecrets() (responseStatus int, responseBody []byte, err error) {
	results := []typesv1.Secret{}

	for _, s := range secrets {
		results = append(results, s)
	}

	resultsJson, marshalErr := json.Marshal(results)
	if marshalErr != nil {
		return http.StatusInternalServerError,
			nil,
			fmt.Errorf("error marshalling secrets to json: %s", marshalErr)

	}

	return http.StatusOK, resultsJson, nil
}

func createNewSecret(body []byte) (responseStatus int, responseBody []byte, err error) {
	var secret typesv1.Secret

	unmarshalErr := json.Unmarshal(body, &secret)
	if unmarshalErr != nil {
		return http.StatusBadRequest, nil, fmt.Errorf(
			"error unmarshalling body to json in secretPostHandler: %s",
			unmarshalErr,
		)
	}

	secrets[secret.Name] = secret
	return http.StatusCreated, nil, nil
}

func updateSecret(body []byte) (responseStatus int, responseBody []byte, err error) {
	var secret typesv1.Secret

	unmarshalErr := json.Unmarshal(body, &secret)
	if unmarshalErr != nil {
		return http.StatusBadRequest, nil, fmt.Errorf(
			"error unmarshalling body to json in secretPostHandler: %s",
			unmarshalErr,
		)
	}

	secrets[secret.Name] = secret
	return http.StatusCreated, nil, nil
}

func deleteSecret(body []byte) (responseStatus int, responseBody []byte, err error) {
	var secret typesv1.Secret

	unmarshalErr := json.Unmarshal(body, &secret)
	if unmarshalErr != nil {
		return http.StatusBadRequest, nil, fmt.Errorf(
			"error unmarshaling secret in secretDeleteHandler: %s",
			unmarshalErr,
		)
	}

	delete(secrets, secret.Name)

	return http.StatusOK, nil, nil
}
