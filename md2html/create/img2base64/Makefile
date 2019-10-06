GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=img2base64
BINARY_WIN=$(BINARY_NAME)_win
BINARY_MAC=$(BINARY_NAME)_mac
BINARY_LINUX=$(BINARY_NAME)_linux

RUN_OPTION=./readme.html ./build_on_win10.gif

RUN_COMMAND :=
ifeq ($(OS),Windows_NT)
	RUN_COMMAND=$(BINARY_WIN) $(RUN_OPTION)
else
	UNAME_S := $(shell uname -s)
	ifeq ($(UNAME_S),Linux)
		RUN_COMMAND=$(BINARY_LINUX) $(RUN_OPTION)
	endif
	ifeq ($(UNAME_S),Darwin)
		RUN_COMMAND=$(BINARY_MAC) $(RUN_OPTION)
	endif
endif


#all: build-linux build-mac build-windows make-symlink
all: build-windows build-mac build-linux

clean:
	$(GOCLEAN)
	rm -f $(BINARY_WIN)
	rm -f $(BINARY_MAC)
	rm -f $(BINARY_LINUX)

# クロスコンパイル
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WIN)
build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MAC)
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX)


.PHONY: run
run:
	./$(RUN_COMMAND)

