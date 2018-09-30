package parse

import "sync"

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

// Parse parses the command line arguments then TODO
func Parse(args []string) (err error) {
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
