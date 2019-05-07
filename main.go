package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"time"

	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	"gopkg.in/resty.v1"

	"github.com/gorilla/mux"
)

var (
	roomID     string
	token      string
	timeout    int
	stdin      *os.File
	webexToken string
)

func init() {
	webexToken = os.Getenv("GRAFANA_BOT_TEAMS_TOKEN")
}

func sendMessage(alert *grafanaAlert, targetAddress string, targetType string) error {

	if err := validateEvent(alert); err != nil {
		return errors.New(err.Error())
	}

	client := resty.New()

	client.SetAuthToken(webexToken)
	Client := webexteams.NewClient(client)

	template := getTemplateNew(alert)

	var markDownMessage *webexteams.MessageCreateRequest
	switch targetType {
	// Email address
	case "emailAddress":
		markDownMessage = &webexteams.MessageCreateRequest{
			Markdown:      template,
			ToPersonEmail: targetAddress,
		}
		// Email from guid
	case "people":
		markDownMessage = &webexteams.MessageCreateRequest{
			Markdown:   template,
			ToPersonID: targetAddress,
		}
		// Room from guid
	case "room":
		markDownMessage = &webexteams.MessageCreateRequest{
			Markdown: template,
			RoomID:   targetAddress,
		}
	}

	if alert.ImageURL != "" {
		markDownMessage.Files = []string{alert.ImageURL}
	}

	newMarkDownMessage, _, err := Client.Messages.CreateMessage(markDownMessage)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("POST:", newMarkDownMessage.ID, newMarkDownMessage.Created)
	return nil
}

func handleWebhook(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	// This is the guid, which can be a roomID, personID, or even a personEmail
	targetData := vars["context"]

	targetType, err := decodeAlertTargetData(targetData)
	if err != nil {
		panic(nil)
	}

	var alert grafanaAlert
	if r.Method == http.MethodPost || r.Method == http.MethodPut {
		reqData, e := ioutil.ReadAll(r.Body)
		if e != nil {
			return
		}
		json.Unmarshal(reqData, &alert)
		sendMessage(&alert, targetData, targetType)
	}
	return
}

func main() {

	mux := mux.NewRouter()

	mux.HandleFunc("/webex/{context}", handleWebhook)

	muxWithMiddlewares := http.TimeoutHandler(mux, time.Second*10, "Timeout!")

	log.Fatal(http.ListenAndServe(":17575", muxWithMiddlewares))

}
