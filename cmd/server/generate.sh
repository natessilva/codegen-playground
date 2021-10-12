#!/usr/bin/env bash

cd ../..

go build -o gen cmd/gen/*.go

./gen -p api/app -n app -t tpl/server.go.plush > pkg/app/server_gen.go
./gen -p api/authn -n authn -t tpl/server.go.plush > pkg/authn/server_gen.go

./gen -p api/app -n app -t tpl/client.ts.plush > client/src/lib/app/client_gen.ts
./gen -p api/authn -n authn -t tpl/client.ts.plush > client/src/lib/authn/client_gen.ts

sqlc generate