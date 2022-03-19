package args

import (
	"errors"
	"fmt"
	"os"
	"regexp"
)

type ScanSetting struct {
	V2               bool
	AuthKey          string
	User             string
	Crawl            string
	ScanType         string
	ProjectId        string
	Fqdn             string
	VerificationCode string
	CrawlId          string
}

/**
 * v1 : auth_key, user, fqdn, crawl, verification_code, scan_type
 * v2 : auth_key, user, crawl, verification_code, scan_type, project_id
 */
func GetArgsFromEnv() (ScanSetting, error) {
	var scanSetting ScanSetting
	var ok1, ok2, ok3, ok4 bool

	scanSetting.VerificationCode, _ = os.LookupEnv("VADDY_VERIFICATION_CODE")
	scanSetting.AuthKey, ok1 = os.LookupEnv("VADDY_TOKEN")
	scanSetting.User, ok2 = os.LookupEnv("VADDY_USER")
	scanSetting.Fqdn, ok3 = os.LookupEnv("VADDY_HOST")
	scanSetting.ProjectId, ok4 = os.LookupEnv("VADDY_PROJECT_ID")
	scanSetting.Crawl, _ = os.LookupEnv("VADDY_CRAWL")
	scanSetting.CrawlId, _ = os.LookupEnv("VADDY_CRAWL")
	scanSetting.ScanType, _ = os.LookupEnv("VADDY_SCAN_TYPE")

	// v1
	if ok1 && ok2 && ok3 {
		scanSetting.V2 = false
		return scanSetting, nil
	}
	// v2
	if ok1 && ok2 && ok4 {
		scanSetting.V2 = true
		return scanSetting, nil
	}

	return scanSetting, errors.New("Missing system environments.")
}

func (s ScanSetting) IsV2() bool {
	return s.V2
}

//Crawlで指定された文字列が数字かどうか判定
//数字であればそのままcrawlIDとして使い、数字以外はcrawl labelとしてcrawl IDを引き上げる処理につなげるため
func (s ScanSetting) CheckNeedToGetCrawlId() bool {
	if len(s.Crawl) == 0 || s.Crawl == "" {
		return false
	}
	var regex string = `[^0-9]`
	return regexp.MustCompile(regex).Match([]byte(s.Crawl))
}

func (s ScanSetting) GetApiVersion() string {
	if s.V2 {
		return "v2"
	}
	return "v1"
}

//スキャン結果のjsonはv1でfqdn、V2でprojectIDがセットされるため切り替えて表示する文字列を返す
func (s ScanSetting) GetHostOrProjectName() string {
	if s.IsV2() {
		return "ProjectId: " + s.ProjectId
	}
	return "Server: " + s.Fqdn
}

//検査項目指定の検査を実行している場合は、指定した検査項目を表示する
func (s ScanSetting) PrintScanTypeSetting() {
	if len(s.ScanType) > 0 {
		fmt.Println("Scan Type: " + s.ScanType)
	}
}
