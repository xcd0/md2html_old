
BIN=build/md2html

build: ./build/*
	mkdir -p build
ifeq ($(shell uname -o),Msys)
	go build -o $(BIN).exe
else
	go build -o $(BIN)
endif

run: build
	./$(BIN) README.md

all: build

cross:
	GOARCH=amd64 GOOS=windows go build -o $(BIN)_win.exe
	GOARCH=amd64 GOOS=darwin go build -o $(BIN)_macOS
	GOARCH=amd64 GOOS=linux go build -o $(BIN)_linux


clean:
	rm -f build/$(BIN)
	rm -f *.html
	rm -f *.mini.css

