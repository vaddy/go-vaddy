#!/bin/sh

cd ../src/vaddy/

go build  -o ../../bin/vaddy-macosx-64bit
#GOOS=linux GOARCH=386 go build  -o ../../bin/vaddy-linux-32bit
GOOS=linux GOARCH=amd64 go build  -o ../../bin/vaddy-linux-64bit
#GOOS=windows GOARCH=386 go build  -o ../../bin/vaddy-win-32bit.exe
GOOS=windows GOARCH=amd64 go build  -o ../../bin/vaddy-win-64bit.exe
#GOOS=freebsd GOARCH=386 go build  -o ../../bin/vaddy-freebsd-32bit
GOOS=freebsd GOARCH=amd64 go build  -o ../../bin/vaddy-freebsd-64bit
