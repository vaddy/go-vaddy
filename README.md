
go-vaddy: VAddy API command tool
=================================

VAddy API command tool using golang  
http://vaddy.net

Go-vaddy can start scan and check the result.

## OS type

You can use exe files on `go-vaddy/bin` directory.
If you use linux(64bit), use vaddy-linux-64bit.  

For example, `./vaddy-linux-64bit  api_key userID FQDN`

| OS            | file               | 
| ------------- |:------------------:| 
| Linux(64bit)  | vaddy-linux-64bit  |
| MacOS(64bit)  | vaddy-macosx64     |
| Windows(64bit)| vaddy-win-64bit.exe|
| Linux(32bit)  | vaddy-linux-32bit  |
| Windows(32bit)| vaddy-win-32bit.exe|



## Usage (start scan and get the result)

### Exit status
Go-vaddy returns 0 (no errors, no vulnerabilities) or 1 (errors, 1 or more vulnerabilities).


###Arguments

Usage: `vaddy-linux-64bit auth_key username(LoginID)  hostname crawl_id(optional)`

    ./vaddy-linux-64bit 123455667789  ichikaway  www.examplevaddy.net 30

CrawlID is optional. If you don't specify it, VAddy uses the latest crawl data.


###Argument with the crawl label keyword

Usage: `vaddy-linux-64bit auth_key username(LoginID)  hostname crawl_label_keyword`

    vaddy-linux-64bit 123455667789  ichikaway  www.examplevaddy.net useredit

This example set "useredit" keyword on crawl label parameter.  
At first, go-vaddy search by "useredit" keyword and get the Crawl ID .  
Then go-vaddy sets the Crawl ID to start scan.
If it can not get Crawl ID, vaddy uses latest Crawl ID.



### ENV
You can also set paramters using OS environment variables.  

    export VADDY_TOKEN="123455667789"  
    export VADDY_USER="ichikaway"  
    export VADDY_HOST="www.examplevaddy.com"  
    export VADDY_CRAWL="30"  


`VADDY_CRAWL` is optional. If you don't specify it, VAddy uses the latest crawl data.  
You can specify crawl label keyword on `VADDY_CRAWL` like this  

    export VADDY_CRAWL="search result pages"  
