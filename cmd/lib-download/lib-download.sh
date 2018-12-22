#!/bin/bash

# make sure we are in the right directory
if [[ "$PWD" != "$GOPATH/src/github.com/es3649/scripturetool" ]]; then
    pushd $GOPATH/src/github.com/es3649/scripturetool
fi

# check if the library is compiled
if [[ ! -e lib/ ]]; then
    echo "creating \`lib' directory"
    mkdir lib
else
    echo "\`lib' directory already exists, assuming library is compiled already"
    exit 0
fi

# print a warning, this isn't a fast process
echo "-----------------------------------------------------------"
echo ""
echo "                     COMPILING LIBRARY                     "
echo "               (this may take several hours)               "
echo "-----------------------------------------------------------"

# colpile the library
./cmd/lib-download/curl-the-scriptures.sh
./cmd/lib-download/clean-verses.sh
./cmd/lib-download/compress-json.sh