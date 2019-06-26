#!/usr/bin/env sh

go get github.com/mitchellh/gox
gox -os="linux darwin windows" -arch="amd64" ${LDFLAGS} -output "dist/placeholder.{{.OS}}.{{.Arch}}" -verbose ./cmd
