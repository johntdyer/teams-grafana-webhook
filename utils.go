package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"html/template"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"

	"github.com/Masterminds/sprig"
	webexteams "github.com/jbogarin/go-cisco-webex-teams/sdk"
)

const myRegex = `^(?:\w*)://(?:us)/(?P<CONTEXT>[A-Z]*)/(?P<GUID>(?:\w|-){36})`

const emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

func (e *AlertEvent) createMessageRequest() (*webexteams.MessageCreateRequest, error) {

	switch e.TargetType {
	// Email address
	case "emailAddress":
		return &webexteams.MessageCreateRequest{
			Markdown:      e.Template,
			ToPersonEmail: e.TargetAddress,
		}, nil
		// Email from guid
	case "people":
		return &webexteams.MessageCreateRequest{
			Markdown:   e.Template,
			ToPersonID: e.TargetAddress,
		}, nil
		// Room from guid
	case "room":
		return &webexteams.MessageCreateRequest{
			Markdown: e.Template,
			RoomID:   e.TargetAddress,
		}, nil

	}
	return &webexteams.MessageCreateRequest{}, errors.New("unable to render message type")
}

func (e *AlertEvent) postMessage() error {

	log.Debug("Starting message post")

	if err := validateEvent(e.GrafanaAlert); err != nil {
		return errors.New(err.Error())
	}

	// markDownMessage := &webexteams.MessageCreateRequest{}

	if e.GrafanaAlert.State == "no_data" && e.IgnoreNoData == true {
		log.WithFields(logrus.Fields{
			"AlertTitle":    e.GrafanaAlert.Title,
			"AlertRuleName": e.GrafanaAlert.RuleName,
			"AlertRuleUrl":  e.GrafanaAlert.RuleURL,
			"AlertRuleID":   e.GrafanaAlert.RuleID,
		}).Debug("no data event, ignoring")
		return nil
	}

	// Render template
	e.renderTemplate()

	markDownMessage, err := e.createMessageRequest()
	if err != nil {
		log.Fatal(err)
	}

	// Ignore images if present in noImage
	if e.NoTags == false {

		if e.GrafanaAlert.ImageURL != "" {
			log.WithFields(logrus.Fields{
				"briefMode": false,
			}).Debug("Adding Image")

			markDownMessage.Files = []string{e.GrafanaAlert.ImageURL}
		}
	} else {
		log.WithFields(logrus.Fields{
			"briefMode": true,
		}).Debug("Brief Mode")
	}

	newMarkDownMessage, _, err := appConfig.webexClient.Messages.CreateMessage(markDownMessage)
	if err != nil {
		log.Fatal(err)
	}

	log.WithFields(logrus.Fields{
		"messageID": newMarkDownMessage.ID,
	}).Debug("Post")
	return nil

}

// Define a template.
const inccidentTemplate string = `
<blockquote class='{{.EventColor}}'> {{.Emoji}} {{.MessageStatus}} <br/>
<b>Check Name:</b> {{.Title}}<br/>
<b>Rule:</b> <a href="{{.RuleURL}}">{{.RuleName}}</a><br/>
{{if (ne .MessageStatus "Resolved") }}
<b>Message:</b> {{.Message}} <br/>
{{end}}
{{ with .EvalMatches }}
{{ range . }}
<li>{{ .Metric }}</li> -> {{ .Value }}</li>
-->{{list .Tags | join "," }}
{{ end }}
{{ end }}

</blockquote>
`

