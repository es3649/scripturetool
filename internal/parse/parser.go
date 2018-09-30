package parse

import "fmt"

type parserState int

type parser struct {
	state   parserState
	inChan  chan token
	errChan chan error
	results []Lookuper
}

func newParser(refs []Lookuper, c chan token, e chan error) *parser {
	return &parser{
		inChan:  c,
		errChan: e,
		results: refs,
	}
}

func (p *parser) parseOrder() {
	defer close(p.errChan)

	// for each token we get
	for tok := range p.inChan {
		fmt.Printf("%#v\n", tok)
	}
	p.errChan <- nil
	fmt.Println("Finished Parsing (errChan closed)")
}
