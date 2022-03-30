package crawl

import (
	"errors"
	"fmt"
	"testing"
	"vaddy/args"
	"vaddy/httpreq"
	"vaddy/test/mock/httpreqmock"
)

func TestName1(t *testing.T) {
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{200, []byte("result body"), nil}
	httpreqmock.HttpResponseMockErr = nil
	//httpreqmock.HttpResponseMockErr = errors.New("aa")

	httpRequestHandler = httpreqmock.HttpRequestData{} //search_crawl.go global var HttpRequestHander
	scanSetting := args.ScanSetting{}
	result, err := doCrawlSearch(scanSetting)
	fmt.Println(result)
	if err != nil {
		t.Fail()
	}
}
func TestName2(t *testing.T) {
	httpreqmock.HttpResponseMock = httpreq.HttpResponseData{200, []byte("result body"), nil}
	httpreqmock.HttpResponseMockErr = errors.New("aa")

	httpRequestHandler = httpreqmock.HttpRequestData{}
	scanSetting := args.ScanSetting{}
	result, err := doCrawlSearch(scanSetting)
	fmt.Println(result)
	if err != nil {
		t.Fail()
	}
}
