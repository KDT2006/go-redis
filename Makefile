run:build
	@./build/goredis

build:
	@go build -o build/goredis .

.PHONY: run build