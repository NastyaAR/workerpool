.PHONY: build run test clean

build:
	go build -o workerpool .

run: build
	./workerpool

test:
	go test ./...

clean:
	go clean
	rm -f workerpool