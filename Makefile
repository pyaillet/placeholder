VERSION=`git describe --tags`

LDFLAGS=-ldflags "-X main.version=${VERSION}"

build: go.mod test
	go build -o dist/placeholder cmd/main.go

go.mod:
	go mod init github.com/pyaillet/placeholder

test:
	go test github.com/pyaillet/placeholder/pkg/... -v -race -coverprofile=coverage.txt -covermode=atomic

watch:
	watch go test github.com/pyaillet/placeholder/pkg -v

clean:
	rm dist/*

release: go.mod test
	go build -o dist/placeholder ${LDFLAGS} cmd/main.go
