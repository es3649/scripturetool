#!/bin/bash

for file in $(ls lib/*/*.json); do
    echo "compressing $file"
    tar czf $file.tar.gz $file
done