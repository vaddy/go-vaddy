package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"
	"vaddy/args"
	"vaddy/common"
	"vaddy/config"
	"vaddy/crawl"
	"vaddy/httpreq"
	"vaddy/notification"
	"vaddy/scan"
)

const SUCCESS_EXIT int = 0
const ERROR_EXIT int = 1

func init() {
	var showVersion bool
	flag.BoolVar(&showVersion, "version", false, "Show version")
	flag.BoolVar(&showVersion, "v", false, "Show version")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s", config.VERSION)
		os.Exit(SUCCESS_EXIT)
	}
}

func main() {
	fmt.Println("==== Start VAddy Scan (Version " + config.VERSION + ")====")

	scanSetting, err := args.GetArgsFromEnv()
	if err != nil {
		fmt.Println("Missing arguments or system env.")
		fmt.Println(" ENV VADDY_USER: " + scanSetting.User)
		fmt.Println(" ENV VADDY_PROJECT_ID(V2): " + scanSetting.ProjectId)
		fmt.Println(" ENV VADDY_HOST(V1): " + scanSetting.ScanType)
		os.Exit(ERROR_EXIT)
	}

	if scanSetting.CheckNeedToGetCrawlId() {
		fmt.Println("Start to get crawl ID from keyword: " + scanSetting.Crawl)
		scanSetting.CrawlId, err = crawl.GetCrawlId(httpreq.HttpRequestData{}, scanSetting)
		if err != nil {
			fmt.Println(err)
			os.Exit(ERROR_EXIT)
		}
	}

	scanId, err := scan.StartScan(httpreq.HttpRequestData{}, scanSetting)
	if err != nil {
		fmt.Println(err)
		os.Exit(ERROR_EXIT)
	}
	fmt.Println("Scan Started. ScanID: " + scanId)

	var waitCount int = 0
	for {
		scanResult, err := scan.GetScanResult(httpreq.HttpRequestData{}, scanSetting, scanId)
		if err != nil {
			fmt.Println(err)
			os.Exit(ERROR_EXIT)
		}

		if scanResult.IsScanRunning() {
			common.PrintDots(waitCount)
		}

		if scanResult.IsCanceled(){
			fmt.Println("Scan canceled.")
			fmt.Println("Result URL: " + scanResult.ScanResultUrl)
			os.Exit(ERROR_EXIT)
		}

		if scanResult.IsFinished(){
			fmt.Println(".")
			fqdnOrProjectIdLabel := scanSetting.GetHostOrProjectName()
			fmt.Println(fqdnOrProjectIdLabel)
			fmt.Println("scanId: " + scanId)
			fmt.Println("Result URL: " + scanResult.ScanResultUrl)
			fmt.Println("Crawl ID: " + strconv.Itoa(scanResult.CrawlId))
			fmt.Println("Crawl Label: " + scanResult.CrawlLabel)
			scanSetting.PrintScanTypeSetting()

			if scanResult.IsVulnExist() {
				fmt.Print("Vulnerabilities: ")
				fmt.Println(scanResult.AlertCount)
				fmt.Println("Warning!!!")
				notification.PostSlackVulnerabilitiesWarning(scanResult.AlertCount, fqdnOrProjectIdLabel, scanId, scanResult.ScanResultUrl)
				os.Exit(ERROR_EXIT)
			} else if scanResult.IsIncomplete() {
				fmt.Printf("Notice: Scan was NOT complete (%d%%).\n", scanResult.Complete)
				fmt.Println("No vulnerabilities.")
				notification.PostSlackIncompleteNotice(fqdnOrProjectIdLabel, scanId, scanResult)
				os.Exit(SUCCESS_EXIT)
			} else {
				fmt.Println("Scan Success. No vulnerabilities!")
				os.Exit(SUCCESS_EXIT)
			}
		}

		sleepSec := common.GetSleepSec(waitCount)
		time.Sleep(time.Duration(sleepSec) * time.Second)

		if isTimeOver(waitCount) {
			fmt.Println("Error: time out")
			os.Exit(ERROR_EXIT)
		}
		waitCount++
	}
}

func isTimeOver(waitCount int) bool {
	if waitCount > config.LIMIT_WAIT_COUNT {
		return true
	}
	return false
}