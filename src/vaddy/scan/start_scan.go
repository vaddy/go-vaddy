package scan

import (
	"net/url"
	"vaddy/args"
	"vaddy/common"
)

type ScanId struct {
	ScanID string `json:"scan_id"`
}

func StartScan(scanSetting args.ScanSetting) (string, error) {
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
	return scanId, nil
}

func getScanId(jsonByteData []byte) string {
	var scan_result ScanId
	common.ConvertJsonToStruct(jsonByteData, &scan_result)
	return scan_result.ScanID
}
