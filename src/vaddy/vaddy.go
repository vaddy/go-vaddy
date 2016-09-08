package vaddy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

const SUCCESS_EXIT int = 0
const ERROR_EXIT int = 1
const LIMIT_WAIT_COUNT int = 540 // 20sec * 540 = 3 hours
const API_SERVER string = "https://api.vaddy.net"

type StartScan struct {
	ScanID string `json:"scan_id"`
}

type ScanResult struct {
	Status        string `json:"status"`
	AlertCount    int    `json:"alert_count"`
	ScanResultUrl string `json:"scan_result_url"`
}

func main() {
	fmt.Println("==== Start VAddy Scan ====")

	var auth_key, user, fqdn, crawl string = getApiParamsFromArgsOrEnv()

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
		fmt.Println("USAGE: vaddy ApiKey UserId FQDN CrawlID(optional)")
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

	apiServer := getApiServerName()
	res, err := http.PostForm(apiServer+"/v1/scan", values)
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

	apiServer := getApiServerName()
	res, err := http.Get(apiServer + "/v1/scan/result?" + values.Encode())
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
	apiserver, ok := os.LookupEnv("VADDY_API_SERVER")
	if ok {
		return apiserver
	}
	return API_SERVER
}
