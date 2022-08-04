#!/bin/bash

#go-stress-testing-linux -c 100 -n 100000 -u http://127.0.0.1:8087
go-stress-testing-linux -c 100 -n 100000 -u http://localhost:8087/product?code=Sticker
#go-stress-testing-linux -c 100 -n 100000 -u http://127.0.0.1:8087/order