all:
	go build -o main

test:
	go test ./...

clean:
	rm -f main
