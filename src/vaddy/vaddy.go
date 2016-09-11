package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
	"regexp"
	"strconv"
)
const VERSION string = "1.0.1"
const SUCCESS_EXIT int = 0
const ERROR_EXIT int = 1
const LIMIT_WAIT_COUNT int = 540 // 20sec * 540 = 3 hours
const API_SERVER string = "https://api.vaddy.net"

type CrawlSearch struct {
	Total int `json:"total"`
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
	ScanResultUrl string `json:"scan_result_url"`
}

func main() {
	fmt.Println("==== Start VAddy Scan (Version " + VERSION + ")====")

	var auth_key, user, fqdn, crawl string = getApiParamsFromArgsOrEnv()

	if checkNeedToGetCrawlId(crawl) {
		fmt.Println("Start to get crawl ID from keyword: " + crawl)
		crawl = getCrawlId(auth_key, user, fqdn, crawl)
	}

	scan_id := startScan(auth_key, user, fqdn, crawl)

	var wait_count int = 0
	for {
		checkScanResult(auth_key, user, fqdn, scan_id)

		time.Sleep(20 * time.Second) //wait 20 second
		wait_count++
		if wait_count > LIMIT_WAIT_COUNT {
			fmt.Println("Error: time out")
			os.Exit(ERROR_EXIT)
		}
	}
}

func getApiParamsFromArgsOrEnv() (string, string, string, string) {
	var auth_key, user, fqdn, crawl string
	if len(os.Args) < 4 {
		return getArgsFromEnv()
	}

	auth_key = os.Args[1]
	user = os.Args[2]
	fqdn = os.Args[3]
	if len(os.Args) >= 5 {
		crawl = os.Args[4]
	}
	return auth_key, user, fqdn, crawl
}

func getArgsFromEnv() (string, string, string, string) {
	var auth_key, user, fqdn, crawl string
	auth_key, ok1 := os.LookupEnv("VADDY_TOKEN")
	user, ok2 := os.LookupEnv("VADDY_USER")
	fqdn, ok3 := os.LookupEnv("VADDY_HOST")
	crawl, _ = os.LookupEnv("VADDY_CRAWL")

	if !ok1 || !ok2 || !ok3 {
		fmt.Println("Missing arguments or system env.")
		fmt.Println("USAGE: vaddy.go ApiKey UserId FQDN CrawlID/Label(optional)")
		os.Exit(ERROR_EXIT)
	}
	return auth_key, user, fqdn, crawl
}

func startScan(auth_key string, user string, fqdn string, crawl string) string {
	values := url.Values{}
	values.Add("auth_key", auth_key)
	values.Add("user", user)
	values.Add("fqdn", fqdn)
	values.Add("action", "start")
	if len(crawl) > 0 {
		values.Add("crawl_id", crawl)
	}

	api_server := getApiServerName()
	res, err := http.PostForm(api_server + "/v1/scan", values)
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

func getScanResult(auth_key string, user string, fqdn string, scan_id string) []byte {
	values := url.Values{}
	values.Add("auth_key", auth_key)
	values.Add("user", user)
	values.Add("fqdn", fqdn)
	values.Add("scan_id", scan_id)

	api_server := getApiServerName()
	res, err := http.Get(api_server + "/v1/scan/result?" + values.Encode())
	if err != nil {
		fmt.Println(err)
		os.Exit(ERROR_EXIT)
	}
	defer res.Body.Close()

	json_response := getResponseData(res)
	return json_response
}

func checkScanResult(auth_key string, user string, fqdn string, scan_id string) {
	json_response := getScanResult(auth_key, user, fqdn, scan_id)

	var scan_result ScanResult
	convertJsonToStruct(json_response, &scan_result)

	status := scan_result.Status
	switch status {
	case "scanning":
		fmt.Println(scan_result.Status)
	case "canceled":
		fmt.Println(scan_result.Status)
		os.Exit(ERROR_EXIT)
	case "finish":
		//fmt.Println(string(json_response) + "\n")
		fmt.Println("Server: " + fqdn)
		fmt.Println("scanId: " + scan_id)
		fmt.Println("Result URL: " + scan_result.ScanResultUrl)

		if scan_result.AlertCount > 0 {
			fmt.Print("Vulnerabilities: ")
			fmt.Println(scan_result.AlertCount)
			fmt.Println("Warning!!!")
			os.Exit(ERROR_EXIT)
		} else {
			fmt.Println("Scan Success. No vulnerabilities!")
			os.Exit(SUCCESS_EXIT)
		}
	}
}


func getCrawlId(auth_key string, user string, fqdn string, search_label string) string {
	json_response := doCrawlSearch(auth_key, user, fqdn, search_label)
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

func doCrawlSearch(auth_key string, user string, fqdn string, search_label string) []byte {
	values := url.Values{}
	values.Add("auth_key", auth_key)
	values.Add("user", user)
	values.Add("fqdn", fqdn)
	values.Add("search_label", search_label)

	api_server := getApiServerName()
	res, err := http.Get(api_server + "/v1/crawl?" + values.Encode())
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
		return api_server
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

