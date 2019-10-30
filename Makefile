
BIN=md2html
DST=build

build: ./build/*
	mkdir -p $(DST)
ifeq ($(shell uname -o),Msys)
	go build -o $(DST)/$(BIN).exe
else
	go build -o $(DST)/$(BIN)
endif

run: build
	./$(BIN) README.md

all: build

release:
	GOARCH=amd64 GOOS=windows go build -o $(BIN)_win.exe
	GOARCH=amd64 GOOS=darwin go build -o $(BIN)_macOS
	GOARCH=amd64 GOOS=linux go build -o $(BIN)_linux

	cd $(DST)
	mv $(BIN)_win.exe $(BIN).exe
	zip md2html_binary_windows.zip md2html.exe
	mv $(BIN)_macOS $(BIN)
	zip md2html_binary_macOS.zip md2html
	mv $(BIN)_linux $(BIN)
	zip md2html_binary_linux.zip md2html



clean:
	rm -f build/$(BIN)
	rm -f *.html
	rm -f *.mini.css

