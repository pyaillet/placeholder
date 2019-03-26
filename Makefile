VERSION=`git describe --tags`

LDFLAGS=-ldflags "-X main.version=${VERSION}"

build: go.mod test
	if [ -z ${VERSION} ]; then go build -o dist/placeholder cmd/main.go; fi
	if [ ! -z ${VERSION} ]; then gox -os="linux darwin windows" -arch="amd64" ${LDFLAGS} -output "dist/placeholder.{{.OS}}.{{.Arch}}" -verbose ./...; fi

go.mod:
	go mod init github.com/pyaillet/placeholder

test:
	go test github.com/pyaillet/placeholder/pkg/... -v -race -coverprofile=coverage.txt -covermode=atomic

watch:
	watch go test github.com/pyaillet/placeholder/pkg -v

clean:
	rm dist/*

release: go.mod test
	gox -os="linux darwin windows" -arch="amd64" ${LDFLAGS} -output "dist/placeholder.{{.OS}}.{{.Arch}}" -verbose ./...
