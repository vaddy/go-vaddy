package scan

import (
	"strconv"
	"testing"
	"vaddy/args"
	"vaddy/httpreq"
	"vaddy/test/mock/httpreqmock"
)

// 正常な検査完了レスポンスが返るテスト
func TestCheckScanResultOk_ScanFinished(t *testing.T) {
	var jsonResponse string = `{ "status":"finish",
  "project_id" : "6eb1f9fcbdb6a5a",
  "scan_id" : "1-837b5f9f-e088-4af5-9491-67f7ce8035a4",
  "scan_count" : 22,
  "alert_count" : 1,
  "timezone" : "UTC",
  "start_time" :  "2014-06-12T11:11:11+0000",
  "end_time" :  "2014-06-12T11:11:11+0000",
  "scan_result_url" : "https://console.vaddy.net/scan_status/1/1-837b5f9f-e088-4af5-9491-67f7ce8035a4",
  "complete" : 100,
  "crawl_id" : 30,
  "crawl_label" : "ユーザ情報修正シナリオ",
  "scan_list" : ["XSS","SQL Injection"]
}`
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{200, []byte(jsonResponse), nil}
	httpRequestHandler = httpreqmock.HttpRequestData{} //scan_http_request_hander.go global var HttpRequestHander
	scanSetting := args.ScanSetting{}
	result, err := GetScanResult(scanSetting, "1-837b5f9f-e088-4af5-9491-67f7ce8035a4")

	if err != nil {
		t.Error("err is not nil. " + err.Error())
	}
	if !result.IsFinished() {
		t.Error("Scan Results not finish")
	}
	if !result.IsVulnExist() {
		t.Error("Scan Results alert count error. alert count: " + strconv.Itoa(result.AlertCount))
	}
}

// 検査がキャンセルされた結果が返るテスト
func TestCheckScanResultOk_ScanCanceled(t *testing.T) {
	var jsonResponse string = `{"status":"canceled"}`
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{200, []byte(jsonResponse), nil}
	httpRequestHandler = httpreqmock.HttpRequestData{} //scan_http_request_hander.go global var HttpRequestHander
	scanSetting := args.ScanSetting{}
	result, err := GetScanResult(scanSetting, "1-837b5f9f-e088-4af5-9491-67f7ce8035a4")

	if err != nil {
		t.Error("err is not nil. " + err.Error())
	}
	if !result.IsCanceled() {
		t.Error("Scan Results not cancel")
	}
}

// 検査実行中の結果が返るテスト
func TestCheckScanResultOk_ScanScanning(t *testing.T) {
	var jsonResponse string = `{"status":"scanning"}`
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{200, []byte(jsonResponse), nil}
	httpRequestHandler = httpreqmock.HttpRequestData{} //scan_http_request_hander.go global var HttpRequestHander
	scanSetting := args.ScanSetting{}
	result, err := GetScanResult(scanSetting, "1-837b5f9f-e088-4af5-9491-67f7ce8035a4")

	if err != nil {
		t.Error("err is not nil. " + err.Error())
	}
	if !result.IsScanRunning() {
		t.Error("Scan Results not running")
	}
}

// ネットワークやサーバエラーで5回リトライをテスト
func TestCheckScanResultOk_ScanError(t *testing.T) {
	var jsonResponse string = `{"error_message":"xxxxxx"}`
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{400, []byte(jsonResponse), nil}
	httpRequestHandler = httpreqmock.HttpRequestData{} //scan_http_request_hander.go global var HttpRequestHander
	scanSetting := args.ScanSetting{}
	result, err := GetScanResult(scanSetting, "1-837b5f9f-e088-4af5-9491-67f7ce8035a4")

	if err == nil {
		t.Error("err is nil.")
	}
	if err.Error() != `Error: get scan result retry over.` {
		t.Error("Expect retry over error message")
	}
	if result.Status != "" {
		t.Error("scan status expect empty")
	}
}
