
BIN=md2html_bin

build:
ifeq ($(shell uname -o),Msys)
	go build -o $(BIN).exe
else
	go build -o $(BIN)
endif

run: build
	./$(BIN) README.md

test: build
	go run test/test.go README.md

all: build


clean:
	rm $(BIN)*
	rm *.html

