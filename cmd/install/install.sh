#!/bin/bash

# create a directory for the files to sit in
if [[ ! -e /usr/local/scripturetool ]]; then
    echo "Creating /usr/local/scripturetool"
    mkdir /usr/local/scripturetool
else
    echo "/usr/local/scripturetool exists"
fi

# download the libraries
if [[ ! -e /usr/local/scripturetool/lib ]]; then
    echo "Downloading libraries"
    curl https://studmane.com/files/untracked/lib.tar.gz | tar xz -C /usr/local/scripturetool
else
    echo "Libraries already installed"
fi

# get the executable
# TODO determine 64-bit or 32-bit
# TODO arm or amd?
if [[ ! -c /usr/local/scripturetool/scripturetool ]]; then
    echo "Downloading executable"

    # determine the os and architecture
    # download the correct executable
    # TODO what about raspberry pi?
    # if [[ "$OSTYPE" == "linux-gnu" ]]; then
        # then Linux

    # elif [[ "$OSTYPE" == "darwin"* ]]; then
        # then mac

    # elif [[ "$OSTYPE" == "msys" ]]; then
        # then wsl (double check "msys")

    # fi
fi

EXECUTABLE_NAME="sct"

# create a link
if [[ "$(which $EXECUTABLE_NAME)" == "" ]]; then
    echo "Creating symlink"
    ln -s /usr/local/scripturetool/scripturetool /usr/local/bin/$EXECUTABLE_NAME 
else
    echo "Symlink exists at $(which $EXECUTABLE_NAME)"
fi

echo "Installation complete!"