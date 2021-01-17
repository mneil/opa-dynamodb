run:
	go build -o /go/bin/opa
	opa run \
		--server \
		--log-level=debug \
		--addr=0.0.0.0:8001 \
		--authentication=token \
		--authorization=basic \
		--set=services.opa.credentials=null

serve-docs:
	# Requires nodejs and docsify
	docsify serve ./docs

unit:
	go test -v -short ./...

integration:
	go test -run Integration ./...

test:
	# Run full test suite
	go test -v -cover -coverprofile cover.out ./...
	go tool cover -func cover.out

cover: test
	go tool cover -html=cover.out -o cover.html

.PHONY: run unit integration
