
VERSION:=$(gobump show | gojq -r .version)

BIN=md2html
DST=../build
GOARCH=amd64

FLAGS_WIN=-ldflags="-X main.version=$(VERSION) -X main.revision=$(git rev-parse --short HEAD) -H windowsgui -w -s -extldflags '-static'"
FLAGS_UNIX=-ldflags="-X main.version=$(VERSION) -X main.revision=$(git rev-parse --short HEAD) -w -s -extldflags '-static'"
FLAGS=-a -tags netgo -installsuffix netgo

#.PHONY: build
.PHONY: all
build:
	bash update_version_info.sh
	mkdir -p $(DST)
ifeq ($(shell uname -o),Msys)
	CGO_ENABLED=0 go build -o $(DST)/$(BIN).exe
else
	CGO_ENABLED=0 go build -o $(DST)/$(BIN)
endif
	cp -rf $(DST)/$(BIN) ..

run: build
	cd .. && ./$(BIN) readme.md

#all: build
all: release-build

release-build:
	bash update_version_info.sh
	go generate
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=windows CGO_ENABLED=0 go build -o $(DST)/$(BIN)_windows.exe $(FLAGS_WIN) $(FLAGS)
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=darwin CGO_ENABLED=0 go build -o $(DST)/$(BIN)_macOS $(FLAGS_UNIX) $(FLAGS)
	CGO_ENABLED=0 GOARCH=$(GOARCH) GOOS=linux CGO_ENABLED=0 go build -o $(DST)/$(BIN)_linux $(FLAGS_UNIX) $(FLAGS)

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
	go get -u -v github.com/x-motemen/gobump/cmd/gobump
	go get -u -v github.com/itchyny/gojq/cmd/gojq
