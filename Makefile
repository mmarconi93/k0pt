.PHONY: build clean

build:
    go build -o kopt cmd/main.go

clean:
    rm -f kopt
