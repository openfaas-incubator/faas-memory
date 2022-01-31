package main

import (
	"os"
	"strings"

	"github.com/yannip1234/faas-memory/handlers"
	"github.com/yannip1234/faas-memory/types"
	"github.com/yannip1234/faas-memory/version"
	bootstrap "github.com/openfaas/faas-provider"

	"github.com/openfaas/faas-provider/logs"
	bootTypes "github.com/openfaas/faas-provider/types"
	log "github.com/sirupsen/logrus"
)

func init() {
	logFormat := os.Getenv("LOG_FORMAT")
	logLevel := os.Getenv("LOG_LEVEL")
	if strings.EqualFold(logFormat, "json") {
		log.SetFormatter(&log.JSONFormatter{
			FieldMap: log.FieldMap{
				log.FieldKeyMsg:  "message",
				log.FieldKeyTime: "@timestamp",
			},
			TimestampFormat: "2006-01-02T15:04:05.999Z07:00",
		})
	} else {
		log.SetFormatter(&log.TextFormatter{
			FullTimestamp: true,
		})
	}

	if level, err := log.ParseLevel(logLevel); err == nil {
		log.SetLevel(level)
	}
}

func main() {

	log.Infof("faas-memory version:%s. Last commit message: %s, commit SHA: %s'", version.BuildVersion(), version.GitCommitMessage, version.GitCommitSHA)

	readConfig := types.ReadConfig{}
	osEnv := types.OsEnv{}
	cfg := readConfig.Read(osEnv)

	bootstrapHandlers := bootTypes.FaaSHandlers{
		FunctionProxy:        handlers.MakeProxy(),
		DeleteHandler:        handlers.MakeDeleteHandler(),
		DeployHandler:        handlers.MakeDeployHandler(),
		FunctionReader:       handlers.MakeFunctionReader(),
		ReplicaReader:        handlers.MakeReplicaReader(),
		ReplicaUpdater:       handlers.MakeReplicaUpdater(),
		UpdateHandler:        handlers.MakeUpdateHandler(),
		HealthHandler:        handlers.MakeHealthHandler(),
		InfoHandler:          handlers.MakeInfoHandler(version.BuildVersion(), version.GitCommitSHA),
		SecretHandler:        handlers.MakeSecretsHandler(),
		LogHandler:           logs.NewLogHandlerFunc(handlers.NewLogRequester(), cfg.WriteTimeout),
		ListNamespaceHandler: handlers.NamespaceLister(),
	}

	bootstrapConfig := bootTypes.FaaSConfig{
		ReadTimeout:     cfg.ReadTimeout,
		WriteTimeout:    cfg.WriteTimeout,
		TCPPort:         &cfg.Port,
		EnableHealth:    true,
		EnableBasicAuth: false,
	}

	log.Infof("listening on port %d ...", cfg.Port)
	bootstrap.Serve(&bootstrapHandlers, &bootstrapConfig)
}
