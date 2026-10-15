BINARY ?= sistema-solar


.PHONY: all build clean

all: build

build:
	go build -buildvcs=false -o $(BINARY) .

clean:
	rm -f $(BINARY)
