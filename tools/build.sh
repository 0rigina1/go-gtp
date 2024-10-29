#!/bin/sh

CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o s11_msg_sender s11_msg_sender.go