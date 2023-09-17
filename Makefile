vet:
	@echo "Checking go files..."
	@echo "Starting golangci-lint..."
	@golangci-lint run ./...
	@echo "Checking go files... DONE"
	