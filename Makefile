run:
	go build -o /go/bin/opa
	opa run \
		--server \
		--log-level=debug \
		--addr=0.0.0.0:8001 \
		--authentication=token \
		--authorization=basic \
		--set=services.opa.credentials=null

unit:
	go test -cover -v -short ./...

integration:
	go test -run Integration ./...

test:
	go test -cover -coverprofile cover.out ./...
	go tool cover -func cover.out

.PHONY: run unit integration
