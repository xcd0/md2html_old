GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=img2base64
BINARY_LINUX=$(BINARY_NAME)_linux
BINARY_MAC=$(BINARY_NAME)_mac
BINARY_WIN=$(BINARY_NAME)_win

#all: build-linux build-mac build-windows make-symlink
all: build-linux build-mac build-windows

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_LINUX)
	rm -f $(BINARY_MAC)
	rm -f $(BINARY_WIN)

# クロスコンパイル
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_LINUX)
build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MAC)
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WIN)

# build.shでOS毎のバイナリをたたきに行くように変更したので不要になった
# シンボリックリンクを張る
#make-symlink:
#ifeq (\$(OS),Windows_NT)
## for Windows
#	cd ..
#	ln -sf base64/base64_win
#	mv base64_win base64
#else
#UNAME = \${shell uname}
#
#ifeq (\$(UNAME),Linux)
## for Linux
#	cd ..
#	ln -sf base64/base64_linux
#	mv base64_linux base64
#endif
#
#ifeq (\$(UNAME),Darwin)
## for MacOSX
#	cd ..
#	ln -sf base64/base64_mac
#	mv base64_mac base64
#endif

