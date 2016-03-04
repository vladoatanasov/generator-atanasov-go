buildit:
	go get && go build -o bin/main
clean:
	rm -fr bin/*
linux:
	go get && GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags "-w" -o bin/main_linux
buildclean: clean buildit
cleanbuild: clean buildit
test:
	go get && go test
