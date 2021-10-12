#!/usr/bin/env bash -e
export PGHOST=localhost
export PGDATABASE=codegen
export PGPORT=5432
export PGUSER=codegen
export PGPASSWORD=codegen
export SERVER_PORT=8000

go build -o server cmd/server/main.go 
./server