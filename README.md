# scripturetool
A tool for accessing Latter-day Saint scripture via command line.

This project is specifically designed for faster scripture lookup

## Installation:

### Build from source
As of v0.3.0, you can actually run
```
sudo make install
```

It will download the library files and build executables to useful places

### Download binaries

Use the installer script by running your favorite version of
```
curl <some-url> | sudo sh
```

### Copyright Disclaimer

I do not own the content that this tool is designed to organize and recall. Intellectual Reserve, Inc. is not affiliated with nor do they endorse this project.

This content is publicly available at churchofjesuschrist.org/scriptures, and I personally invite everyone to read it and learn more about Jesus Christ and his restored Gospel.

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
