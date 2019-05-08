package main

import (
	"github.com/jbogarin/go-cisco-webex-teams/sdk"
)

// Struct for Grafana Alert
type grafanaAlert struct {
	Title       string `json:"title"`
	RuleID      int    `json:"ruleId"`
	RuleName    string `json:"ruleName"`
	RuleURL     string `json:"ruleUrl"`
	State       string `json:"state"`
	ImageURL    string `json:"imageUrl"`
	Message     string `json:"message"`
	EvalMatches []struct {
		Metric string `json:"metric"`
		Tags   struct {
			Name string `json:"name"`
		} `json:"tags"`
		Value int `json:"value"`
	} `json:"evalMatches"`
}

type applicationConfig struct {
	webexClient *webexteams.Client
	webexToken  string
}

// An alert struct which we will operate on per request
type AlertEvent struct {
	GrafanaAlert  *grafanaAlert
	Template      string
	TargetAddress string
	TargetType    string
	IgnoreNoData  bool
	NoTags        bool
	NoImages      bool
}
