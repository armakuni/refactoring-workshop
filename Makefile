clean:
	@rm -rf *.json coverage.out

test:
	go test

coverage:
	go test -coverprofile coverage.out
	go tool cover -html=coverage.out
