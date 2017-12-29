#!/bin/bash
GOBIN=$PWD

if [ -f "netapp" ]; then
 rm netapp
fi
go install netapp.go setup.go utils.go
