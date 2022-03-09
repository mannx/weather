#!/bin/sh

version=0.1.0
time=$(date)

# if we have a -t flag, simply output the version for use as a tag in the main build.sh script
if [ $# -ne 0 ]; then
	echo $version
	exit 0
fi

go build -ldflags="-X 'main.BuildTime=$time' -X 'main.BuildVersion=$version'" -o /weather .
