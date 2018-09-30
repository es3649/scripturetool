package parse

import (
	"strings"
	"unicode"
)

type analyzerState int

const (
	aStartState     analyzerState = 0
	aSemicolonState analyzerState = 1
	aColonState     analyzerState = 2
	aCommaState     analyzerState = 3
	aDashState      analyzerState = 4
	aNumberState    analyzerState = 5
	aBookState      analyzerState = 6
	aUndefinedState analyzerState = 7
	aEOFState       analyzerState = 8
)

type analyzer struct {
	value         string
	curState      analyzerState
	toParse       []string
	parseString   string
	pos           int
	outputChannel chan token
	stopChan      chan error
}

func newAnalyzer(parseStrs []string, c chan token, errChan chan error) *analyzer {
	return &analyzer{
		toParse:       parseStrs,
		outputChannel: c,
		stopChan:      errChan,
	}
}

func (a *analyzer) makeToken(tokType analyzerState) {
	a.outputChannel <- token{
		Value: strings.ToLower(a.value),
		Type:  a.curState,
	}
	a.value = ""
}

func (a *analyzer) analyze() error {
	for _, curString := range a.toParse {
		// no length 0 strings!
		if curString == "" {
			continue
		}
		if err := a.analyzeOne(curString); err != nil {
			close(a.outputChannel)
			return err
		}
	}
	close(a.outputChannel)
	return <-a.stopChan
}

func (a *analyzer) analyzeOne(curString string) error {
	a.curState = aStartState
	a.parseString = curString
	a.pos = -1 // -1 because the first thing we do is increment it
	// analyze each character in the string
	for {
		// stop analysis if we get a stop signal form the parser
		select {
		case err := <-a.stopChan:
			return err
		default:
			a.pos++
			if a.pos >= len(curString) {
				return a.finish()
			}
			c := rune(a.parseString[a.pos])
			switch a.curState {
			case aStartState:
				a.start(c)
			case aSemicolonState:
				a.semicolon(c)
			case aColonState:
				a.colon(c)
			case aCommaState:
				a.comma(c)
			case aDashState:
				a.dash(c)
			case aNumberState:
				a.number(c)
			case aBookState:
				a.book(c)
			}
		}
	}
}

func (a *analyzer) finish() error {
	switch a.curState {
	case aNumberState:
		a.makeToken(aNumberState)
	case aBookState:
		a.makeToken(aBookState)
	}
	return nil
}

func (a *analyzer) semicolon(c rune) {
	a.makeToken(aSemicolonState)
	a.start(c)
}

func (a *analyzer) colon(c rune) {
	a.makeToken(aColonState)
	a.start(c)
}

func (a *analyzer) comma(c rune) {
	a.makeToken(aCommaState)
	a.start(c)
}

func (a *analyzer) dash(c rune) {
	a.makeToken(aDashState)
	a.start(c)
}

func (a *analyzer) number(c rune) {
	if unicode.IsDigit(c) {
		a.value += string(c)
		return
	}
	a.makeToken(aNumberState)
	a.start(c)
}

func (a *analyzer) book(c rune) {
	if unicode.IsLetter(c) {
		a.value += string(c)
		return
	}
	a.makeToken(aBookState)
	a.start(c)
}

func (a *analyzer) start(c rune) {
	a.value += string(a.parseString[a.pos])
	switch a.parseString[a.pos] {
	case ':':
		a.curState = aColonState
	case ';':
		a.curState = aSemicolonState
	case '-':
		a.curState = aDashState
	case ',':
		a.curState = aCommaState
	default:
		if unicode.IsDigit(c) {
			a.curState = aNumberState
			return
		} else if unicode.IsLetter(c) {
			a.curState = aBookState
			return
		} else if unicode.IsSpace(c) {
			a.curState = aStartState
			a.value = ""
			return
		}
		// got an unknown character
		a.makeToken(aUndefinedState)
	}
}
