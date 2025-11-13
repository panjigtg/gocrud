.PHONY: run docs build

run:
	go run main.go

docs:
	swag init -g main.go -o ./docs

build:
	go build -o app main.go
