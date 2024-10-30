#!/bin/bash

FOLDER1=/home/mahdi/Desktop/go/socketProgram/server
FOLDER2=/home/mahdi/Desktop/go/socketProgram/client



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

build_executables "$FOLDER1" "server"
#
build_executables "$FOLDER2" "client"

move_executable(){
  local file_path="/home/mahdi/Desktop/go/socketProgram/build/linux"
  username=noori
  server_ip=192.168.100.97
  server_path=/home/noori
  password=noori@1
  echo "Moving executable files in $folder ...."

  sshpass -p "$password" scp -r "$file_path" "$username@$server_ip:$server_path"

  if [ $? -eq 0 ]; then
      echo "Folder uploaded successfully to $username@$server_ip:$server_path"
  else
      echo "Folder upload failed."
  fi
  }
move_executable
