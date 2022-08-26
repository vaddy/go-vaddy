package httpreqmock

import (
	"net/url"
	"vaddy/args"
	"vaddy/httpreq"
)

type HttpRequestData struct {
}

var HttpResponseMock = httpreq.HttpResponseData{}
var HttpResponseMockErr error = nil

func (hrd HttpRequestData) HttpGet(urlpath string, scanSetting args.ScanSetting, values url.Values) (httpreq.HttpResponseData, error) {
	return HttpResponseMock, HttpResponseMockErr
}
func (hrd HttpRequestData) HttpPost(urlpath string, scanSetting args.ScanSetting, values url.Values) (httpreq.HttpResponseData, error) {
	return HttpResponseMock, HttpResponseMockErr
}
