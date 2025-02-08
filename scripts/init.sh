#!/bin/sh

# pre-commit
echo "Checking pre-commit installation..."
if ! command -v pre-commit >/dev/null 2>&1; then
    echo "pre-commit is not installed. Installing..."

    if command -v pipx >/dev/null 2>&1; then
        pipx install pre-commit
    else
        echo "pipx is not installed. Installing..."

        if command -v brew >/dev/null 2>&1; then
            brew install pipx && pipx ensurepath
        elif command -v apt >/dev/null 2>&1; then
            sudo apt install -y pipx
        fi

        pipx install pre-commit
    fi
else
    echo "pre-commit is already installed."
fi

# mockgen
echo "Checking mockgen installation..."
if ! command -v mockgen >/dev/null 2>&1; then
    echo "mockgen is not installed. Installing..."
    go install github.com/golang/mock/mockgen@latest
else
    echo "mockgen is already installed."
fi

# go path
if ! echo "$PATH" | grep -q "$HOME/go/bin"; then
    echo "Adding ~/go/bin to PATH..."
    echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.bashrc
    echo 'export PATH=$HOME/go/bin:$PATH' >> ~/.zshrc
    export PATH=$HOME/go/bin:$PATH
    echo "Done! Please restart your terminal."
fi
