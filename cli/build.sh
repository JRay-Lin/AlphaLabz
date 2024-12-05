#!/bin/bash

# Create builds directory if it doesn't exist
mkdir -p builds

# Build for Windows
echo "Building for Windows..."
GOOS=windows GOARCH=amd64 go build -o builds/elimt-cli.exe
if [ $? -eq 0 ]; then
    echo "✓ Windows build successful"
else
    echo "× Windows build failed"
fi

# Build for macOS
echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o builds/elimt-cli-mac
if [ $? -eq 0 ]; then
    chmod +x builds/elimt-cli-mac
    echo "✓ macOS build successful"
else
    echo "× macOS build failed"
fi

# Build for Linux
echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o builds/elimt-cli-linux
if [ $? -eq 0 ]; then
    chmod +x builds/elimt-cli-linux
    echo "✓ Linux build successful"
else
    echo "× Linux build failed"
fi

echo "Adding execute permissions to builds..."
# Add execute permissions to all builds (including Windows .exe for consistency)
chmod +x builds/*