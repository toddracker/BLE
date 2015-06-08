#!/bin/bash

go get github.com/devicehive/gatt
GOOS=linux go build -o bin/x86_64/devicehive-ble ../../../devicehive-ble/devicehive-ble.go
GOOS=linux GOARCH=arm GOARM=7 go build -o bin/armhf/devicehive-ble ../../../devicehive-ble/devicehive-ble.go