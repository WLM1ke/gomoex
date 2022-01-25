new:
	brew install go
	brew install golangci-lint
update:
	brew upgrade go
	brew upgrade golangci-lint
lint:
	golangci-lint run --no-config --enable-all --fix --tests=true --exclude-use-default=false \
	--disable golint \
	--disable maligned \
	--disable scopelint \
    --disable interfacer \
    --disable testpackage \
    --disable paralleltest \
    --disable exhaustivestruct
test:
	go get -u -t -v ./...
	go mod tidy -v
	make lint
	go test -v -covermode=atomic -race ./...