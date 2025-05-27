all: build
# .DEFAULT_GOAL:=build
# variable

build: 
	@echo 'building app...'
	rm -rf ./fileserver
	rm -rf bin
	GOOS=linux GOARCH=amd64 go build -o bin/fileserver-linux main.go
	GOOS=windows GOARCH=amd64 go build -o bin/fileserver-windows.exe main.go
	GOOS=darwin GOARCH=amd64 go build -o bin/fileserver-darwin main.go

clean:
	@echo 'removing binary..'
	rm -rf ./fileserver

dev:
	@echo 'starting dev server'
	CompileDaemon -build="go build -o ./fileserver main.go" -command="./fileserver"