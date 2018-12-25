package parse

import (
	"fmt"
	"os"
	"strconv"
)

// flags type holds flags used for reference lookup
type flags struct {
	Context      int // not yet implemented
	Footnotes    bool
	JST          bool // not yet implemented
	Link         bool // not yet implemented
	Headings     bool
	HeadingsOnly bool
	Refs         bool
	Paragraphs   bool // not yet implemented
}

// Flags holds info needed for lookup,
// This is declared one once, populated by the cmd package, then passed into
// the Lookup methods
var Flags flags

// Lookuper interfaces references accepting chapter selections, or
// verse selections within a chapter. Lookup methods ought also to
// lookup with respect to command line args
type Lookuper interface {
	Lookup(flags) error
}

// makeRanve takes two numbers (as strings) and creates a range of ints-in-strings
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

// ReferenceVerses stores a scripture reference containing a series of verses
// all in the same chapter. It satisfies the Lookuper interface
type ReferenceVerses struct {
	Book    string
	Chapter string
	Verse   []string
}

// Lookup will look up the chapters referenced in the ReferenceChapters object
func (r *ReferenceVerses) Lookup(f flags) error {
	path := buildChapterRef(r.Book, r.Chapter)
	_, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("failed to stat %s %s (does that book/chapter exist?)", r.Book, r.Chapter)
	}

	chap, err := ReadChapter(path)

	if f.Refs {
		fmt.Printf("----%s %s----\n", r.Book, r.Chapter)
	}

	if f.Headings {
		fmt.Println(chap.Heading)
	}

	for _, vs := range r.Verse {
		if chap.Verses[vs].Text == "" {
			fmt.Printf("Failed: Chapter %s has no verse %s\n", r.Chapter, vs)
		}
		if f.Refs {
			fmt.Printf(" %s", vs)
		}
		if f.Footnotes {
			fmt.Printf(" %s\n", chap.Verses[vs].putFootnotes())
			fmt.Printf("    %s\n", chap.Verses[vs].formatFootnotes())
		} else {
			fmt.Printf(" %s\n", chap.Verses[vs].Text)
		}
	}

	// print an extra newline for good measure
	fmt.Print("\n")

	return nil
}

// ReferenceChapters stores a scripture reference that displays whole chapters.
// It satisfies the Lookuper interface
type ReferenceChapters struct {
	Book    string
	Chapter []string
}

// Lookup will look up the verses referenced in the ReferenceChapters object and
// display them
func (r *ReferenceChapters) Lookup(f flags) error {
	for _, lookupChap := range r.Chapter {
		path := buildChapterRef(r.Book, lookupChap)
		_, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("failed to stat %s %s (does that book/chapter exist?)", r.Book, lookupChap)
		}

		chap, err := ReadChapter(path)

		if f.Refs {
			fmt.Printf("----%s %s----\n", r.Book, lookupChap)
		}

		if f.Headings || f.HeadingsOnly {
			fmt.Println(chap.Heading)
		}

		// if we only wanted the headings...
		if f.HeadingsOnly {
			// then we're done
			fmt.Print("\n")
			continue
		}

		// print the text of the chapter
		for i := 1; i <= len(chap.Verses); i++ {
			verse := chap.Verses[strconv.Itoa(i)]
			if f.Refs {
				fmt.Printf(" %d", i)
			}
			if f.Footnotes {
				fmt.Printf(" %s\n", verse.putFootnotes())
				fmt.Printf("    %s\n\n", verse.formatFootnotes())
			} else {
				fmt.Printf(" %s\n", verse.Text)
			}
		}
	}

	return nil
}
