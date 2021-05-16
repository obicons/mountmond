GOSRC=$(shell find ./ -iname *.go)
TARGET=./bin/mountmond

$(TARGET): $(GOSRC)
	go build -o $(TARGET) ./cmd/mountmond/

.PHONY: build
build: $(TARGET)

.PHONY: test
test:
	go test -v ./...

.PHONY: clean
clean:
	rm -rf ./bin

.PHONY: install
install:
	./scripts/install.sh
