package parse

import (
	"fmt"
	"strings"
	"sync"

	"github.com/sirupsen/logrus"
)

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
		args[i] = PutAbbrevs(args[i])
	}

	log.WithFields(logrus.Fields{"where": "parse", "args": fmt.Sprintf("%v", args)}).Info("Parsing these arguments")

	analysisResultsChan := make(chan token, 50)

	a := newAnalyzer(args, analysisResultsChan)
	p := newParser(refs, analysisResultsChan)

	var w sync.WaitGroup

	var aErr, pErr error

	// run analysis in a thread
	w.Add(1)
	go func() {
		aErr = a.analyze()
		w.Done()
		return
	}()

	// parse in another thread
	w.Add(1)
	go func() {
		// TODO get results from here
		pErr = p.parseOrder()
		w.Done()
		return
	}()

	w.Wait()

	if aErr != nil {
		return aErr
	} else if pErr != nil {
		return pErr
	}

	// lookup all the references we got from the parser
	for _, reference := range p.Results {
		if err := reference.Lookup(Flags); err != nil {
			fmt.Printf("Error looking up the given reference: %#v\n", reference)
			fmt.Printf("  Error is: %v\n", err)
		}
	}

	return nil
}
