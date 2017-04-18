
VAddy for Private Net command tool
======================================

VAddy for Private Net command tool.  
[https://vaddy.net](https://vaddy.net)

You can check web application security for your local web server(local PC, VM and Intranet servers).


![screen](../images/screen.png "screen")


## System requirements 

This tool needs Java, shell and Mac/Linux, and uses ssh, ssh-keygen, ps commands.



## Download

Download zip from [https://github.com/vaddy/go-vaddy/releases](https://github.com/vaddy/go-vaddy/releases),  
or git clone.

    git clone https://github.com/vaddy/go-vaddy.git



## Set up

Get your VAddy WebAPI key on [https://console.vaddy.net/user/webapi](https://console.vaddy.net/user/webapi)  

Create `privatenet/conf/vaddy.conf` file from `privatenet/conf/vaddy.conf.example` and edit it.  


## How to work

This tool creates ssh tunnel(remort port forwarding) to open your local web server port to us.  
VAddy can access your local web server through ssh tunnel.  
Don't worry, your local server opens for only VAddy.  


## Exit status
This tool returns 0 (no errors, no vulnerabilities) or 1 (errors, 1 or more vulnerabilities).


## Usage 

###Arguments

    Usage: ./vaddy_privatenet.sh action [-crawl crawl_id or crawl_label]   


| action        |                                                           | 
| ------------- |:---------------------------------------------------------:| 
| connect       | connect to VAddy ssh server with remote port forwarding.  |
| disconnect    | cut off ssh tunnel connection.                            |
| scan          | connect, scan and disconnect.                             |


#### Example1 (make connection)

    ./vaddy_privatenet.sh connect

#### Example2 (start scan)
With scan action, this tool connects and disconnects VAddy server automatically. 

     ./vaddy_privatenet.sh scan


#### Example3 (start scan with crawlID 1234)
     ./vaddy_privatenet.sh scan -crawl 1234


CrawlID is optional. If you don't specify it, VAddy uses the latest crawl data.


