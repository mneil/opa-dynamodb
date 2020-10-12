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

test: unit integration

.PHONY: run unit integration
