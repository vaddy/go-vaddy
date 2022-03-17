package notification

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"vaddy/config"
	"vaddy/scan"
)

func PostSlackVulnerabilitiesWarning(alertCount int, fqdnOrProjectId string, scanID string, scanResultURL string) {
	title := fmt.Sprintf("VAddy Scan found Vulnerabilities: %d Warning!!!\n", alertCount)
	text := fqdnOrProjectId + "\n"
	text += "Scan ID: " + scanID + "\n"
	text += "Result URL: " + scanResultURL

	postSlack(title, text)
}

func PostSlackIncompleteNotice(fqdnOrProjectId string, scanID string, scanResult scan.ScanResult) {
	title := fmt.Sprintf("Notice: VAddy Scan was NOT complete (%d%%).\n", scanResult.Complete)
	text := fqdnOrProjectId + "\n"
	text += "Scan ID: " + scanID + "\n"
	text += "Result URL: " + scanResult.ScanResultUrl

	postSlack(title, text)
}

func postSlack(title, text string) {
	slackWebhookURL, ok1 := os.LookupEnv("SLACK_WEBHOOK_URL")

	if ok1 {
		slackUsername, ok2 := os.LookupEnv("SLACK_USERNAME")
		if !ok2 {
			slackUsername = ""
		}
		slackChannel, ok3 := os.LookupEnv("SLACK_CHANNEL")
		if !ok3 {
			slackChannel = ""
		}
		iconEmoji, ok4 := os.LookupEnv("SLACK_ICON_EMOJI")
		if !ok4 {
			iconEmoji = ""
		}
		iconURL, ok5 := os.LookupEnv("SLACK_ICON_URL")
		if !ok5 {
			iconURL = ""
		}

		type attachments struct {
			Color string `json:"color"`
			Title string `json:"title"`
			Text  string `json:"text"`
		}

		type slack struct {
			Username     string        `json:"username"`
			IconEmoji    string        `json:"icon_emoji"`
			IconURL      string        `json:"icon_url"`
			Channel      string        `json:"channel"`
			Text         string        `json:"text"`
			Attachements []attachments `json:"attachments"`
		}

		webhooks := slack{
			Username:  slackUsername,
			IconEmoji: iconEmoji,
			IconURL:   iconURL,
			Channel:   slackChannel,
			Text:      "VAddy Scan (Version " + config.VERSION + ")",
			Attachements: []attachments{
				{
					Color: "warning",
					Title: title,
					Text:  text,
				},
			},
		}

		params, _ := json.Marshal(webhooks)
		resp, err := http.PostForm(
			slackWebhookURL,
			url.Values{"payload": {string(params)}},
		)
		if err == nil {
			defer resp.Body.Close()
		}
	}
}