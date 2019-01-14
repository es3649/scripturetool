package parse

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

var refs []Lookuper

// Parse parses the command line arguments then TODO
func Parse(args []string) (err error) {
	// expand the args
	for i := range args {
		args[i], err = expandBraces(args[i])
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
