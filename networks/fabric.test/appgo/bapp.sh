#!/bin/bash
GOBIN=$PWD
rm app
go install app.go appclient.go server.go utils.go version.go
