# scripturetool
A tool for accessing Latter-day Saint scripture via command line.

This project is specifically designed for faster scripture lookup

## Installation:

Turns out you can just do 
```
make install
```

It will download the library files from lds.org and clean them up, but it will take several hours to do so.

### Copyright Disclaimer

I do not own the content that this tool is designed to organize and recall. Intellectual Reserve, Inc. is not affiliated with nor do they endorse this project.

This content is publicly available on lds.org/scriptures, and I personally invite everyone to read it and learn more about Jesus Christ and his restored Gospel.

## Status:
 * Parses correctly formatted references
    * Looks up entire books or entire tomes
 * Default pages to `less`
    * Reverts to stdout with the `-o` flag
 * Library compiles
    * lib-download.sh works but takes HOURS to run, fix this
       * multithreading?
       * better router?
 * Add boundaries on the verses.
    * No chapter greater than 150 (Psalms)
    * No verse greater than 176 (Ps 119) 
 * Implement all cmd line flags
 * Implement the cool features
    * Brace expansion {Matt,Mark,Luke}
    * Use of wildcard (*)