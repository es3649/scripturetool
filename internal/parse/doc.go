// Package parse - Eric Steadman 2019
//
// This package handles the parsing of references from command line
// arguments and looking them up.
//
// It has a number of packages used for formatting.
//
// Package TODOs:
//
// Flags should become a struct literal, and should be removed from
// the specification of the Lookuper interface
//
// Wildcards and bash-style curly brace expansion needs to be implemented
// This will require changing the way that references are parsed, because
// '*', '{', '}', and ',' will probably need to be added to make sure
// parsing happens correctly
// we might also need markers within classes indicating the location of
// the * wildcard
//
// The grammar needs to be modified so that multiple references can
// be given as one argument. The grammar should also be recorded somewhere,
// here is probably as good as anywhere.
//
// TODO add bash-style curly brace lists
// add wildcards
// consider adding tomes
// <start>			::= <reference> <referencelist>
// <referencelist> 	::= ;<reference> <refrencelist> | ; | <lambda>
// <reference> 		::= <book> | <book> <chapter> <chapterlist> | <book> <chapter> : <verse> <verselist>
// <book>			::= <string>
// <string>			::= [::alnum::[]-] <string> | <\lambda>
// <lambda> 		::=
// <integer>		::= [0-9] <integer> | <lambda>
// <chapter>		::= <integer> | <integerrange>
// <chapterlist>	::= , <chapter> <chapterlist> | <lambda>
// <integerrange>	::= <integer> - <integer>
// <verse>			::= <integer> | <integerrange>
// <verselist>		::= , <verse> <verserange> | <lambda>
//
// It would also be good to fix the multithreading. As it was it would
// deadlock whenever there was a parse error.
// This might be fixed by passing tokens into the out channel IMMEDIATELY
// after the select statement deciding whether there was an error to handle.
package parse
