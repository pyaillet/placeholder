VERSION=`git describe --tags`

LDFLAGS=-ldflags "-X main.version=${VERSION}"

placeholder:
	go build -o dist/placeholder cmd/main.go

build: go.mod test
	if [ -z ${VERSION} ]; then go build -o dist/placeholder cmd/main.go; fi
	if [ ! -z ${VERSION} ]; then gox -os="linux darwin windows" -arch="amd64" ${LDFLAGS} -output "dist/placeholder.{{.OS}}.{{.Arch}}" -verbose ./cmd; fi

go.mod:
	go mod init github.com/pyaillet/placeholder

test:
	go test ./... -v -race -coverprofile=coverage.txt -covermode=atomic

watch:
	watch go test github.com/pyaillet/placeholder/pkg -v

clean:
	rm dist/*

release: go.mod test
	gox -os="linux darwin windows" -arch="amd64" ${LDFLAGS} -output "dist/placeholder.{{.OS}}.{{.Arch}}" -verbose ./...

.DEFAULT_GOAL := build
