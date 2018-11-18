package parse

import "fmt"

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
	pSemicolon        parserState = 16
)

type parser struct {
	curState parserState
	inChan   chan token
	errChan  chan error
	results  []Lookuper
	curBook  string
	curChap  []int
}

func newParser(refs []Lookuper, c chan token, e chan error) *parser {
	return &parser{
		inChan:  c,
		errChan: e,
		results: refs,
	}
}

func (p *parser) parseOrder() {

	// for each token we get
	for tok := range p.inChan {
		fmt.Printf("%v\n", tok)
		switch p.curState {
		case pStartState:
			p.startState(tok)
		case pBookNameNum:
			p.BookNameNum(tok)
		case pBookNameDash:
			p.BookNameDash(tok)
		case pBookName:
			p.BookName(tok)
		case pChapNum:
			p.ChapNum(tok)
		case pChapRangeDash:
			p.ChapRangeDash(tok)
		case pChapRangeNum:
			p.ChapRangeNum(tok)
		case pChapComma:
			p.ChapComma(tok)
		case pChapListNum:
			p.ChapListNum(tok)
		case pChapListDash:
			p.ChapListDash(tok)
		case pChapListRangeNum:
			p.ChapListRangeNum(tok)
		case pColon:
			p.Colon(tok)
		case pVerseNum:
			p.VerseNum(tok)
		case pVerseRangeDash:
			p.VerseRangeDash(tok)
		case pVerseRangeNum:
			p.VerseRangeNum(tok)
		case pVerseComma:
			p.VerseComma(tok)
		case pSemicolon:
			p.semicolon(tok)
		}
	}

	close(p.errChan)
	log.WithField("where", "parseOrder").Info("Finished Parsing (errChan closed)")
}

func (p *parser) startState(tok token) {
	switch tok.Type {
	case aNumberState:
		p.curState = pBookNameNum
	case aBookState:
		p.curState = pBookName
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) BookNameNum(tok token) {
	switch tok.Type {
	case aDashState:
		p.curState = pBookNameDash
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) BookNameDash(tok token) {
	switch tok.Type {
	case aBookState:
		p.curState = pBookName
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) BookName(tok token) {
	switch tok.Type {
	case aNumberState:
		p.curState = pChapNum
	case aSemicolonState:
		p.curState = pSemicolon
	case aEOFState:
		p.errChan <- nil
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) ChapNum(tok token) {
	switch tok.Type {
	case aColonState:
		p.curState = pColon
	case aDashState:
		p.curState = pChapRangeDash
	case aCommaState:
		p.curState = pChapComma
	case aSemicolonState:
		p.curState = pSemicolon
	case aEOFState:
		p.errChan <- nil
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) ChapRangeDash(tok token) {
	switch tok.Type {
	case aNumberState:
		p.curState = pChapRangeNum
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) ChapRangeNum(tok token) {
	switch tok.Type {
	case aSemicolonState:
		p.curState = pSemicolon
	case aColonState:
		p.curState = pColon
	case aCommaState:
		p.curState = pChapComma
	case aEOFState:
		p.errChan <- nil
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) ChapComma(tok token) {
	switch tok.Type {
	case aNumberState:
		p.curState = pChapListNum
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) ChapListNum(tok token) {
	switch tok.Type {
	case aCommaState:
		p.curState = pChapComma
	case aDashState:
		p.curState = pChapListDash
	case aSemicolonState:
		p.curState = pSemicolon
	case aEOFState:
		p.errChan <- nil
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) ChapListDash(tok token) {
	switch tok.Type {
	case aNumberState:
		p.curState = pChapListRangeNum
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) ChapListRangeNum(tok token) {
	switch tok.Type {
	case aCommaState:
		p.curState = pChapComma
	case aSemicolonState:
		p.curState = pSemicolon
	case aEOFState:
		p.errChan <- nil
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) Colon(tok token) {
	switch tok.Type {
	case aNumberState:
		p.curState = pVerseNum
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) VerseNum(tok token) {
	switch tok.Type {
	case aDashState:
		p.curState = pVerseRangeDash
	case aCommaState:
		p.curState = pVerseComma
	case aSemicolonState:
		p.curState = pSemicolon
	case aEOFState:
		p.errChan <- nil
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) VerseRangeDash(tok token) {
	switch tok.Type {
	case aNumberState:
		p.curState = pVerseRangeNum
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) VerseRangeNum(tok token) {
	switch tok.Type {
	case aCommaState:
		p.curState = pVerseComma
	case aSemicolonState:
		p.curState = pSemicolon
	case aEOFState:
		p.errChan <- nil
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) VerseComma(tok token) {
	switch tok.Type {
	case aNumberState:
		p.curState = pVerseNum
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}

func (p *parser) semicolon(tok token) {
	switch tok.Type {
	case aNumberState:
		p.curState = pBookNameNum
	case aBookState:
		p.curState = pBookName
	case aSemicolonState:
		// Then stay in the state. This lets the expander work
	case aEOFState:
		p.errChan <- nil
	default:
		p.errChan <- fmt.Errorf("Invalid token received: %#v", tok)
	}
}
