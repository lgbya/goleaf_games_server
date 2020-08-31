BINARY_NAME=server
MAIN_FILE=cmd/main.go
build:
	go build  -v -o $(BINARY_NAME) $(MAIN_FILE) 
	
start:	
	./$(BINARY_NAME)

rebuild:
	make build
	make start