
all: shell

list:
	@echo Targets:
	@echo "  build"
	@echo "  test"
	@echo "  shell"
	@echo "  clean"
	@echo default: shell

shell: build
	./iceShell

build: test
	go build

test:
	go test ./ -v -failfast -timeout 5s

clean:
	rm -f ice  .log  log  dump
