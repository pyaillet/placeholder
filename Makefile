build: go.mod test
	go build -o dist/placeholder cmd/main.go

go.mod:
	go mod init github.com/pyaillet/placeholder

test:
	go test github.com/pyaillet/placeholder/pkg/... -v

watch:
	watch go test github.com/pyaillet/placeholder/pkg -v

clean:
	rm dist/*

