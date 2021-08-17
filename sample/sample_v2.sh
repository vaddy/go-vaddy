#/bin/sh

export VADDY_TOKEN="123455667789"
export VADDY_USER="test user"
export VADDY_PROJECT_ID="your project id"

#### crawl ID (optional)
#export VADDY_CRAWL="30"

#### slack setting (optional)
#export SLACK_WEBHOOK_URL="webhook url"
#export SLACK_USERNAME="your user (optional)"
#export SLACK_CHANNEL="your channel (optional)"
#export SLACK_ICON_EMOJI=":smile: (optional)"
#export SLACK_ICON_URL="icon url (optional)"

../bin/vaddy-linux-64bit

