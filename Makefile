lint:
	clear
	golangci-lint -j 1 -c ./.golangci.yaml run ./...

fmt:
	gofumpt -l -w .
	go mod tidy
	goimports -w -local=github.com/icdb37/kypd .

test:
	go test ./...

build-debug:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -gcflags='all=-N -l' -v -o kypd

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -ldflags "-w -s" -v -o kypd

