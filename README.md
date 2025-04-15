# go-monitoring-htmx
Monitoring Dynamic Data from Server/PC and get also Weather/Crypto infromation real time with GO, HTMX and Websocket.. Data refreshs 5secs automatically...

![preview](./static/pic-selected-250316-0000-51.png)

## Code
[mozkaya1@github.com](https://github.com/mozkaya1/go-monitoring-htmx)

# Real Time Data
- Current Time
 
>Below feed is coming from my another go-api service.. It can be adjusted what you need with changing api url/query on internal/api/btc.go file. Details also on [go-api](https://github.com/mozkaya1/go-api#) repo page.
- Location
- Weather Temp
- Weather Description
- Bitcoin/USD
- Etherium/USD
  
>These system data from the pc/server which this api running on.
- System Specs
- Disk
- CPU
- Load
- System Temperature
- Docker Source(if any)

# Installing / Running Services 
Clone repository with submodule in it :
```bash
git clone --recurse-submodules git@github.com:mozkaya1/go-monitoring-htmx.git
```
And simply you can run it with make command.
> Docker system needs sudo rights to fetch docker status in the system
```bash
cd go-monitoring-htmx
sudo make run
```
**Output**
```bash
[musti@musti-deputyp25 go-monitoring-htmx]$ sudo make run
Server started on port 8080
Starting monitor server on port 8000
open browser at http://localhost:8000
Added subscriber &{0xc0000d0000}
```
Api Service running on `:8080` port and Monitoring Service is running on `:8000` port.. Just open http://localhost:8000 link on your browser and enjoy monitoring system every 5 secs...

