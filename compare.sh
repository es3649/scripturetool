#!/bin/bash

if [[ "$1" == "" || "$2" == "" ]]; then
    echo "Usage: compare.sh [REFERENCE1] [REFERENCE2]"
    exit 1
fi

./scripturetool $1 | tr " " "\n" > $1.txt
./scripturetool $2 | tr " " "\n" > $2.txt
diff -y $1.txt $2.txt | less
rm $1.txt $2.txt
