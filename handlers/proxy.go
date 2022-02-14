package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/btittelbach/go-bbhw"
)

type response struct {
	Function     string
	ResponseBody string
	HostName     string
}

func gpio_turn_on(pin_num uint) error {
	pin, err := bbhw.NewSysfsGPIO(pin_num, bbhw.OUT)
	err = pin.SetState(true)
	time.Sleep(500 * time.Millisecond)
	err = pin.SetState(false)
	time.Sleep(500 * time.Millisecond)
	err = pin.SetState(true)
	return err
}
// MakeProxy creates a proxy for HTTP web requests which can be routed to a function.
func MakeProxy() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		name := vars["name"]
		log.Info("proxy request: " + name)

		v, okay := functions[name]
		if !okay {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("{ \"status\" : \"Not found\"}"))
			log.Errorf("%s not found", name)
			return
		}



		//worker_list := [2]uint{66, 48}

		worker_list := map[int]uint{
			1: 66,
			2: 48,
		}

		gpio_turn_on(worker_list[1])

		v.InvocationCount = v.InvocationCount + 1

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)

		hostName, _ := os.Hostname()
		d := &response{
			Function:     name,
			ResponseBody: string(body),
			HostName:     hostName,
		}

		responseBody, err := json.Marshal(d)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			log.Errorf("error invoking %s. %v", name, err)
			return
		}

		w.Write(responseBody)

		log.Info("!!!!!proxy request: %s completed.", name)
	}
}
