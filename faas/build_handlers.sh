#!/bin/sh

if ! [ -d $1 ];
then
   echo "error: $1 is not a directory";
   exit;
fi

for f in $(find "$1" -type f -name "*.go");
do
    echo "building $f";
    GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/ $f
done