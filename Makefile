all: test

test:
	@echo "test coap"
	@go test -v -coverprofile c.out coap
	@go tool cover -html=c.out -o /tmp/c.html
	@rm c.out
