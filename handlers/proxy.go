package handlers

import (
	"encoding/json"
	// b64 "encoding/base64"
	// "fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/btittelbach/go-bbhw"
	"bytes"
	// "strconv"

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
type Payload struct{
	Fid string `json:"fid"`
	Src string `json:"src"`
	Params string `json:"params,omitempty"`
	Lang string `json:"lang"`
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

		// Working GPIO pins
		worker_list := map[int]uint{
			1: 48, // works
			2: 67, // works
			3: 68, // works
		}

		gpio_turn_on(worker_list[3])

		v.InvocationCount = v.InvocationCount + 1

		defer r.Body.Close()
		body, _ := ioutil.ReadAll(r.Body)
		// // body_str := string(body)
		// // log.Info(body_str)/
		// log.Info(string(body))
		// var payload Payload
		// json.Unmarshal([]byte(body), &payload)
		// // log.Info(payload.Src)
		// data, _ := b64.StdEncoding.DecodeString(payload.Src)
		// log.Info(data)
		// s2, _ := strconv.Unquote(string(data))
		// log.Info(s2)


		resp, err := http.Post("http://128.197.176.240:8080/run", "application/json",
			bytes.NewBuffer(body))

		if err != nil {
			log.Fatal(err)
		}
		resp_body, _ := ioutil.ReadAll(resp.Body)
		log.Info(string(resp_body))


		hostName, _ := os.Hostname()
		d := &response{
			Function:     name,
			ResponseBody: "'" + string(body) + "'",
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
