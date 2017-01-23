build:
	govendor generate
	go build -o build/dp-dd-search-api

debug: build
	HUMAN_LOG=1 ./build/dp-dd-search-api

.PHONY: build debug
