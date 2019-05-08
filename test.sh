# curl -X POST -d '{
#   "title": "My alert",
#   "ruleId": 1,
#   "ruleName": "Load peaking",
#   "ruleUrl": "http://testing.com/db/dashboard/my_dashboard?panelId=2",
#   "state": "alerting",
#   "imageUrl": "http://docs.grafana.org/img/docs/v43/heatmap_histogram.png",
#   "message": "Load is peaking. Make sure the traffic is real and spin up more webfronts",
#   "evalMatches": [
#     {
#       "metric": "requests",
#       "tags": { "name": "fireplace_chimney" },
#       "value": 122
#     }
#   ]
# }' http://localhost:17575/webex/johndye@cisco.com

curl -X POST -d '{
    "evalMatches": [],
    "ruleId": 19385,
    "ruleName": "es-logs : unassigned_shards",
    "ruleUrl": "https://grafana-aiad2.wbx2.com/d/lvinkKkWk/platform-services?fullscreen=true&edit=true&tab=alert&panelId=76&orgId=1",
    "state": "no_data",
    "title": "[No Data] es-logs : unassigned_shards"
}' http://localhost:17575/webex/johndye@cisco.com$1
#?ignoreNoData=true

#\
#?briefMode=true

#Y2lzY29zcGFyazovL3VzL1BFT1BMRS8xOWE5YzhhNS1lZmRjLTRjNTgtODM4Yy1kNzU2OTk0YjQwN2E?briefMode=true
#