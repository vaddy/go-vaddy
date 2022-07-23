package crawl

import (
	"errors"
	"net/url"
	"strconv"
	"vaddy/args"
	"vaddy/common"
	"vaddy/httpreq"
)

type CrawlSearch struct {
	Total int               `json:"total"`
	Items []CrawlSearchItem `json:"items"`
}

type CrawlSearchItem struct {
	CrawlId int `json:"id"`
}

var httpRequestHandler httpreq.HttpReqInterface = httpreq.HttpRequestData{}

func GetCrawlId(scanSetting args.ScanSetting) (string, int, error) {
	json_response, err := doCrawlSearch(scanSetting)
	//fmt.Println(string(json_response))
	if err != nil {
		return "", 0, err
	}

	var crawl_result CrawlSearch
	json_err := common.ConvertJsonToStruct(json_response, &crawl_result)
	if json_err != nil {
		return "", 0, errors.New("Can not convert json response in crawl search.")
	}
	if crawl_result.Total == 0 {
		return "", 0, errors.New("Can not find crawl id. using latest crawl id.")
	}
	var crawl_id int = crawl_result.Items[0].CrawlId
	return strconv.Itoa(crawl_id), crawl_result.Total, nil
}

func doCrawlSearch(scanSetting args.ScanSetting) ([]byte, error) {
	values := url.Values{}
	values.Add("user", scanSetting.User)
	values.Add("fqdn", scanSetting.Fqdn)
	values.Add("search_label", scanSetting.Crawl)
	values.Add("verification_code", scanSetting.VerificationCode)
	values.Add("project_id", scanSetting.ProjectId)

	result, err := httpRequestHandler.HttpGet("/crawl", scanSetting, values)
	if err != nil {
		return nil, err
	}
	err2 := common.CheckHttpResponse(result)
	return result.Body, err2
}
