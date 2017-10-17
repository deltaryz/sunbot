#!/bin/sh
go get -u github.com/techniponi/sunbot
mv $GOPATH/src/github.com/techniponi/sunbot/img ./img
sunbot