
BIN=md2html
DST=../build
GOARCH=amd64
FLAGS_WIN=-ldflags='-H windowsgui -w -s -extldflags "-static"' -a -tags netgo -installsuffix netgo
FLAGS=-ldflags='-w -s -extldflags "-static"' -a -tags netgo -installsuffix netgo

#.PHONY: build
.PHONY: all
build:
	mkdir -p $(DST)
ifeq ($(shell uname -o),Msys)
	go build -o $(DST)/$(BIN).exe
else
	go build -o $(DST)/$(BIN)
endif
	cp -rf $(DST)/$(BIN) ..

run: build
	cd .. && ./$(BIN) readme.md

#all: build
all: release-build

release-build:
	go generate
	GOARCH=$(GOARCH) GOOS=windows go build -o $(DST)/$(BIN)_windows.exe $(FLAGS_WIN)
	GOARCH=$(GOARCH) GOOS=darwin go build -o $(DST)/$(BIN)_macOS $(FLAGS)
	GOARCH=$(GOARCH) GOOS=linux go build -o $(DST)/$(BIN)_linux $(FLAGS)

release: release-build
	# for windows
	cd $(DST) && \
	mv $(BIN)_windows.exe $(BIN).exe && \
	zip md2html_binary_$(GOARCH)_windows.zip md2html.exe && \
	mv $(BIN).exe $(BIN)_windows.exe
	# for mac
	cd $(DST) && \
	mv $(BIN)_macOS $(BIN) && \
	zip md2html_binary_$(GOARCH)_macOS.zip md2html && \
	mv $(BIN) $(BIN)_macOS
	# for linux
	cd $(DST) && \
	mv $(BIN)_linux $(BIN) && \
	zip md2html_binary_$(GOARCH)_linux.zip md2html && \
	mv $(BIN) $(BIN)_linux
	# clean
	cd $(DST) && \
	rm $(DST)/$(BIN)_windows.exe && \
	rm $(DST)/$(BIN)_macOS && \
	rm $(DST)/$(BIN)_linux

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
	go get -u -v github.com/josephspurrier/goversioninfo/cmd/goversioninfo
