package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"time"
)

const VERSION string = "1.0.6"
const SUCCESS_EXIT int = 0
const ERROR_EXIT int = 1
const LIMIT_WAIT_COUNT int = 600 // 20sec * 600 = 3.3 hours
const API_SERVER string = "https://api.vaddy.net"

type CrawlSearch struct {
	Total int               `json:"total"`
	Items []CrawlSearchItem `json:"items"`
}

type CrawlSearchItem struct {
	CrawlId int `json:"id"`
}

type StartScan struct {
	ScanID string `json:"scan_id"`
}

type ScanResult struct {
	Status        string `json:"status"`
	AlertCount    int    `json:"alert_count"`
	ScanCount     int    `json:"scan_count"`
	ScanResultUrl string `json:"scan_result_url"`
	Complete      int    `json:"complete"`
}

func (s ScanResult) IsIncomplete() bool {
	return s.AlertCount == 0 && s.Complete < 100
}

func init() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("Current version is %s", VERSION)
		os.Exit(SUCCESS_EXIT)
	}
}

func main() {
	fmt.Println("==== Start VAddy Scan (Version " + VERSION + ")====")

	var auth_key, user, fqdn, crawl, verification_code, project_id string = getApiParamsFromArgsOrEnv()

	if checkNeedToGetCrawlId(crawl) {
		fmt.Println("Start to get crawl ID from keyword: " + crawl)
		crawl = getCrawlId(auth_key, user, fqdn, crawl, verification_code, project_id)
	}

	scan_id := startScan(auth_key, user, fqdn, crawl, verification_code, project_id)

	var wait_count int = 0
	var sleep_sec int = 20

	for {
		checkScanResult(auth_key, user, fqdn, scan_id, wait_count, verification_code, project_id)

		sleep_sec = 20
		if wait_count < 10 {
			sleep_sec = 3
		}
		time.Sleep(time.Duration(sleep_sec) * time.Second)

		wait_count++
		if wait_count > LIMIT_WAIT_COUNT {
			fmt.Println("Error: time out")
			os.Exit(ERROR_EXIT)
		}
	}
}

func getApiParamsFromArgsOrEnv() (string, string, string, string, string, string) {
	var auth_key, user, fqdn, crawl, verification_code string
	verification_code, _ = os.LookupEnv("VADDY_VERIFICATION_CODE")

	if len(os.Args) < 4 {
		return getArgsFromEnv(verification_code)
	}

	auth_key = os.Args[1]
	user = os.Args[2]
	fqdn = os.Args[3]
	if len(os.Args) >= 5 {
		crawl = os.Args[4]
	}
	return auth_key, user, fqdn, crawl, verification_code, ""
}

func getArgsFromEnv(verification_code string) (string, string, string, string, string, string) {
	var auth_key, user, fqdn, crawl string
	auth_key, ok1 := os.LookupEnv("VADDY_TOKEN")
	user, ok2 := os.LookupEnv("VADDY_USER")
	fqdn, ok3 := os.LookupEnv("VADDY_HOST")
	crawl, _ = os.LookupEnv("VADDY_CRAWL")
	project_id, ok4 := os.LookupEnv("VADDY_PROJECT_ID")

	// v1
	if ok1 && ok2 && ok3 {
		return auth_key, user, fqdn, crawl, verification_code, ""
	}
	// v2
	if ok1 && ok2 && ok4 {
		return auth_key, user, "", crawl, verification_code, project_id
	}

	fmt.Println("Missing arguments or system env.")
	fmt.Println("USAGE: vaddy.go ApiKey UserId FQDN CrawlID/Label(optional)")
	os.Exit(ERROR_EXIT)

	return "", "", "", "", "", ""
}

func startScan(auth_key string, user string, fqdn string, crawl string, verification_code, project_id string) string {
	values := url.Values{}
	values.Add("auth_key", auth_key)
	values.Add("user", user)
	values.Add("fqdn", fqdn)
	values.Add("action", "start")
	values.Add("verification_code", verification_code)
	if len(crawl) > 0 {
		values.Add("crawl_id", crawl)
	}
	if len(project_id) > 0 {
		values.Add("project_id", project_id)
	}

	api_server := getApiServerName()
	api_version := detectApiVersion(project_id)
	res, err := http.PostForm(api_server+"/"+api_version+"/scan", values)
	if err != nil {
		fmt.Println(err)
		os.Exit(ERROR_EXIT)
	}
	defer res.Body.Close()
	json_response := getResponseData(res)
	scanId := getScanId(json_response)
	//fmt.Println("scanId: " + scanId)
	return scanId
}

func getScanResult(auth_key string, user string, fqdn string, scan_id string, verification_code, project_id string) []byte {
	values := url.Values{}
	values.Add("auth_key", auth_key)
	values.Add("user", user)
	values.Add("fqdn", fqdn)
	values.Add("scan_id", scan_id)
	values.Add("verification_code", verification_code)
	if len(project_id) > 0 {
		values.Add("project_id", project_id)
	}

	api_server := getApiServerName()
	api_version := detectApiVersion(project_id)
	res, err := http.Get(api_server + "/" + api_version + "/scan/result?" + values.Encode())
	if err != nil {
		fmt.Println(err)
		os.Exit(ERROR_EXIT)
	}
	defer res.Body.Close()

	json_response := getResponseData(res)
	return json_response
}

