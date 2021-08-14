build:
	go build -o bin/tcl tcl.go
build-arm:
	GOOS=linux GOARCH=arm go build bin/tcl-arm tcl.go
run:
	go run tcl.go
clean:
	rm -r bin/
	rm *.log