func (e *AlertEvent) renderTemplate() {

	log.WithFields(logrus.Fields{
		"AlertTitle":    e.GrafanaAlert.Title,
		"AlertRuleName": e.GrafanaAlert.RuleName,
		"AlertRuleUrl":  e.GrafanaAlert.RuleURL,
		"AlertRuleID":   e.GrafanaAlert.RuleID,
	}).Debug("rendering template")

	e.minimizeTemplate(inccidentTemplate)
	// templateWithoutNewlines := stringMinifier(inccidentTemplate)

	eventColor, emoji, formatedStatus := stateToEmojifier(e.GrafanaAlert)

	t := template.Must(template.New("inccident").Funcs(sprig.FuncMap()).Parse(e.Template))

	var tpl bytes.Buffer

	localStruct := struct {
		Title         string
		RuleID        int
		RuleName      string
		RuleURL       string
		State         string
		ImageURL      string
		Message       string
		EventColor    string
		Emoji         string
		MessageStatus string
		EvalMatches   []struct {
			Metric string `json:"metric"`
			Tags   struct {
				Name string `json:"name"`
			} `json:"tags"`
			Value int `json:"value"`
		} `json:"evalMatches"`
	}{
		e.GrafanaAlert.Title,
		e.GrafanaAlert.RuleID,
		e.GrafanaAlert.RuleName,
		e.GrafanaAlert.RuleURL,
		e.GrafanaAlert.State,
		e.GrafanaAlert.ImageURL,
		e.GrafanaAlert.Message,
		eventColor,
		emoji,
		formatedStatus,
		e.GrafanaAlert.EvalMatches,
	}

	err := t.Execute(&tpl, localStruct)
	if err != nil {
		log.Fatal(err)
		// panic(err)
	}
	e.Template = tpl.String()
	//return tpl.String()

}

func (e *AlertEvent) minimizeTemplate(in string) {
	// var out string
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				e.Template = e.Template + " "
			}
			white = true
		} else {
			e.Template = e.Template + string(c)
			white = false
		}
	}

}

// Get tje correct unicode emoji for alert
func stateToEmojifier(event *grafanaAlert) (string, string, string) {

	log.WithFields(logrus.Fields{
		"eventState": event.State,
	}).Debug("state to emoji method")

	switch event.State {

	case "ok":
		return "success", "✅", "Success"
	case "alerting":
		return "danger", "🚨", "Critical"
	case "paused":
		return "warning", "️⚠️", "Warning"
	case "pending":
		return "secondary", "️⚠️", "Pending"
	case "no_data":
		return "info", "⁉️", "No Data"
	default:
		return "primary", "Unknown", "Unknown"
	}
}

func parseTime(input time.Time) string {
	return input.Format("Monday 01/02/2006 - 15:04:05 MST")
}

// func formattedEventAction(event *types.Event) string {
// 	switch event.Check.Status {
// 	case 0:
// 		return "RESOLVED"
// 	default:
// 		return "ALERT"
// 	}
// }

func validateEvent(alert *grafanaAlert) error {

	log.WithFields(logrus.Fields{
		"AlertTitle":    alert.Title,
		"AlertRuleName": alert.RuleName,
		"AlertRuleUrl":  alert.RuleURL,
	}).Debug("Starting validation of even")

	if alert.Title == "" {
		return errors.New("tile is missing")
	}
	if alert.RuleName == "" {
		return errors.New("Rule name is missing from event")
	}
	if alert.RuleURL == "" {
		return errors.New("Rule URL is missing from event")
	}

	return nil
}

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)
	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}

// decodeAlertTargetData Used to determine the type of message we're sending based on the path. Base64 encoded personID, SpaceID, or a valid email address
func decodeAlertTargetData(targetData string) (string, error) {
	log.WithFields(logrus.Fields{
		"targetAddress": targetData,
	}).Debug("decode target address")

	// Check if its an email
	re := regexp.MustCompile(emailRegex)
	if re.MatchString(targetData) {
		return "emailAddress", nil
	}

	// Otherwise a guid
	base64Text := make([]byte, base64.RawStdEncoding.DecodedLen(len(targetData)))

	n, err := base64.RawStdEncoding.Decode(base64Text, []byte(targetData))
	if err != nil {
		panic(err)
	}

	r := regexp.MustCompile(myRegex)

	f := findNamedMatches(r, string(base64Text[:n]))
	if f["CONTEXT"] == "" {
		return "", errors.New("Unknown context")

	}
	return strings.ToLower(f["CONTEXT"]), nil

}
