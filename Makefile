
BIN=md2html
DST=build

.PHONY: build
build:
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
	GOARCH=amd64 GOOS=windows go build -o $(DST)/$(BIN)_windows.exe
	GOARCH=amd64 GOOS=darwin go build -o $(DST)/$(BIN)_macOS
	GOARCH=amd64 GOOS=linux go build -o $(DST)/$(BIN)_linux

	cd $(DST) && \
	mv $(BIN)_windows.exe $(BIN).exe && \
	zip md2html_binary_windows.zip md2html.exe && \
	mv $(BIN).exe $(BIN)_windows.exe

	cd $(DST) && \
	mv $(BIN)_macOS $(BIN) && \
	zip md2html_binary_macOS.zip md2html && \
	mv $(BIN) $(BIN)_macOS

	cd $(DST) && \
	mv $(BIN)_linux $(BIN) && \
	zip md2html_binary_linux.zip md2html && \
	mv $(BIN) $(BIN)_linux



clean:
	rm -rf build
	rm -f *.html
	rm -f *.mini.css

