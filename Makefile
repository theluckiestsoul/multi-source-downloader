install:
	@echo "Installing..."
	go install
	@echo "Done!"
.PHONY: install


build:
	@echo "Building..."
	go build -o multi-source-downloader
	@echo "Done!"
.PHONY: build

test:
	@echo "Testing..."
	go test -race -v ./...
	@echo "Done!"
.PHONY: test

clean:
	@echo "Cleaning..."
	go clean
	@echo "Done!"
