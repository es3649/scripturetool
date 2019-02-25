package parse

import (
	"fmt"
	"strconv"

	"github.com/es3649/scripturetool/internal/lookup"
	"github.com/es3649/scripturetool/pkg/log"
	"github.com/sirupsen/logrus"
)

type parserState int

const (
	pStartState       parserState = 0
	pBookNameNum      parserState = 1
	pBookStar         parserState = 2
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
	pChapStar         parserState = 16
	pVerseStar        parserState = 17
)

// parser takes tokens from an analyzer. It takes tokens from inChan
// (passed from the analyzer) and returns an error (or nil)
type parser struct {
	curState parserState
	inChan   chan token
	Results  []lookup.Lookuper
	curChRef lookup.ReferenceChapters
	curVsRef lookup.ReferenceVerses
	curBook  string
	curChap  string
	curVerse string
}

func newParser(refs []lookup.Lookuper, c chan token) *parser {
	return &parser{
		inChan:  c,
		Results: refs,
	}
}

// parseOrder ensures that the tokens are received in the correct order
func (p *parser) parseOrder() error {
	p.curVerse = ""
	p.curState = pStartState

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
				p.curBook = tok.Value
			case aBookState:
				p.curState = pBookName
				p.curBook = tok.Value
			case aStarState:
				p.curState = pBookStar
				p.curBook = tok.Value
			case aSemicolonState:
				// continue
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pBookNameNum:
			// a number must be followed by a dash
			switch tok.Type {
			case aDashState:
				p.curState = pBookName
				p.curBook = p.curBook + tok.Value
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pBookStar:
			// a dash must be followed by a book Name
			switch tok.Type {
			case aNumberState:
				p.curState = pChapNum
				p.curChap = tok.Value
			case aStarState:
				p.curState = pChapStar
				p.curChap = tok.Value
			case aSemicolonState:
				p.curState = pStartState
				curBook := lookup.ReferenceBook(p.curBook)
				p.Results = append(p.Results, &curBook)
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pBookName:
			// a book name can be the end of a reference (semicolon or EOF)
			// or be followed by a colon
			switch tok.Type {
			case aDashState:
				p.curBook = p.curBook + tok.Value
			case aBookState:
				p.curBook = p.curBook + tok.Value
			case aNumberState:
				p.curState = pChapNum
				p.curChap = tok.Value
			case aSemicolonState:
				curBook := lookup.ReferenceBook(p.curBook)
				p.Results = append(p.Results, &curBook)
				p.curState = pStartState
			case aStarState:
				// then we have a star for the chapter nameq
				p.curState = pChapStar
				p.curChap = tok.Value
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
				// create a ReferenceVerses object
				p.curVsRef = lookup.ReferenceVerses{
					Book:    p.curBook,
					Chapter: p.curChap,
				}

			case aDashState:
				p.curState = pChapRangeDash
				// get this chapter
				// we'll get the end chapter when we finish
				p.curChRef = lookup.ReferenceChapters{
					Book:    p.curBook,
					Chapter: append(make([]string, 0), p.curChap),
				}

			case aCommaState:
				p.curState = pChapComma
				// log this chapter, then we'll get the rest of them later
				p.curChRef = lookup.ReferenceChapters{
					Book:    p.curBook,
					Chapter: append(make([]string, 0), p.curChap),
				}

			case aSemicolonState:
				p.curState = pStartState
				// we're finished
				p.Results = append(p.Results, &lookup.ReferenceChapters{
					Book:    p.curBook,
					Chapter: append(make([]string, 0), p.curChap),
				})

			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapRangeDash:
			switch tok.Type {
			case aNumberState:
				p.curState = pChapRangeNum
				chaps, err := makeRange(p.curChap, tok.Value)
				if err != nil {
					return err
				}
				p.curChRef.Chapter = append(p.curChRef.Chapter, chaps...)
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapRangeNum:
			switch tok.Type {
			case aSemicolonState:
				p.curState = pStartState
				p.Results = append(p.Results, &p.curChRef)
			case aCommaState:
				p.curState = pChapComma
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapComma:
			switch tok.Type {
			case aNumberState:
				p.curState = pChapListNum
				p.curChap = tok.Value
				p.curChRef.Chapter = append(p.curChRef.Chapter, tok.Value)
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
				p.Results = append(p.Results, &p.curChRef)
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapListDash:
			switch tok.Type {
			case aNumberState:
				p.curState = pChapListRangeNum
				chaps, err := makeRange(p.curChap, tok.Value)
				if err != nil {
					return err
				}
				p.curChRef.Chapter = append(p.curChRef.Chapter, chaps...)
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapListRangeNum:
			switch tok.Type {
			case aCommaState:
				p.curState = pChapComma
			case aSemicolonState:
				p.curState = pStartState
				p.Results = append(p.Results, &p.curChRef)
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pChapStar:
			switch tok.Type {
			case aColonState:
				p.curState = pColon
			case aSemicolonState:
				p.Results = append(p.Results, &p.curChRef)
				p.curState = pStartState
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pColon:
			switch tok.Type {
			case aNumberState:
				p.curState = pVerseNum
				p.curVerse = tok.Value
				p.curVsRef = lookup.ReferenceVerses{
					Book:    p.curBook,
					Chapter: p.curChap,
					Verse:   append(make([]string, 0), p.curVerse),
				}
			case aStarState:
				// Then we have a star for a verse number, and we're done!
				p.curState = pStartState
				p.Results = append(p.Results, &lookup.ReferenceVerses{
					Book:    p.curBook,
					Chapter: p.curChap,
					Verse:   append(make([]string, 0), tok.Value),
				})
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

			/////////////
			// TODO all this below: the verse numbers don't get saved into Lookupers yet on multiple iterations
			/////////////
		case pVerseNum:
			switch tok.Type {
			case aDashState:
				p.curState = pVerseRangeDash
			case aCommaState:
				p.curState = pVerseComma
			case aSemicolonState:
				p.curState = pStartState
				p.Results = append(p.Results, &p.curVsRef)
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pVerseRangeDash:
			switch tok.Type {
			case aNumberState:
				p.curState = pVerseRangeNum
				// # we got a range
				verses, err := makeRange(p.curVerse, tok.Value)
				if err != nil {
					return err
				}
				p.curVsRef.Verse = append(p.curVsRef.Verse, verses...)
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pVerseRangeNum:
			switch tok.Type {
			case aCommaState:
				p.curState = pVerseComma
			case aSemicolonState:
				p.curState = pStartState
				p.Results = append(p.Results, &p.curVsRef)
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pVerseComma:
			switch tok.Type {
			case aNumberState:
				p.curState = pVerseNum
				p.curVerse = tok.Value
				p.curVsRef.Verse = append(p.curVsRef.Verse, p.curVerse)
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}

		case pVerseStar:
			switch tok.Type {
			case aSemicolonState:
				p.Results = append(p.Results, &p.curVsRef)
				p.curState = pStartState
			default:
				return fmt.Errorf("Invalid token received: %#v", tok)
			}
		}
	}
	log.Log.WithFields(logrus.Fields{"where": "parseOrder", "status": "success"}).Info("Finished Parsing")
	if p.curState != pStartState {
		return fmt.Errorf("End of line while parsing reference")
	}

	// p.curVerse should be empty if we parsed a chapter reference
	// if p.curVerse == "" {
	// 	p.Results = append(p.Results, &p.curChRef)
	// 	log.Log.WithFields(logrus.Fields{"where": "parseOrder", "reference": fmt.Sprintf("%#v", p.curChRef)}).Info("Logged a Chapter")
	// } else {
	// 	p.Results = append(p.Results, &p.curVsRef)
	// 	log.Log.WithFields(logrus.Fields{"where": "parseOrder", "reference": fmt.Sprintf("%#v", p.curVsRef)}).Info("Logged a Verse")
	// 	p.curVerse = ""
	// }

	return nil
}

// makeRange takes two numbers (as strings) and creates a range of ints-in-strings
// from the lower to the upper (if it's actually lower)
func makeRange(lower, upper string) ([]string, error) {
	var list []string
	l, _ := strconv.ParseInt(lower, 10, 64)
	lo := int(l)
	u, _ := strconv.ParseInt(upper, 10, 64)
	up := int(u)

	// bound the numbers to [1,176]
	if up > 176 {
		up = 176
	}
	if lo < 1 {
		lo = 1
	}

	if u <= l {
		return nil, fmt.Errorf("error in range: %d-%d", l, u)
	}

	for i := lo + 1; i <= up; i++ {
		list = append(list, strconv.Itoa(i))
	}

	return list, nil
}
