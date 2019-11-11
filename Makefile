
BIN=md2html
DST=build
FLAGS=-ldflags='-w -s -extldflags "-static"' -a -tags netgo -installsuffix netgo

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

release-build:
	GOARCH=amd64 GOOS=windows go build -o $(DST)/$(BIN)_windows.exe $(FLAGS)
	GOARCH=amd64 GOOS=darwin go build -o $(DST)/$(BIN)_macOS $(FLAGS)
	GOARCH=amd64 GOOS=linux go build -o $(DST)/$(BIN)_linux $(FLAGS)

release: release-build
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

get:
	go get -u -v github.com/google/go-github/github
	go get -u -v github.com/russross/blackfriday
	go get -u -v github.com/shurcooL/github_flavored_markdown
	go get -u -v github.com/tdewolff/minify
	go get -u -v github.com/tdewolff/minify/css
	go get -u -v github.com/xcd0/go-nkf
