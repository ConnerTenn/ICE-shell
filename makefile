
all: clean build run

list:
	@echo Targets:
	@echo "  build"
	# @echo "  test"
	@echo "  clean"
	@echo default: clean build shell

run: build
	./iceShell

build: clean
	go build

test: build
	go test ./ -v -failfast -timeout 5s

clean:
	rm -f ice  .log  log  dump
