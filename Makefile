.PHONY: build
build:
	@go build -o ./bin/ github.com/LilithGames/nrpc/protoc-gen-nrpc/...

.PHONY: clean
clean:
	@rm -f bin/*
	@rm -f example/bin/*

.PHONY: proto
proto: clean
	@protoc -I=. --go_out=paths=source_relative:. proto/nrpc.proto

.PHONY: example
example: proto
	@protoc -I=. --go_out=paths=source_relative:. --nrpc_out=paths=source_relative:. example/proto/test.proto

.PHONY: tag
tag:
	@git tag $$(svu next)
