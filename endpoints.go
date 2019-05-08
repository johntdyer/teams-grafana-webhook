package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// HealthCheckHandler handler
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	// A very simple health check.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

// GrafanaAlertHandler handler
func GrafanaAlertHandler(w http.ResponseWriter, r *http.Request) {

	event := &AlertEvent{
		NoTags:       false,
		IgnoreNoData: false,
		NoImages:     false,
	}

	vars := mux.Vars(r)

	// This is the guid, which can be a roomID, personID, or even a personEmail
	targetAddressRaw := vars["context"]

	if r.URL.Query().Get("noTags") != "" {
		event.NoTags = true
		log.WithFields(logrus.Fields{
			"targetAddressRaw": targetAddressRaw}).Debug("Detected Brief / no tags mode")
	}

	if r.URL.Query().Get("ignoreNoData") != "" {
		event.IgnoreNoData = true
		log.WithFields(logrus.Fields{
			"targetAddressRaw": targetAddressRaw}).Debug("Detected IgnoreNoData mode, will ignore alerts w/ action type 'no_data'")
	}

	if r.URL.Query().Get("noImages") != "" {
		event.NoImages = true
		log.WithFields(logrus.Fields{
			"targetAddressRaw": targetAddressRaw}).Debug("Detected NoImage mode")

	}
	targetType, err := decodeAlertTargetData(targetAddressRaw)
	if err != nil {

		log.Fatal(err)
	}

	log.WithFields(logrus.Fields{
		"targetAddressRaw": targetAddressRaw, "targetType": targetType}).Debug("target")

	event.TargetAddress = targetAddressRaw
	event.TargetType = targetType

	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		reqData, e := ioutil.ReadAll(r.Body)
		if e != nil {
			return
		}
		json.Unmarshal(reqData, &event.GrafanaAlert)

		err := event.postMessage()
		if err != nil {
			log.Fatal(err)
		}

	}
	return
}
