# scripturetool
A tool for accessing Latter-day Saint scripture via command line

### Status:
 * Parses correctly formatted references
 * Library compiles
    * lib-download.sh works but takes HOURS to run
 * Add boundaries on the verses.
    * No chapter breater than 150 (Psalms)
    * No verse greater than 176 (Ps 119) 
 * Implement all cmd line flags
 * Multiple args doesn't work
 * Implement the cool features
    * Tome classes
    * Brace expansion {Matt,Mark,Luke}
    * Use of wildcard (*)
    * Semicolon stuff (reevaluate grammar)