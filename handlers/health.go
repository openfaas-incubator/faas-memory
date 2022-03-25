package handlers

import (
	"net/http"
	"net"

	log "github.com/sirupsen/logrus"
)

// MakeHealthHandler returns 200/OK when healthy
func MakeHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip, _, err:=net.SplitHostPort(r.RemoteAddr)
		if(err != nil){
			log.Info("error during health check");
		}
		userIP := net.ParseIP(ip)
		defer r.Body.Close()
		log.Info("health check request")
		log.Info("healthy from " + (userIP.String()))
		for _, worker := range allWorkers{
			if (worker.ip == userIP.String()){
				worker.status = READY
				log.Info(worker.ip+"is READY")
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
