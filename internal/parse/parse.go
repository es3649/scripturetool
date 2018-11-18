package parse

import (
	"fmt"
	"strings"
	"sync"
)

type token struct {
	Type  analyzerState
	Value string
}

// Reference stores a scripture reference. It needs to be interfaced
// with a lookup option, and we need types containing just chapters,
// and ones with verses, and maybe ones with ranges of chapters and
// or verses
type Reference struct {
	Book    string
	Chapter int
	Verse   int
}

// Lookuper interfaces references accepting chapter selections, or
// verse selections within a chapter. Lookup methods ought also to
// lookup with respect to command line args
type Lookuper interface {
	Lookup() error
}

var refs []Lookuper

// TODO this only works with one set of braces, redo it with graph theory
// and recursion
// [matt     [1
//  mark  1:  2  ;
//  luke]     3]
func expand(s string) (string, error) {
	log.WithField("value", s).Info("Expanding the argument")
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
			log.WithField("value", piece).Debug("Building with piece")
			s += prefix + piece + suffix + ";"
		}
		// reset the conditions to be sure we're good
		i = strings.Index(s, "{")
		j = strings.Index(s, "}")
		log.WithField("value", s).Debug("completed an expansion")
	}
	log.WithField("value", s).Info("Expanded the argument")
	return s, nil
}

// Parse parses the command line arguments then TODO
func Parse(args []string) (err error) {
	// expand the args
	for i := range args {
		args[i], err = expand(args[i])
		if err != nil {
			return err
		}
	}

	fmt.Print(args, "\n")

	analysisResultsChan := make(chan token)
	errChan := make(chan error)

	a := newAnalyzer(args, analysisResultsChan, errChan)
	p := newParser(refs, analysisResultsChan, errChan)

	var w sync.WaitGroup

	// run analysis in a thread
	w.Add(1)
	go func() {
		err = a.analyze()
		w.Done()
	}()

	// parse in another thread
	w.Add(1)
	go func() {
		// TODO get results from here
		p.parseOrder()
		w.Done()
	}()

	w.Wait()

	return err
}
