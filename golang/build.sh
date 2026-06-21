#!/bin/bash

set -e

BINARY_NAME="cors-proxy"
BUILD_DIR="build"

echo "Building CORS Proxy..."

# Create build directory
mkdir -p "$BUILD_DIR"

# Build for current platform
go build -o "$BUILD_DIR/$BINARY_NAME" .

echo "Build successful: $BUILD_DIR/$BINARY_NAME"

# Cross-compile for common platforms
echo ""
echo "Cross-compiling for multiple platforms..."

GOOS=linux GOARCH=amd64 go build -o "$BUILD_DIR/${BINARY_NAME}-linux-amd64" .
GOOS=linux GOARCH=arm64 go build -o "$BUILD_DIR/${BINARY_NAME}-linux-arm64" .
GOOS=darwin GOARCH=amd64 go build -o "$BUILD_DIR/${BINARY_NAME}-darwin-amd64" .
GOOS=darwin GOARCH=arm64 go build -o "$BUILD_DIR/${BINARY_NAME}-darwin-arm64" .
GOOS=windows GOARCH=amd64 go build -o "$BUILD_DIR/${BINARY_NAME}-windows-amd64.exe" .

echo "Cross-compilation complete!"
echo ""
echo "Build artifacts:"
ls -lh "$BUILD_DIR/"
echo ""
echo "To run:"
echo "  ./build/$BINARY_NAME"
