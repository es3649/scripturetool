package parse

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type parserState int

const (
	pStartState       parserState = 0
	pBookNameNum      parserState = 1
	pBookNameDash     parserState = 2
	pBookName         parserState = 3
	pChapNum          parserState = 4
	pChapRangeDash    parserState = 5
	pChapRangeNum     parserState = 6
	pChapComma        parserState = 7
	pChapListNum      parserState = 8
	pChapListDash     parserState = 9
	pChapListRangeNum parserState = 10
	pColon            parserState = 11
	pVerseNum         parserState = 12
	pVerseRangeDash   parserState = 13
	pVerseRangeNum    parserState = 14
	pVerseComma       parserState = 15
)

// parser takes tokens from an analyzer. It takes tokens from inChan
// (passed from the analyzer) and retruns an error (or nil) to the
// analyzer to terminate it if needed via errChan.
type parser struct {
	curState parserState
	inChan   chan token
	results  []Lookuper
	curBook  string
	curChap  []int
}

func newParser(refs []Lookuper, c chan token) *parser {
	return &parser{
		inChan:  c,
		results: refs,
	}
}

// parseOrder ensures that the tokens are received in the correct order
func (p *parser) parseOrder() error {

	// for each token we get
	for tok := range p.inChan {
		// fmt.Printf("%v\n", tok)
		// fmt.Printf("parser state: %d\n", p.curState)
		switch p.curState {

		case pStartState:
			// it can start with a bookNum or a BookName
			switch tok.Type {
			case aNumberState:
				p.curState = pBookNameNum
			case aBookState:
				p.curState = pBookName
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pBookNameNum:
			// a number must be followed by a dash
			switch tok.Type {
			case aDashState:
				p.curState = pBookNameDash
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pBookNameDash:
			// a dash must be followed by a book Name
			switch tok.Type {
			case aBookState:
				p.curState = pBookName
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pBookName:
			// a book name can be the end of a reference (semicolon or EOF)
			// or be followed by a colon
			switch tok.Type {
			case aNumberState:
				p.curState = pChapNum
			case aSemicolonState:
				p.curState = pStartState
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapNum:
			// we can get a colon, then we move to verses
			// we can get a dash, indicating a chapter range
			// we can get a comma, indicating a list
			// this can be the end (EOF or semicolon)
			switch tok.Type {
			case aColonState:
				p.curState = pColon
			case aDashState:
				p.curState = pChapRangeDash
			case aCommaState:
				p.curState = pChapComma
			case aSemicolonState:
				p.curState = pStartState
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapRangeDash:
			switch tok.Type {
			case aNumberState:
				p.curState = pChapRangeNum
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapRangeNum:
			switch tok.Type {
			case aSemicolonState:
				p.curState = pStartState
			case aColonState:
				p.curState = pColon
			case aCommaState:
				p.curState = pChapComma
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapComma:
			switch tok.Type {
			case aNumberState:
				p.curState = pChapListNum
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapListNum:
			switch tok.Type {
			case aCommaState:
				p.curState = pChapComma
			case aDashState:
				p.curState = pChapListDash
			case aSemicolonState:
				p.curState = pStartState
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapListDash:
			switch tok.Type {
			case aNumberState:
				p.curState = pChapListRangeNum
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapListRangeNum:
			switch tok.Type {
			case aCommaState:
				p.curState = pChapComma
			case aSemicolonState:
				p.curState = pStartState
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pColon:
			switch tok.Type {
			case aNumberState:
				p.curState = pVerseNum
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pVerseNum:
			switch tok.Type {
			case aDashState:
				p.curState = pVerseRangeDash
			case aCommaState:
				p.curState = pVerseComma
			case aSemicolonState:
				p.curState = pStartState
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pVerseRangeDash:
			switch tok.Type {
			case aNumberState:
				p.curState = pVerseRangeNum
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pVerseRangeNum:
			switch tok.Type {
			case aCommaState:
				p.curState = pVerseComma
			case aSemicolonState:
				p.curState = pStartState
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pVerseComma:
			switch tok.Type {
			case aNumberState:
				p.curState = pVerseNum
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}
		}
	}
	log.WithFields(logrus.Fields{"where": "parseOrder", "status": "success"}).Info("Finished Parsing (errChan closed)")
	if p.curState != pStartState {
		return fmt.Errorf("End of line while parsing reference")
	}

	return nil
}
