#!/bin/bash

pushd $GOPATH/src/github.com/es3649/scripturetool

if [[ ! -e lib/ ]]; then
    echo "creating \`lib' directory"
    mkdir lib
else
    ./cmd/lib-download/clean-lib.sh
fi

./cmd/lib-download/curl-the-scriptures.sh
./cmd/lib-download/clean-verses.sh
./cmd/lib-download/compress-json.sh