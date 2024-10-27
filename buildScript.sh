


#!/bin/bash

FOLDER1=/home/mahdi/Desktop/go/socketProgram/ServerSocket
FOLDER2=/home/mahdi/Desktop/go/socketProgram/ClientSocket


 
# Function to build Go files with specific naming
build_executables() {
    local folder=$1
    local prefix=$2
    echo "Building executables in $folder with prefix '$prefix'..."

    # Navigate to the folder
    cd "$folder" || exit

    # Create a build directory if it doesn't exist

    # Build for Linux
    echo "Building for Linux..."
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o "/home/mahdi/Desktop/go/socketProgram/build/linux/${prefix}-linux"
    if [ $? -eq 0 ]; then
        echo "Linux build succeeded in $folder."
    else
        echo "Linux build failed in $folder."
    fi

    # Build for Windows
    echo "Building for Windows..."
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o "/home/mahdi/Desktop/go/socketProgram/build/windows/${prefix}-windows.exe"
    if [ $? -eq 0 ]; then
        echo "Windows build succeeded in $folder."
    else
        echo "Windows build failed in $folder."
    fi
}

# Build executables in Folder 1 with 'server' prefix
build_executables "$FOLDER1" "server"

# Build executables in Folder 2 with 'client' prefix
build_executables "$FOLDER2" "client"
