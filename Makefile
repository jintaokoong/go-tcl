all: build run

build:
	go build -o bin/tcl tcl.go
build-arm:
	GOOS=linux GOARCH=arm go build -o bin/tcl tcl.go
run:
	./bin/tcl
clean:
	rm -r bin/
	rm *.log