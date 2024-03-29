package scan

import (
	"errors"
	"net/url"
	"vaddy/args"
	"vaddy/common"
)

type ScanId struct {
	ScanID string `json:"scan_id"`
}

func StartScan(scanSetting args.ScanSetting) (scanIdString string, errorVal error) {
	values := url.Values{}
	values.Add("action", "start")
	values.Add("user", scanSetting.User)
	values.Add("fqdn", scanSetting.Fqdn)
	values.Add("verification_code", scanSetting.VerificationCode)
	values.Add("project_id", scanSetting.ProjectId)
	values.Add("scan_type", scanSetting.ScanType)
	values.Add("crawl_id", scanSetting.CrawlId)

	result, err := httpRequestHandler.HttpPost("/scan", scanSetting, values)
	if err != nil {
		return "", err
	}
	err2 := common.CheckHttpResponse(result)

	if err2 != nil {
		return "", err2
	}
	scanId := getScanId(result.Body)
	//fmt.Println("scanId: " + scanId)
	if scanId == "" {
		return "", errors.New("StartScan Error: No Scan ID found.")
	}
	return scanId, nil
}

func getScanId(jsonByteData []byte) (scanIdString string) {
	var scan_result ScanId
	common.ConvertJsonToStruct(jsonByteData, &scan_result)
	return scan_result.ScanID
}
