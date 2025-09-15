#!/bin/bash

# Web4 API Bridge Demo - Build Script for Linux/macOS

set -e  # Exit on any error

echo "========================================"
echo "Web4 API Bridge Demo - Build Script"
echo "========================================"

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
echo "Checking prerequisites..."

if ! command_exists cmake; then
    echo "ERROR: CMake not found. Please install CMake."
    echo "Ubuntu/Debian: sudo apt-get install cmake"
    echo "macOS: brew install cmake"
    exit 1
fi

if ! command_exists git; then
    echo "ERROR: Git not found. Please install Git."
    echo "Ubuntu/Debian: sudo apt-get install git"
    echo "macOS: brew install git"
    exit 1
fi

# Detect OS
if [[ "$OSTYPE" == "linux-gnu"* ]]; then
    OS="linux"
    echo "Detected OS: Linux"
elif [[ "$OSTYPE" == "darwin"* ]]; then
    OS="macos"
    echo "Detected OS: macOS"
else
    echo "WARNING: Unknown OS type: $OSTYPE"
    OS="unknown"
fi

# Install system dependencies
echo "Installing system dependencies..."

if [[ "$OS" == "linux" ]]; then
    # Ubuntu/Debian
    if command_exists apt-get; then
        sudo apt-get update
        sudo apt-get install -y build-essential cmake git libssl-dev pkg-config
    # CentOS/RHEL
    elif command_exists yum; then
        sudo yum groupinstall -y "Development Tools"
        sudo yum install -y cmake git openssl-devel pkgconfig
    # Fedora
    elif command_exists dnf; then
        sudo dnf groupinstall -y "Development Tools"
        sudo dnf install -y cmake git openssl-devel pkgconfig
    else
        echo "WARNING: Could not detect package manager. Please install build tools manually."
    fi
elif [[ "$OS" == "macos" ]]; then
    if command_exists brew; then
        brew install cmake openssl pkg-config
    else
        echo "WARNING: Homebrew not found. Please install build tools manually."
    fi
fi

# Set build directory
BUILD_DIR="build"
VCPKG_DIR="$HOME/vcpkg"

# Install vcpkg if not exists
if [[ ! -d "$VCPKG_DIR" ]]; then
    echo "Installing vcpkg..."
    git clone https://github.com/Microsoft/vcpkg.git "$VCPKG_DIR"
    cd "$VCPKG_DIR"
    ./bootstrap-vcpkg.sh
    ./vcpkg integrate install
    cd ..
fi

# Install dependencies
echo "Installing dependencies..."
cd "$VCPKG_DIR"

# Detect architecture
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
    TRIPLET="x64-linux"
elif [[ "$ARCH" == "aarch64" ]] || [[ "$ARCH" == "arm64" ]]; then
    TRIPLET="arm64-linux"
else
    echo "WARNING: Unknown architecture $ARCH, using x64-linux"
    TRIPLET="x64-linux"
fi

./vcpkg install nlohmann-json:$TRIPLET
./vcpkg install grpc:$TRIPLET
./vcpkg install protobuf:$TRIPLET

cd ..

# Create build directory
if [[ ! -d "$BUILD_DIR" ]]; then
    mkdir "$BUILD_DIR"
fi
cd "$BUILD_DIR"

# Configure with CMake
echo "Configuring with CMake..."
cmake .. -DCMAKE_TOOLCHAIN_FILE="$VCPKG_DIR/scripts/buildsystems/vcpkg.cmake" -DCMAKE_BUILD_TYPE=Release

if [[ $? -ne 0 ]]; then
    echo "ERROR: CMake configuration failed."
    exit 1
fi

# Build the project
echo "Building project..."
make -j$(nproc 2>/dev/null || sysctl -n hw.ncpu 2>/dev/null || echo 4)

if [[ $? -ne 0 ]]; then
    echo "ERROR: Build failed."
    exit 1
fi

# Copy configuration file
if [[ ! -f "bin/config.json" ]]; then
    cp "../config.json" "bin/"
fi

echo "========================================"
echo "Build completed successfully!"
echo "========================================"
echo
echo "Executable location: $BUILD_DIR/bin/APIBridgeDemo"
echo
echo "To run the demo:"
echo "  cd $BUILD_DIR/bin"
echo "  ./APIBridgeDemo"
echo
echo "To run with custom config:"
echo "  ./APIBridgeDemo --config config.json"
echo

# Ask if user wants to run the demo
read -p "Do you want to run the demo now? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Starting demo..."
    cd bin
    ./APIBridgeDemo
fi 