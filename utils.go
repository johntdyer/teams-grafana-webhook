package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"regexp"
	"strings"
	"time"
	"unicode"

	"html/template"

	"github.com/Masterminds/sprig"
)

//#.com/alecthomas/template"
const myRegex = `^(?:\w*)://(?:us)/(?P<CONTEXT>[A-Z]*)/(?P<GUID>(?:\w|-){36})`

const emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"

// stringMinifier remove whitespace before sending message to teams
func stringMinifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}

func stateToEmojifier(event *grafanaAlert) (string, string, string) {
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

// Define a template.
const inccidentTemplate = `
<blockquote class='{{.EventColor}}'> {{.Emoji}} {{.MessageStatus}} <br/>
<b>Check Name:</b> {{.Title}}
<b>Rule:</b> <a href="{{.RuleURL}}">{{.RuleName}}</a><br/>
{{if (ne .MessageStatus "Resolved") }}
<b>Message:</b> {{.Message}} <br/>
{{end}}
<br/>
{{ with .EvalMatches }}
{{ range . }}

<li>{{ .Metric }}</li> -> {{ .Value }}</li>

	--> {{list .Tags | join "," }}
{{ end }}
{{ end }}

</blockquote>





`

func getTemplateNew(ga *grafanaAlert) string {

	templateWithoutNewlines := stringMinifier(inccidentTemplate)

	eventColor, emoji, formatedStatus := stateToEmojifier(ga)
	//.Funcs(template.FuncMap{"parseTime": parseTime,})
	t := template.Must(template.New("inccident").Funcs(sprig.FuncMap()).Parse(templateWithoutNewlines))

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
		ga.Title,
		ga.RuleID,
		ga.RuleName,
		ga.RuleURL,
		ga.State,
		ga.ImageURL,
		ga.Message,
		eventColor,
		emoji,
		formatedStatus,
		ga.EvalMatches,
	}

	err := t.Execute(&tpl, localStruct)
	if err != nil {
		panic(err)
	}

	return tpl.String()

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

func decodeAlertTargetData(targetData string) (string, error) {

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
