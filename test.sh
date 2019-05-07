curl -X POST -d '{
  "title": "My alert",
  "ruleId": 1,
  "ruleName": "Load peaking",
  "ruleUrl": "http://testing.com/db/dashboard/my_dashboard?panelId=2",
  "state": "alerting",
  "imageUrl": "http://docs.grafana.org/img/docs/v43/heatmap_histogram.png",
  "message": "Load is peaking. Make sure the traffic is real and spin up more webfronts",
  "evalMatches": [
    {
      "metric": "requests",
      "tags": { "name": "fireplace_chimney" },
      "value": 122
    }
  ]
}' http://localhost:17575/webex/johndye@cisco.com

#Y2lzY29zcGFyazovL3VzL1BFT1BMRS8xOWE5YzhhNS1lZmRjLTRjNTgtODM4Yy1kNzU2OTk0YjQwN2E
#