package crawl

import (
	"errors"
	"fmt"
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

func GetCrawlId(httpReq httpreq.HttpRequestData, scanSetting args.ScanSetting) (string, error) {
	json_response, err := doCrawlSearch(httpReq, scanSetting)
	//fmt.Println(string(json_response))
	if err != nil {
		return "", err
	}

	var crawl_result CrawlSearch
	json_err := common.ConvertJsonToStruct(json_response, &crawl_result)
	if json_err != nil {
		return "", errors.New("Can not convert json response in crawl search.")
	}
	if crawl_result.Total == 0 {
		return "", errors.New("can not find crawl id. using latest crawl id.")
	}
	var crawl_id int = crawl_result.Items[0].CrawlId
	fmt.Printf("Found %d results. Using CrawlID: %d \n\n", crawl_result.Total, crawl_id)
	return strconv.Itoa(crawl_id), nil
}

func doCrawlSearch(httpReq httpreq.HttpRequestData, scanSetting args.ScanSetting) ([]byte, error) {
	values := url.Values{}
	values.Add("user", scanSetting.User)
	values.Add("fqdn", scanSetting.Fqdn)
	values.Add("search_label", scanSetting.Crawl)
	values.Add("verification_code", scanSetting.VerificationCode)
	values.Add("project_id", scanSetting.ProjectId)

	result, err := httpReq.HttpGet("/crawl", scanSetting, values)
	if err != nil {
		return nil, err
	}
	err2 := common.CheckHttpResponse(result)
	return result.Body, err2
}
