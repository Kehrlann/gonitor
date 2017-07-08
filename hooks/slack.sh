#!/bin/bash
#
# Hook for Gonitor, for posting to slack channels.
# For more details, see : https://github.com/Kehrlann/gonitor/blob/master/hooks/README.md
#

SLACK_URL=http://your.slack/url

json="{\"text\" : \"[$4] $1, with codes $3\"}"
curl -X POST -H 'Content-type: application/json' --data "$json" $SLACK_URL
