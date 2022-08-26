package scan

import (
	"errors"
	"testing"
	"vaddy/args"
	"vaddy/httpreq"
	"vaddy/test/mock/httpreqmock"
)

func TestStartScanOK(t *testing.T) {
	var jsonResponse string = `{"scan_id": "vaddy-scanid-12345"}`
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{200, []byte(jsonResponse), nil}
	httpRequestHandler = httpreqmock.HttpRequestData{} //scan_http_request_hander.go global var HttpRequestHander
	scanSetting := args.ScanSetting{}
	scanId, err := StartScan(scanSetting)

	if err != nil {
		t.Error("err is not nil")
	}
	if scanId != "vaddy-scanid-12345" {
		t.Error("scanId check error")
	}
}

//HTTP statusが400で、あえてレスポンスのMockErrにnilをセットしてエラーを検証
func TestStartScanNG_checkResponseStatus400(t *testing.T) {
	var jsonResponse string = `{"error_message":"error response"}`
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{400, []byte(jsonResponse), nil}
	httpreqmock.HttpResponseMockErr = nil

	scanSetting := args.ScanSetting{}
	scanId, err := StartScan(scanSetting)

	if err == nil {
		t.Error("err is nil")
	}
	if scanId != "" {
		t.Error("scanId not empty. expect empty")
	}

	errorMessage := `Network/Auth error
{"error_message":"error response"}`

	if err.Error() != errorMessage {
		t.Fatal(err.Error())
	}
}

// レスポンス200だがScanIDが見つからないレスポンスbodyを返してエラーをチェック
func TestStartScanNG_checkBrokenJson(t *testing.T) {
	var jsonResponse string = `{"error_message"`
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{200, []byte(jsonResponse), nil}
	httpreqmock.HttpResponseMockErr = nil

	scanSetting := args.ScanSetting{}
	scanId, err := StartScan(scanSetting)

	if err == nil {
		t.Error("err is nil")
	}
	if err.Error() != `StartScan Error: No Scan ID found.` {
		t.Fatal(err.Error())
	}
	if scanId != "" {
		t.Error("scanId not empty. expect empty")
	}
}

//MockErrにエラーをセットしてエラーチェック
func TestStartScanNG2_checkErrorReturn(t *testing.T) {
	var jsonResponse string = `{"error_message":"test error"}`
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{400, []byte(jsonResponse), nil}
	httpreqmock.HttpResponseMockErr = errors.New("Error test")
	scanSetting := args.ScanSetting{}
	scanId, err := StartScan(scanSetting)

	if err == nil {
		t.Error("err is nil")
	}
	if scanId != "" {
		t.Error("scanId not empty. expect empty")
	}
	if err.Error() != `Error test` {
		t.Fatal(err.Error())
	}
}
