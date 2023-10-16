build:
	go build -o out/api -v
run: build
	./out/api
test:
	go test -v ./...