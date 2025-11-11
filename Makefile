# Display all package documentation with examples
doc:
	@go doc -all utils

# Run all example tests
test-examples:
	@echo "Running example tests..."
	@go test -v ./utils -run Example

# Run all tests
test:
	@go test -v ./...