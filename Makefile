default: test

test:
	go test -v -race ./...

run:
	go run main.go