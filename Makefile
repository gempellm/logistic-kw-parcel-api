.PHONY: build
build:
	go build cmd/logistic-kw-parcel-api/main.go

.PHONY: test
test:
	go test -v ./...