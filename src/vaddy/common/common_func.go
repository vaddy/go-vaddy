package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"vaddy/httpreq"
)

// HTTPレスポンスのステータスコードを見てエラー判定
func CheckHttpResponse(response httpreq.HttpResponseData) error {
	status_code := response.Status
	//fmt.Println(status_code)
	if status_code != 200 {
		return errors.New("Network/Auth error\n" + string(response.Body))
	}

	if response.Error != nil {
		return response.Error
	}
	return nil
}

//Json文字列を任意の構造体にマッピングする
func ConvertJsonToStruct(jsonByteData []byte, structData interface{}) error {
	err := json.Unmarshal(jsonByteData, structData)
	if err != nil {
		return err
	}
	return nil
}

// 検査開始直後はポーリング間隔を短くする
// 検査がすぐに終わるケースの方が多いため
func GetSleepSec(waitCount int) int {
	var sleepSec int = 20
	if waitCount < 10 {
		sleepSec = 5
	}
	return sleepSec
}

func PrintDots(count int ) {
	if count > 0 && (count%60 == 0) { //wrap every 60 dots.
		fmt.Println(".")
	} else {
		fmt.Print(".")
	}
}