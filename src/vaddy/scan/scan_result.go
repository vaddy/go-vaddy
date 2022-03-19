package scan

type ScanResult struct {
	Status        string `json:"status"`
	AlertCount    int    `json:"alert_count"`
	ScanCount     int    `json:"scan_count"`
	ScanResultUrl string `json:"scan_result_url"`
	Complete      int    `json:"complete"`
	CrawlId       int    `json:"crawl_id"`
	CrawlLabel    string `json:"crawl_label"`
}

func (s ScanResult) IsIncomplete() bool {
	return s.AlertCount == 0 && s.Complete < 100
}

func (s ScanResult) IsScanRunning() bool {
	if s.Status == "scanning" {
		return true
	}
	return false
}

func (s ScanResult) IsCanceled() bool {
	if s.Status == "canceled" {
		return true
	}
	return false
}

func (s ScanResult) IsFinished() bool {
	if s.Status == "finish" {
		return true
	}
	return false
}

func (s ScanResult) IsVulnExist() bool {
	if s.AlertCount > 0 {
		return true
	}
	return false
}
