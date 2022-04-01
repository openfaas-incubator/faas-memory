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
		for i := range allWorkers{
			if (allWorkers[i].ip == userIP.String()){
				allWorkers[i].status = READY
				log.Info(allWorkers[i].ip+" is READY")
				log.Info(allWorkers[0].status, allWorkers[1].status, allWorkers[2].status)
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
