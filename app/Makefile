# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Clean up
clean:
	rm -f coverage.out
	rm -rf image_cache/
	rm -rf managers/cache/image_cache
	rm -rf managers/cache/test_cache