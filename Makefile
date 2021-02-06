TARGET=ibex

all: clean build install

build:
	@go build -o build/ibex cmd/ibex/*.go

clean:
	@rm -rf $(TARGET)
	@rm -rf build

install:
	@go install ./cmd/ibex
