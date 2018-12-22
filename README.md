# scripturetool
A tool for accessing Latter-day Saint scripture via command line

## Installation:

Turns out you can just do 
```
make install
```

It will download the library files from lds.org and clean them up, but it will take several hours to do so.

### Copyright Disclaimer

I don't own the content that this tool is designed to organize and recall, however the content is made publicly available on lds.org/scriptures and all are invited to read it and learn more about Jesus Christ and his restored Gospel.

## Status:
 * Parses correctly formatted references
 * Library compiles
    * lib-download.sh works but takes HOURS to run, fix this
       * multithreading?
       * better router?
 * Add boundaries on the verses.
    * No chapter greater than 150 (Psalms)
    * No verse greater than 176 (Ps 119) 
 * Implement all cmd line flags
 * Multiple args doesn't work
 * Implement the cool features
    * Tome classes
    * Brace expansion {Matt,Mark,Luke}
    * Use of wildcard (*)
    * Semicolon stuff (reevaluate grammar)