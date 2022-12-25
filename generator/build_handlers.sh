#!/usr/bin/env bash

if ! [ -d $1 ];
then
   echo "error: $1 is not a directory";
   exit;
fi

for f in $(find "$1" -type f -name "*.go");
do
    echo "building $f";
    env GOOS=linux go build -o ../faas/bin/ $f
done