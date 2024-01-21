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
