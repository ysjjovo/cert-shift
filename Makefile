run:
	go run main.go
build:
	GOOS=linux go build main.go
	zip function.zip main