package httpreq

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"runtime"
	"strings"
	"vaddy/args"
	"vaddy/config"
)

type HttpReqInterface interface {
	HttpGet(urlpath string, scanSetting args.ScanSetting, values url.Values) (HttpResponseData, error)
	HttpPost(urlpath string, scanSetting args.ScanSetting, values url.Values) (HttpResponseData, error)
}

type HttpRequestData struct {
}

type HttpResponseData struct {
	Status int
	Body   []byte
	Error  error
}

func (hrd HttpRequestData) HttpGet(urlpath string, scanSetting args.ScanSetting, values url.Values) (HttpResponseData, error) {
	api_server := GetApiServerName()

	params := values.Encode()
	var endpoint string = api_server + "/" + scanSetting.GetApiVersion() + urlpath + "?" + params

	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return HttpResponseData{}, err
	}
	req.Header.Add("X-API-KEY", scanSetting.AuthKey)
	req.Header.Add("User-Agent", createUserAgentString())

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	return getResponseData(resp, err), err
}

func (hrd HttpRequestData) HttpPost(urlpath string, scanSetting args.ScanSetting, values url.Values) (HttpResponseData, error) {
	api_server := GetApiServerName()

	var endpoint string = api_server + "/" + scanSetting.GetApiVersion() + urlpath

	req, err := http.NewRequest(http.MethodPost, endpoint, strings.NewReader(values.Encode()))
	//fmt.Println(strings.NewReader(params))
	//fmt.Println(req)
	if err != nil {
		return HttpResponseData{}, err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-API-KEY", scanSetting.AuthKey)
	req.Header.Add("User-Agent", createUserAgentString())

	client := &http.Client{}
	resp, err := client.Do(req)
	defer resp.Body.Close()

	return getResponseData(resp, err), err
}

func getResponseData(resp *http.Response, httpError error) HttpResponseData {
	body, _ := ioutil.ReadAll(resp.Body)

	return HttpResponseData{
		Status: resp.StatusCode,
		Body:   body,
		Error:  httpError,
	}
}

func GetApiServerName() string {
	api_server, ok := os.LookupEnv("VADDY_API_SERVER")
	if ok {
		//if api server does not contain protocol name, add https://
		var regex string = `^https://.*$`
		if regexp.MustCompile(regex).Match([]byte(api_server)) {
			return api_server
		}
		return "https://" + api_server
	}
	return config.API_SERVER
}

func createUserAgentString() string {
	return fmt.Sprintf("%s: %s (%s, %s)", config.USER_AGENT, config.VERSION, runtime.GOOS, runtime.GOARCH)
}
