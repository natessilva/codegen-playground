#!/usr/bin/env bash
export TEST=true
export PGHOST=localhost
export PGDATABASE=codegen
export PGPORT=5433
export PGUSER=codegen
export PGPASSWORD=codegen
export SERVER_PORT=8001

function finish {
  kill $pid
  wait $pid
  docker stop db-test
}
trap finish EXIT

if ! [ -f db/docker/test.tar ]; then
  db/docker/build-test.sh
fi

rm -fr db/docker/test-data/* 2> /dev/null
tar xf db/docker/test.tar -C db/docker/test-data

db/docker/up
db/docker/dbmate wait

go build -o server cmd/server/main.go 
./server &
pid="$!"

npm test