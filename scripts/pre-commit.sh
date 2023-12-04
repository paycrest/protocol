#!/bin/sh

# Pre-commit hook for running linter and tests in a Go project

# Run linter (using a linter like golangci-lint)
echo "Running linter..."
if [ "$(uname)" = "Darwin" ]; then
    # Perform actions for macOS
    if ! brew ls --versions golangci-lint >/dev/null; then
        echo "golangci-lint is not installed. Attempting to install it"
        # Add your desired commands or actions here
        brew install golangci-lint || brew upgrade golangci-lint
    fi
    # Add your desired commands or actions here
    golangci-lint run ./...
else
    # Perform actions for other operating systems
    # Add your desired commands or actions here
    golangci-lint run ./...
fi

# Run tests
echo "Running tests..."
go test "$(go list ./... | grep -v /ent | grep -v /config | grep -v /database | grep -v /routers)" -coverprofile=coverage.out ./...

# If any errors occurred, abort the commit
if [ $? -ne 0 ]; then
    echo "Linter or Tests failed. Please fix the issues before committing."
    exit 1
fi

# If everything passed, allow the commit to proceed
exit 0
