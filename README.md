
go-vaddy: VAddy API Command-Line Tool
=================================

VAddy API Command-Line Tool using golang  
https://vaddy.net

Go-vaddy can start scan and check the result.

日本語ドキュメントはこちらです。
https://github.com/vaddy/go-vaddy/blob/master/README_ja.md


## OS type

You can use exe files on `go-vaddy/bin` directory.
If you use linux(64bit), use vaddy-linux-64bit.  

For example, `./vaddy-linux-64bit  api_key userID FQDN`

| OS                  | file               |
| ------------------- | ------------------ |
| Linux(64bit)        | vaddy-linux-64bit  |
| MacOS(64bit Intel)  | vaddy-macosx-64bit |
| Windows(64bit)      | vaddy-win-64bit.exe|
| FreeBSD(64bit)      | vaddy-freebsd-64bit|




## Usage (start scan and get the result)

### Exit status
Go-vaddy returns 0 (no errors, no vulnerabilities) or 1 (errors, 1 or more vulnerabilities).




### ENV
You can check V1/V2 project on the dashboard screen after login.

#### for V1 Project

    export VADDY_TOKEN="123455667789"  
    export VADDY_USER="ichikaway"  
    export VADDY_HOST="www.examplevaddy.com"  
    #export VADDY_CRAWL="30"
    #export VADDY_SCAN_TYPE="SQLI,XSS,..."

#### for V2 Project

    export VADDY_TOKEN="123455667789"
    export VADDY_USER="ichikaway"
    export VADDY_PROJECT_ID="your project id"
    #export VADDY_CRAWL="30"
    #export VADDY_SCAN_TYPE="SQLI,XSS"

* `VADDY_USER` is VAddy login ID.

* `VADDY_CRAWL` is optional. If you don't specify it, VAddy uses the latest crawl data.  
You can specify crawl label keyword on `VADDY_CRAWL` like this  

    export VADDY_CRAWL="search result pages"  

* `VADDY_SCAN_TYPE` is optional to specify a specific scan type. [Scan type list document](https://github.com/vaddy/WebAPI-document/blob/master/VAddy-WebApi-ScanType.md)
Without this option, all scan will be performed. If you specify an item that does not exist or an item that does not exist in your plan, the error `Invalid scan type selected` will be returned.



### Command Execution

    cd bin
    ./vaddy-linux-64bit


### Slack Integration
Setting these OS environment variables,
Post message to the slack when VAddy found vulnerabilities.  

    export SLACK_WEBHOOK_URL="webhook url"
    export SLACK_USERNAME="your user (optional)"
    export SLACK_CHANNEL="your channel (optional)"
    export SLACK_ICON_EMOJI=":smile: (optional)"
    export SLACK_ICON_URL="icon url (optional)"
