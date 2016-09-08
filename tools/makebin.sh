#!/bin/sh

cd ../src/vaddy/

go build  -o ../../bin/vaddy-macosx64
GOOS=linux GOARCH=386 go build  -o ../../bin/vaddy-linux-32bit
GOOS=linux GOARCH=amd64 go build  -o ../../bin/vaddy-linux-64bit
GOOS=windows GOARCH=386 go build  -o ../../bin/vaddy-win-32bit
GOOS=windows GOARCH=amd64 go build  -o ../../bin/vaddy-win-64bit