func checkScanResult(auth_key string, user string, fqdn string, scan_id string, count int, verification_code, project_id string) {
	json_response := getScanResult(auth_key, user, fqdn, scan_id, verification_code, project_id)

	var scan_result ScanResult
	convertJsonToStruct(json_response, &scan_result)

	status := scan_result.Status
	switch status {
	case "scanning":
		if count > 0 && (count%60 == 0) { //wrap every 60 dots.
			fmt.Println(".")
		} else {
			fmt.Print(".")
		}
	case "canceled":
		fmt.Println(scan_result.Status)
		os.Exit(ERROR_EXIT)
	case "finish":
		//fmt.Println(string(json_response) + "\n")
		fmt.Println(".")
		fmt.Println("Server: " + fqdn)
		fmt.Println("scanId: " + scan_id)
		fmt.Println("Result URL: " + scan_result.ScanResultUrl)

		if scan_result.AlertCount > 0 {
			fmt.Print("Vulnerabilities: ")
			fmt.Println(scan_result.AlertCount)
			fmt.Println("Warning!!!")
			postSlackVulnerabilitiesWarning(scan_result.AlertCount, fqdn, scan_id, scan_result.ScanResultUrl)
			os.Exit(ERROR_EXIT)
		} else if scan_result.IsIncomplete() {
			fmt.Printf("Notice: Scan was NOT complete (%d%%).\n", scan_result.Complete)
			fmt.Println("No vulnerabilities.")
			postSlackIncompleteNotice(fqdn, scan_id, scan_result)
			os.Exit(SUCCESS_EXIT)
		} else if scan_result.ScanCount == 0 {
			fmt.Println("ERROR: VAddy was not able to scan your sever. Check the result on the Result URL.")
			os.Exit(ERROR_EXIT)
		} else {
			fmt.Println("Scan Success. No vulnerabilities!")
			os.Exit(SUCCESS_EXIT)
		}
	}
}

func getCrawlId(auth_key string, user string, fqdn string, search_label string, verification_code, project_id string) string {
	json_response := doCrawlSearch(auth_key, user, fqdn, search_label, verification_code, project_id)
	//fmt.Println(string(json_response))

	var crawl_result CrawlSearch
	convertJsonToStruct(json_response, &crawl_result)
	if crawl_result.Total == 0 {
		fmt.Println("can not find crawl id. using latest crawl id.")
		return ""
	}
	var crawl_id int = crawl_result.Items[0].CrawlId
	fmt.Printf("Found %d results. Using CrawlID: %d \n\n", crawl_result.Total, crawl_id)
	return strconv.Itoa(crawl_id)
}

func doCrawlSearch(auth_key string, user string, fqdn string, search_label string, verification_code, project_id string) []byte {
	values := url.Values{}
	values.Add("auth_key", auth_key)
	values.Add("user", user)
	values.Add("fqdn", fqdn)
	values.Add("search_label", search_label)
	values.Add("verification_code", verification_code)
	if len(project_id) > 0 {
		values.Add("project_id", project_id)
	}

	api_server := getApiServerName()
	api_version := detectApiVersion(project_id)
	res, err := http.Get(api_server + "/" + api_version + "/crawl?" + values.Encode())
	if err != nil {
		fmt.Println(err)
		os.Exit(ERROR_EXIT)
	}
	defer res.Body.Close()

	json_response := getResponseData(res)
	return json_response
}

func getResponseData(resp *http.Response) []byte {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(string(body))
		os.Exit(ERROR_EXIT)
	}
	status_code := resp.StatusCode
	//fmt.Println(status_code)
	if status_code != 200 {
		fmt.Println("Network/Auth error\n" + string(body))
		os.Exit(ERROR_EXIT)
	}
	return []byte(body)
}

func getScanId(jsonByteData []byte) string {
	var scan_result StartScan
	convertJsonToStruct(jsonByteData, &scan_result)
	return scan_result.ScanID
}

func convertJsonToStruct(jsonByteData []byte, structData interface{}) {
	err := json.Unmarshal(jsonByteData, structData)
	if err != nil {
		fmt.Println(err)
		os.Exit(ERROR_EXIT)
	}
}

func getApiServerName() string {
	api_server, ok := os.LookupEnv("VADDY_API_SERVER")
	if ok {
		//if api server does not contain protocol name, add https://
		var regex string = `^https://.*$`
		if regexp.MustCompile(regex).Match([]byte(api_server)) {
			return api_server
		}
		return "https://" + api_server
	}
	return API_SERVER
}

func checkNeedToGetCrawlId(str string) bool {
	if len(str) == 0 || str == "" {
		return false
	}
	var regex string = `[^0-9]`
	return regexp.MustCompile(regex).Match([]byte(str))
}

func postSlackVulnerabilitiesWarning(alertCount int, fqdn string, scanID string, scanResultURL string) {
	title := "VAddy Scan Vulnerabilities: " + string(alertCount) + " Warning!!!\n"
	text := "Server: " + fqdn + "\n"
	text += "Scan ID: " + scanID + "\n"
	text += "Result URL: " + scanResultURL

	postSlack(title, text)
}

func postSlackIncompleteNotice(fqdn string, scanID string, scanResult ScanResult) {
	title := fmt.Sprintf("Notice: VAddy Scan was NOT complete (%d%%).\n", scanResult.Complete)
	text := "Server: " + fqdn + "\n"
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
			Text:      "VAddy Scan (Version " + VERSION + ")",
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

func detectApiVersion(project_id string) string {
	if len(project_id) > 0 {
		return "v2"
	}

	return "v1"
}
