
VAddy PrivateNet Command-Line Tool
======================================

VAddy PrivateNet Command-Line Tool.
[https://vaddy.net](https://vaddy.net)

This tool allows you to scan for vulnerabilities on VMs, your local development machine, servers on your intranet, and other environments that are not exposed to the public Internet.


![screen](../images/screen.png "screen")


## How It Works
The command-line tool creates an SSH tunnel to use for remote port forwarding. This will expose a port on your local web server to VAddy’s servers through the SSH tunnel, allowing VAddy to scan it for vulnerabilities.  
Though this does expose your local web server’s port outside of your intranet, access to that port is restricted to the VAddy servers at the other end of the SSH tunnel.

![screen](../images/privatenet.png "privatenet")

In the figure above, the VAddy PrivateNet command-line tool has connected local port 443 on a server on the customer’s intranet to a dedicated remote port (3210) for that customer on VAddy’s SSH server. When the customer initiates a scan, VAddy will send its requests through port 3210.  
Note that the VAddy PrivateNet command-line tool creates an SSH tunnel with a remote port (like 3210 above) that is automatically assigned to each server based on its domain name.

## Supported Environments
This tool uses the `ssh`, `ssh-keygen`, and `ps` commands; it also requires a Java Runtime Environment. As a result, the tool is only supported on macOS and Linux.  
Any machine running this tool must also be able to access pfd.vaddy.net on port 22 (SSH).


## Downloading the Tool
You can download a ZIP file of the project from [https://github.com/vaddy/go-vaddy/releases](https://github.com/vaddy/go-vaddy/releases) or clone the Git repository:

    $ git clone https://github.com/vaddy/go-vaddy.git


## Configuration
First, [generate a WebAPI key from VAddy’s admin console](https://console.vaddy.net/user/webapi).  

Next, open the project's `privatenet/conf` subdirectory and copy `vaddy.conf.example` to a new file named `vaddy.conf`. Edit all of the relevant settings, as described below.  
`vaddy.conf` uses environment variables to set the values required to initiate a scan. If you’ve already set any of these environment variables elsewhere in your command prompt or using a continuous integration service like CircleCI or Travis CI, you can comment them out in this configuration file.

### vaddy.conf

| Environment Variable    | Description                                    |
| ----------------------- |:----------------------------------------------:|
| VADDY_AUTH_KEY          | WebAPI key generated by VAddy’s admin console. |
| VADDY_FQDN              | Server name (FQDN) registered with VAddy. (e.g. www.example.com) |
| VADDY_VERIFICATION_CODE | Verification code generated by VAddy’s admin console when registering a server. |
| VADDY_USER              | Username for authenticating with the VAddy service. |
| VADDY_YOUR_LOCAL_IP     | The IP address to scan. This could be localhost or the address of a server on your intranet. (e.g. 172.16.1.10) |
| VADDY_YOUR_LOCAL_PORT   | The port to scan. (e.g. 80 or 443) |
| VADDY_CRAWL             | Optional: ID or label specifying which crawl data to use. |
| VADDY_HTTPS_PROXY       | Optional: Connect to the VAddy WebAPI server through user's proxy server.  Set "IP:Port". ex. "127.0.0.1:8443" |

You can only specify a single port with `VADDY_YOUR_LOCAL_PORT`. As a result, your application will be scanned over either an HTTP or HTTPS connection depending on the port number you have specified for `VADDY_YOUR_LOCAL_PORT`.


## Exit Status Codes
If no errors or vulnerabilities are found, this tool returns the same exit code ordinarily used for successful commands: 0. If any errors or vulnerabilities are found, however, this tool returns an exit code of 1.

## Usage
### Arguments

    Usage: ./vaddy_privatenet.sh action [-crawl crawl_id or crawl_label]

| action     |                                                           |
| ---------- |:---------------------------------------------------------:|
| connect    | Opens an SSH tunnel to VAddy’s server.                    |
| disconnect | Closes the SSH tunnel.                                    |
| scan       | Opens an SSH tunnel to VAddy’s server, scans for vulnerabilities, and then closes the tunnel. |


After opening an SSH tunnel in the background, the `connect` action checks the connection through VAddy’s server and then returns to the command prompt with a successful exit status if it did not encounter any problems.

You can use the `ps` command to confirm that the SSH process is running:

     $ ps aux | grep 'ssh -i vaddy/ssh/'

#### Example #1 (Connecting to VAddy’s Server)
    $ ./vaddy_privatenet.sh connect

#### Example #2 (Disconnecting to VAddy’s Server)
    $ ./vaddy_privatenet.sh disconnect


#### Example #3 (Starting a Scan)
The `scan` action automatically opens an SSH tunnel before it starts scanning and automatically closes the tunnel after it stops scanning.

    $ ./vaddy_privatenet.sh scan

#### Example #4 (Scanning with a Crawl ID of 1234)
    $ ./vaddy_privatenet.sh scan -crawl 1234

The `crawl` option is not required; if you omit it, VAddy will scan for vulnerabilities using the latest crawl ID data.
