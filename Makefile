VERSION=1.0

build:
	go build -o url-checker cmd/url-checker/main.go

test:
	go test -v ./internal/... -coverprofile=coverage.out

coverage:
	go tool cover -html=coverage.out

build-image:
	docker build -t url-checker:$(VERSION) .