package handlers

import (
	"context"

	"github.com/openfaas/faas-provider/logs"
)

// LogRequester implements the Requester interface
type LogRequester struct{}

// NewLogRequester returns a Requestor instance that can be used in the function logs endpoint
func NewLogRequester() logs.Requester {
	return &LogRequester{}
}

// Query implements the actual Swarm logs request logic for the Requester interface
func (l LogRequester) Query(ctx context.Context, r logs.Request) (<-chan logs.Message, error) {
	msgStream := make(chan logs.Message, 1)
	msgStream <- logs.Message{
		Name:      r.Name,
		Namespace: r.Namespace,
		Text:      "memory log line",
	}
	close(msgStream)
	return msgStream, nil
}
