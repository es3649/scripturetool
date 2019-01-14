package parse

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/sirupsen/logrus"
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
	aStarState      analyzerState = 8
)

type token struct {
	Type  analyzerState
	Value string
}

// analyzer handles analysis of characters in an input string. It groups them into
// lexical token (parse.token) objects and passes them into outputChan. It constantly
// checks for errors from the stopChan, and stops processing if it gets one.
type analyzer struct {
	value       string
	curState    analyzerState
	toParse     []string
	parseString string
	pos         int
	outputChan  chan token // puts into this channel
}

// newAnalyzer creates a new analyzer object
func newAnalyzer(parseStrs []string, c chan token) *analyzer {
	return &analyzer{
		toParse:    parseStrs,
		outputChan: c,
	}
}

func (a *analyzer) makeToken(tokType analyzerState) {
	tok := token{
		Value: strings.ToLower(a.value),
		Type:  tokType,
	}
	log.WithField("token", fmt.Sprintf("%#v", tok)).Info("Token found")
	a.outputChan <- tok
	a.value = ""
}

func (a *analyzer) analyze() (err error) {
	for _, curString := range a.toParse {
		// no length 0 strings!
		if curString == "" {
			fmt.Println("Got empty arg")
			continue
		}
		log.WithFields(logrus.Fields{"where": "analyze", "arg": curString}).Info("Analyzing argument")
		if err = a.analyzeOne(curString); err != nil {
			log.WithFields(logrus.Fields{"where": "analyze", "status": "error"}).Info("Finished Analyzing (outputChan closed)")
			close(a.outputChan)
			fmt.Print("analyze-returning\n")
			return err
		}
	}
	// fmt.Print("finishing all analysis\n")
	a.makeToken(aSemicolonState)
	// fmt.Print("made final token\n")
	close(a.outputChan)
	log.WithFields(logrus.Fields{"where": "analyze", "status": "success"}).Info("Finished Analyzing (outputChan closed)")
	return nil
}

func (a *analyzer) analyzeOne(curString string) error {

	a.curState = aStartState
	a.parseString = curString
	a.pos = -1 // -1 because the first thing we do is increment it
	// analyze each character in the string
	for {
		a.pos++
		if a.pos >= len(curString) {
			return a.finish()
		}
		c := rune(a.parseString[a.pos])
		log.WithFields(logrus.Fields{"where": "analyze", "character": string(c), "state": a.curState}).Debug("Parsed character")
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
		case aStarState:
			a.star(c)
		case aNumberState:
			a.number(c)
		case aBookState:
			a.book(c)
		}
	}
}

func (a *analyzer) finish() error {
	switch a.curState {
	case aNumberState:
		a.makeToken(aNumberState)
		// fmt.Print("finish analysis-number state\n")
		return nil
	case aBookState:
		// fmt.Print("finish analysis-book state\n")
		a.makeToken(aBookState)
		return nil
	default:
		// fmt.Print("finish analysis-bad state\n")
		return fmt.Errorf("Invalid end of string: %s", a.value)
	}
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

func (a *analyzer) star(c rune) {
	a.makeToken(aStarState)
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
	if unicode.IsLetter(c) || c == '[' || c == ']' {
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
	case '*':
		a.curState = aStarState
	case '[':
		a.curState = aBookState
	case ']':
		a.curState = aBookState
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
