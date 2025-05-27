all: build
# .DEFAULT_GOAL:=build
# variable

build: 
	@echo 'building app...'
	rm -rf ./fileserver&&go build -o ./fileserver main.go

clean:
	@echo 'removing build dir..'
	rm -rf ./fileserver

dev:
	@echo 'starting dev server'
	CompileDaemon -build="go build -o ./fileserver main.go" -command="./fileserver"