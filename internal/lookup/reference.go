package lookup

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/es3649/scripturetool/internal/scriptures"
	"github.com/es3649/scripturetool/pkg/log"
	"github.com/sirupsen/logrus"
)

// flags type holds flags used for reference lookup
type flags struct {
	Context      int // not yet implemented
	Footnotes    bool
	JST          bool   // not yet implemented
	Language     string // not yet implemented
	Link         bool   // not yet implemented
	Headings     bool
	HeadingsOnly bool
	Refs         bool
	RefsFull     bool // not yet implemented
	Paragraphs   bool // not yet implemented
	UseStdout    bool
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

		// if it's empty, then just say an error
		if chap.Verses[vs].Text == "" {
			fmt.Printf("Failed: Chapter %s has no verse %s\n", r.Chapter, vs)
		}

		// handle the references
		if f.RefsFull {
			// RefsFull prints the book and chapter for the
			fmt.Printf(" [%s %s:%s]", r.Book, r.Chapter, vs)
		} else if f.Refs {
			// Refs means print the verse number before the verse
			fmt.Printf(" %s", vs)
		}

		// should we put the footnotes in?
		if f.Footnotes {
			fmt.Printf(" %s\n", chap.Verses[vs].putFootnotes())
			fmt.Printf("    %s\n", chap.Verses[vs].formatFootnotes())
		} else {
			fmt.Printf(" %s\n", chap.Verses[vs].Text)
		}
	}

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

		err = lookupChapterFromPath(r.Book, lookupChap, path)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceChapters.Lookup", "path": path}).Error("Failed to lookup the given chapter")
		}
	}

	return nil
}

func lookupChapterFromPath(book, chapnum, path string) error {
	chap, err := ReadChapter(path)
	if err != nil {
		return err
	}

	if Flags.Refs {
		fmt.Printf("----%s %s----\n", book, chapnum)
	}

	if Flags.Headings || Flags.HeadingsOnly {
		fmt.Println(chap.Heading)
	}

	// if we only wanted the headings...
	if Flags.HeadingsOnly {
		// then we're done
		fmt.Print("\n")
		return nil
	}

	// print the text of the chapter
	for i := 1; i <= len(chap.Verses); i++ {
		verse := chap.Verses[strconv.Itoa(i)]

		// handle references
		if Flags.RefsFull {
			// RefsFull prints the book and chapter for the
			fmt.Printf(" [%s %s:%d]", book, chapnum, i)
		} else if Flags.Refs {
			// Refs means print the verse numbers
			fmt.Printf(" %d", i)
		}

		// should we put the footnotes?
		if Flags.Footnotes {
			fmt.Printf(" %s\n", verse.putFootnotes())
			fmt.Printf("    %s\n\n", verse.formatFootnotes())
		} else {
			fmt.Printf(" %s\n", verse.Text)
		}
	}
	return nil
}

// ReferenceBook stores a scripture reference that references a whole book
// It satisfies the Lookuper interface
type ReferenceBook string

// Lookup will look up and display the verses references in the ReferenceBook
// object
func (r *ReferenceBook) Lookup(f flags) error {
	// check if the book is a tome, if it is, look it alllllll up
	var err error
	switch string(*r) {
	case "[all]":
		err = lookupTome(scriptures.OldTestament)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.Lookup", "tome": "OldTestament"}).Error(fmt.Sprintf("%v", err))
		}
		err = lookupTome(scriptures.NewTestament)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.Lookup", "tome": "NewTestament"}).Error(fmt.Sprintf("%v", err))
		}
		err = lookupTome(scriptures.BookOfMormon)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.Lookup", "tome": "BookOfMormon"}).Error(fmt.Sprintf("%v", err))
		}
		err = lookupTome(scriptures.DoctrineAndCovenants)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.Lookup", "tome": "DoctrineAndCovenants"}).Error(fmt.Sprintf("%v", err))
		}
		err = lookupTome(scriptures.PearlOfGreatPrice)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.Lookup", "tome": "PearlOfGreatPrice"}).Error(fmt.Sprintf("%v", err))
		}
	case "[ot]":
		err = lookupTome(scriptures.OldTestament)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.Lookup", "tome": "OldTestament"}).Error(fmt.Sprintf("%v", err))
		}
	case "[nt]":
		err = lookupTome(scriptures.NewTestament)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.Lookup", "tome": "NewTestament"}).Error(fmt.Sprintf("%v", err))
		}
	case "[bofm]":
		err = lookupTome(scriptures.BookOfMormon)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.Lookup", "tome": "BookOfMormon"}).Error(fmt.Sprintf("%v", err))
		}
	case "[dc]":
		err = lookupTome(scriptures.DoctrineAndCovenants)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.Lookup", "tome": "DoctrineAndCovenants"}).Error(fmt.Sprintf("%v", err))
		}
	case "[pgp]":
		err = lookupTome(scriptures.PearlOfGreatPrice)
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.Lookup", "tome": "PearlOfGreatPrice"}).Error(fmt.Sprintf("%v", err))
		}
	default:
		return r.fetch()
	}
	return nil
}

func (r *ReferenceBook) fetch() error {
	bookPath := "./lib/" + string(*r)
	files, err := ioutil.ReadDir(bookPath)
	if err != nil {
		return err
	}

	var chaps []string
	for i := 1; i <= len(files); i++ {
		chaps = append(chaps, strconv.Itoa(i))
	}
	// build a ReferenceChapter object
	b := &ReferenceChapters{
		Book:    string(*r),
		Chapter: chaps,
	}

	// call Lookup on that
	err = b.Lookup(Flags)
	if err != nil {
		log.Log.WithFields(logrus.Fields{"where": "ReferenceBook.fetch", "book": *r}).Error("Failed to lookup the given chapter")
	}

	return nil
}

func lookupTome(tome scriptures.Tome) error {
	for _, book := range tome {
		// create a ReferenceBook for each book in that tome
		var r = ReferenceBook(book)
		// lookup that book
		err := r.Lookup(Flags)
		// if we get an error, log it
		if err != nil {
			log.Log.WithFields(logrus.Fields{"where": "lookupTome", "book": book}).Error(fmt.Sprintf("%v", err))
		}
	}
	return nil
}
