package parse

import (
	"fmt"
	"strings"

	"github.com/es3649/scripturetool/pkg/log"
)

// TODO this only works with one set of braces, redo it with graph theory
// and recursion
// {matt|    {1|
// |mark| 1: |2| ;
// |luke}    |3}
func expandBraces(s string) (string, error) {
	log.Log.WithField("value", s).Info("Expanding the argument")
	i := strings.Index(s, "{")
	j := strings.Index(s, "}")
	// for each ocurrence of '{' and '}', expand them
	for i >= 0 && j >= 0 {
		// make sure there is no mismatch, if one of them exists,
		// we need to have the other to match with it
		if (i == -1) != (j == -1) {
			return "", fmt.Errorf("Mismatched braces: '{' at pos %d and '}' at pos %d", i, j)
		}
		// slice the string at the given indicies and expand it
		prefix := s[:i]
		suffix := s[j+1:]
		pieces := strings.Split(s[i+1:j], ",")
		s = ""
		for _, piece := range pieces {
			log.Log.WithField("value", piece).Debug("Building with piece")
			s += prefix + piece + suffix + ";"
		}
		// reset the conditions to be sure we're good
		i = strings.Index(s, "{")
		j = strings.Index(s, "}")
		log.Log.WithField("value", s).Debug("completed an expansion")
	}
	log.Log.WithField("value", s).Info("Expanded the argument")
	return s, nil
}

// expandTomeClasses will expand references with tome wildcards in them, such as
// [ot], [nt], [bofm], [dc], [pgp]
// TODO decide on good functionality for this function
// Should we call it on unprocessed args, or on Lookupers?
func expandTomeClasses(s string) (string, error) {
	return s, nil
}
