package scan

import (
	"errors"
	"fmt"
	"net/url"
	"time"
	"vaddy/args"
	"vaddy/common"
	"vaddy/config"
	"vaddy/httpreq"
)

func GetScanResult(httpReq httpreq.HttpRequestData, scanSetting args.ScanSetting, scanId string) (ScanResult, error) {
	values := url.Values{}
	values.Add("user", scanSetting.User)
	values.Add("fqdn", scanSetting.Fqdn)
	values.Add("verification_code", scanSetting.VerificationCode)
	values.Add("project_id", scanSetting.ProjectId)
	values.Add("scan_id", scanId)

	var scanResult ScanResult
	var retryCount int = 0
	for {
		result, err := httpReq.HttpGet("/scan/result", scanSetting, values)
		if err != nil {
			fmt.Print("HTTP Get Request Error: /scan/result. ")
			fmt.Print(err)
			//リトライするためにerrがあってもprintのみする
		}

		err2 := common.CheckHttpResponse(result)
		if err2 != nil {
			fmt.Print("HTTP Get Response Error: /scan/result. ")
			fmt.Print(err2)
			//リトライするためにerr2があってもprintのみする
		}

		if err == nil && err2 == nil {
			common.ConvertJsonToStruct(result.Body, &scanResult)
			return scanResult, nil
		}

		retryCount++
		if retryCount > config.NETWORK_RETRY_COUNT {
			fmt.Printf("-- getScanResult() retry max count: %d exit. --", retryCount)
			fmt.Println(err)
			return scanResult, errors.New("Error: get scan result retry over.")
		}
		fmt.Printf("-- getScanResult() HTTP GET error: count %d --\n", retryCount)
		time.Sleep(time.Duration(config.NETWORK_RETRY_WAIT_TIME) * time.Second)
	}
}