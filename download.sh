#!/bin/bash

# Function to detect os and architecture
function detect_architecture() {
    if [[ "$(uname)" == "Darwin" ]]; then
        # macOS
        arch=$(uname -m)
        case $arch in
            x86_64)
                echo "darwin_amd64"
                ;;
            arm64)
                echo "darwin_arm64"
                ;;
            i386)
                echo ""
                ;;
            *)
                echo ""
                ;;
        esac
    elif [[ "$(uname)" == "Linux" ]]; then
        # Linux
        arch=$(uname -m)
        case $arch in
            x86_64)
                echo "linux_amd64"
                ;;
            i386|i686)
                echo "linux_386"
                ;;
            armv7l)
                echo ""
                ;;
            aarch64)
                echo "linux_arm64"
                ;;
            *)
                echo ""
                ;;
        esac
    elif [[ "$OSTYPE" == "cygwin" || "$OSTYPE" == "msys" ]]; then
        # Cygwin or MSYS (Windows)
        arch=$(uname -m)
        case $arch in
            x86_64)
                echo "windows_amd64"
                ;;
            i686)
                echo ""
                ;;
            *)
                echo ""
                ;;
        esac
    elif command -v sw_vers > /dev/null 2>&1; then
        # Check for macOS using sw_vers
        arch=$(uname -m)
        case $arch in
            x86_64)
                echo "darwin_amd64"
                ;;
            i386)
                echo ""
                ;;
            *)
                echo ""
                ;;
        esac
    else
        echo ""
    fi
}

GITHUB_API="https://api.github.com/repos/teleology-io/foundation-cli/releases/latest"
OS_FILE=$(detect_architecture)

if [ -z "$OS_FILE" ]; then
    echo "foundation-cli is not supported on this architecture"
fi

release_info=$(curl -s "$GITHUB_API")

# Check if curl was successful
if [ $? -ne 0 ]; then
    echo "Error fetching release info"
    exit 1
fi

# Extract the tar.gz URL for Linux (amd64). Change according to your needs.
tarball_url=$(echo "$release_info" | grep "browser_download_url" | grep $OS_FILE | cut -d '"' -f 4)

# Check if URL was found
if [ -z "$tarball_url" ]; then
    echo "Tarball URL not found"
    exit 1
fi

# Download the tar.gz file
curl -L -o foundation.tar.gz "$tarball_url" 

# Check if curl was successful
if [ $? -ne 0 ]; then
    echo "Error downloading the file"
    exit 1
fi

# extract cli
tar -xvzf foundation.tar.gz

# cleanup
rm foundation.tar.gz