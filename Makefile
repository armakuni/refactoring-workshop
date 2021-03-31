all: test complexity

setup:
	go get
	go get -t
	go get github.com/uudashr/gocognit/cmd/gocognit

test:
	go test

coverage:
	go test -coverprofile coverage.out
	go tool cover -html=coverage.out

complexity:
	gocognit transformer.go

clean:
	@rm -rf coverage.out

