#!/bin/bash

# Simple Go Installation Script for Ubuntu 22.04
#
# This script installs the Go programming language using the recommended 'snap' method.
# It also verifies that the installation was successful.

# --- Style Functions for Output ---
# Adds a little color to the output to make it easier to read.
print_info() {
    echo -e "\n\e[34m\e[1m[INFO]\e[0m $1"
}

print_success() {
    echo -e "\e[32m\e[1m[SUCCESS]\e[0m $1"
}

print_error() {
    echo -e "\e[31m\e[1m[ERROR]\e[0m $1"
}

# --- Installation ---

# Exit immediately if a command exits with a non-zero status.
set -e

print_info "Starting Go installation..."

# Step 1: Install Go using snap.
# We use 'sudo' because installing software requires administrator privileges.
# The '--classic' flag is necessary to allow the Go snap to access system
# resources, which is required for compilers and development tools.
print_info "Installing Go via snap. This may take a moment..."
sudo snap install go --classic

# Step 2: Verify the installation.
# This step checks that the 'go' command is now available and prints its version.
# If the installation failed, the script would have exited on the previous line.
print_info "Verifying Go installation..."

# Use 'command -v' to check if 'go' is in the system's PATH.
if command -v go &> /dev/null
then
    # Get the installed Go version.
    GO_VERSION=$(go version)
    print_success "Go has been installed successfully!"
    print_success "Version details: $GO_VERSION"
    print_info "You can now use the 'go' command in your terminal."
else
    print_error "Go installation failed. The 'go' command could not be found."
    exit 1
fi

echo "" # Add a final newline for clean output.