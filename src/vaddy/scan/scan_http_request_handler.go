package scan

import "vaddy/httpreq"

// scanパッケージ内で共有するHTTPリクエスト/レスポンスを扱うグローバル変数
var httpRequestHandler httpreq.HttpReqInterface = httpreq.HttpRequestData{}
