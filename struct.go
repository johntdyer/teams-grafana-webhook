package main

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
