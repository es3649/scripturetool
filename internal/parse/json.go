package parse

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// Chapter stores a chapter or scripture
type Chapter struct {
	Book    string
	Chapter int
	Heading string
	Verses  []Verse
}

// Verse stores a verse within a chapter
type Verse struct {
	Text      string
	Footnotes []Footnote
}

// PutFootnotes returns the Text field of the verse, but inserts a footnote
// identifier at each position indicated in the []Footnotes
func (v *Verse) PutFootnotes() string {
	return ""
}

// Footnote stores a footnote within a verse
type Footnote struct {
	Position  int
	Reference string
}

// ReadChapter opens a filepath, unzips and reads the json, and returns the chapter
func ReadChapter(path string) (*Chapter, error) {
	// open the file
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ReadChapter--failed to open path:\n%v", err)
	}

	// get a new gzip reader
	zReader, err := gzip.NewReader(f)
	defer func() {
		if err := zReader.Close(); err != nil {
			log.WithFields(logrus.Fields{"where": "anonymous closer", "error": err}).Errorf("Failed to close gReader")
		}
	}()

	// get data
	var data []byte
	_, err = zReader.Read(data)
	if err != nil {
		return nil, fmt.Errorf("ReadChapter--failed to uncompress data:\n%v", err)
	}

	// unmarshal the data into a new chapter struct
	var ch *Chapter
	err = json.Unmarshal(data, &ch)
	if err != nil {
		return nil, fmt.Errorf("ReadChapter--failed to unmarshal json:\n%v", err)
	}

	return ch, nil
}
