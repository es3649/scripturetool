package parse

import (
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/es3649/scripturetool/internal/lookup"
	"github.com/es3649/scripturetool/pkg/log"
	"github.com/sirupsen/logrus"
)

var refs []lookup.Lookuper

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

	log.Log.WithFields(logrus.Fields{"where": "parse", "args": fmt.Sprintf("%v", args)}).Info("Parsing these arguments")

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

	if lookup.Flags.UseStdout {
		displayAll(p.Results)
		return nil
	}
	return pipeToLess(p.Results)
}

func displayAll(refs []lookup.Lookuper) {

	// lookup all the references we got from the parser
	for _, reference := range refs {
		log.Log.WithFields(logrus.Fields{"where": "displayAll", "reference": fmt.Sprintf("%#v", reference)}).Info("Looking up Refrence")
		if err := reference.Lookup(lookup.Flags); err != nil {
			fmt.Printf("Error looking up the given reference: %#v\n", reference)
			fmt.Printf("  Error is: %v\n", err)
		}
	}
}

func pipeToLess(refs []lookup.Lookuper) error {
	// check the parse flags to see about using less

	r, w, err := os.Pipe()
	if err != nil {
		return fmt.Errorf("Error opening pipe to less: %v", err)
	}

	stdout := os.Stdout
	os.Stdout = w

	// open an instance of less
	less := exec.Command("less")
	less.Stdin = r
	less.Stdout = stdout
	less.Stderr = os.Stderr

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {

		if err := less.Run(); err != nil {
			log.Log.WithField("where", "displayAll").Error(fmt.Sprintf("Failed to execute less: %v", err))
		}

		if err = r.Close(); err != nil {
			log.Log.WithField("where", "displayAll").Warn("Failed to close pipe (read end)")
		}
		wg.Done()

	}()

	// lookup all the references we got from the parser
	for _, reference := range refs {
		if err := reference.Lookup(lookup.Flags); err != nil {
			fmt.Printf("Error looking up the given reference: %#v\n", reference)
			fmt.Printf("  Error is: %v\n", err)
		}
	}

	if err = w.Close(); err != nil {
		log.Log.WithField("where", "displayAll").Warn("Failed to close pipe (write end)")
	}
	// restore stdout
	os.Stdout = stdout

	wg.Wait()

	return nil
}
