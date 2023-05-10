#!/bin/sh

echo "============ Build and run database initialization for the baseline ============  ";
go run ../hotel_baseline/db/init/setup.go;

echo "============ Run nubes generator to initialize database for the nubes project ============ ";
export PATH=${PATH}:`go env GOPATH`/bin
go build -o ./../../generator/main ./../../generator/main.go 
./../../generator/main handlers -t=..//hotel//types -o=..//hotel// -m=github.com/Astenna/Nubes/evaluation/hotel -i=true -g=false;

echo "============ Run nubes generator to initialize database for the nubes project ============ ";
go run ../hotel_seeder/*.go