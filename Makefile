APP_NAME := netwatchdog

.PHONY: all build run clean

all: build

build:
	go build -o $(APP_NAME) .

run:
	go run .

clean:
	rm -f $(APP_NAME)