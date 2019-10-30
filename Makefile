
BIN=md2html

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
	rm -f $(BIN)
	rm -f *.html

