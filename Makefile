COV_FILE=cover.out
SRC=*.go
CMD=search

all: $(SRC)
	go build -o $(CMD) $(SRC)

linux:
	GOOS=linux GOARCH=amd64 make all

darwin:
	GOOS=darwin GOARCH=amd64 make all

vet:
	go vet ./...

test: vet
	go test ./...

test-cov: 
	go test -cover -o=$(COV_FILE)

clean:
	rm -rf $(CMD)

.PHONY: clean test test-cov vet
