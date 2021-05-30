hello:
	echo hello

build-server:
	go build -o bin/cchk-server github.com/johnllao/remoteproc/creditcheck/server

build-client:
	go build -o bin/cchk-cli github.com/johnllao/remoteproc/creditcheck/cli

build-html:
	go build -o bin/cchk-html github.com/johnllao/remoteproc/creditcheck/html

build: build-server build-client build-html