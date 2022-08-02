# C10k-Go
This article is about the C10k-Go project.

# Prerequisites
## Go

## C10k

## linux (Ubuntu)
* netstat
* top
* vim 
* wget


## Install-Script
### install go-stress-testing
[go-stress-testing](https://github.com/link1st/go-stress-testing)

# Usage
just run the following command

## 1 start server 
```bash
cd cmd/server && GOMAXPROCS=4 go run server.go
```

## 2 start client
```bash
cd cmd && bash go-stress-testing.sh
```

## 3 check the result && watch the system load
```bash
netstat -tlnpu # check the port
top   # check the system load
```