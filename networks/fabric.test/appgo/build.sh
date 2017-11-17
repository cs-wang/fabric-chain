#!/bin/bash
GOBIN=$PWD
rm netapp
go install netapp.go setup.go utils.go
