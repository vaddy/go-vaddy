package crawl

import (
	"errors"
	"testing"
	"vaddy/args"
	"vaddy/httpreq"
	"vaddy/test/mock/httpreqmock"
)

//Crawl検索で3件取得し、その中のItemの先頭のid:2が返ることを確認する
// total countが3件というのも確認する
func TestGetCrawlIdOkCrawlId2(t *testing.T) {
	var jsonResponse string = `{"total":3,"page":1,"limit":30,
	"items":[
		{"id":2, "label": "all-data", "start":"2016-03-12T11:11:11+0000", "end":"2016-03-12T12:00:00+0000"},
		{"id":1, "label": null, "start":"2016-03-12T11:11:11+0000", "end":"2016-03-12T12:00:00+0000"},
		{"id":3, "label": null, "start":"2016-03-12T11:11:11+0000", "end":"2016-03-12T12:00:00+0000"}
	]}`
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{200, []byte(jsonResponse), nil}
	httpreqmock.HttpResponseMockErr = nil

	httpRequestHandler = httpreqmock.HttpRequestData{} //search_crawl.go global var HttpRequestHander
	scanSetting := args.ScanSetting{}
	crawlId, count, err := GetCrawlId(scanSetting)
	if err != nil {
		t.Error("err is not nil")
	}
	if crawlId != "2" {
		t.Error("crawlId not 2")
	}
	if count != 3 {
		t.Error("crawl count not 3")
	}
}

//Crawl検索で0件だったためエラーとなるのをテスト
func TestGetCrawlIdNoItemError(t *testing.T) {
	var jsonResponse string = `{"total":0,"page":1,"limit":30,
	"items":[
	]}`
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{200, []byte(jsonResponse), nil}
	httpreqmock.HttpResponseMockErr = nil

	httpRequestHandler = httpreqmock.HttpRequestData{} //search_crawl.go global var HttpRequestHander
	scanSetting := args.ScanSetting{}
	_, _, err := GetCrawlId(scanSetting)
	if err == nil {
		t.Error("err is nil")
	}
	if err.Error() != `Can not find crawl id. using latest crawl id.` {
		t.Fatal("Expect no item error")
	}
}

//レスポンスのJSONがJSON形式ではないためJSON convert errorとなるかチェック
func TestGetCrawlIdNotJsonDataError(t *testing.T) {
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{200, []byte("not json data"), nil}
	httpreqmock.HttpResponseMockErr = nil

	httpRequestHandler = httpreqmock.HttpRequestData{} //search_crawl.go global var HttpRequestHander
	scanSetting := args.ScanSetting{}
	_, _, err := GetCrawlId(scanSetting)
	if err == nil {
		t.Error("err is nil")
	}
	if err.Error() != `Can not convert json response in crawl search.` {
		t.Fatal("Expect JSON convert error")
	}
}

// HTTPレスポンスがエラーとなっているのをテスト
func TestGetCrawlIdReturnError(t *testing.T) {
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{400, []byte("result body"), nil}
	httpreqmock.HttpResponseMockErr = errors.New("Error test")

	httpRequestHandler = httpreqmock.HttpRequestData{}
	scanSetting := args.ScanSetting{}
	_, _, err := GetCrawlId(scanSetting)
	if err == nil {
		t.Error("Expect error response status.")
	}
}
