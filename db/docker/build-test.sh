#!/usr/bin/env bash

mkdir -p db/docker/test-data
rm -rf db/docker/test-data/*

export TEST=true

db/docker/up

db/docker/dbmate wait

db/docker/dbmate up

docker stop db-test

tar cf db/docker/test.tar -C db/docker/test-data .