package main

import (
	"net/http"
	"os"

	"time"

	runtime "github.com/banzaicloud/logrus-runtime-formatter"
	"github.com/gorilla/mux"
	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
	"github.com/sirupsen/logrus"
	resty "gopkg.in/resty.v1"
)

var (
	roomID     string
	token      string
	timeout    int
	stdin      *os.File
	log        = logrus.New()
	appConfig  = &applicationConfig{}
	listenPort = "17575"
)

func init() {

	childFormatter := logrus.TextFormatter{}
	runtimeFormatter := &runtime.Formatter{ChildFormatter: &childFormatter}
	log.Formatter = runtimeFormatter

	if os.Getenv("DEBUG") == "" {
		log.Level = logrus.InfoLevel

	} else {
		log.Level = logrus.DebugLevel
	}

	if os.Getenv("GRAFANA_BOT_TEAMS_TOKEN") == "" {
		log.Fatal("GRAFANA_BOT_TEAMS_TOKEN environment variable not found, aborting")
	}

	// Create
	client := resty.New()
	client.SetAuthToken(os.Getenv("GRAFANA_BOT_TEAMS_TOKEN"))
	appConfig.webexClient = webexteams.NewClient(client)
}

func main() {

	mux := mux.NewRouter()

	mux.HandleFunc("/webex/{context}", GrafanaAlertHandler)
	mux.HandleFunc("/health", HealthCheckHandler)

	muxWithMiddlewares := http.TimeoutHandler(mux, time.Second*10, "Timeout!")
	log.Infof("Ready, Listening on :%s\n", listenPort)
	log.Fatal(http.ListenAndServe(":"+listenPort, muxWithMiddlewares))

}
