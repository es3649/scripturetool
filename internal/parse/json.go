package parse

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Chapter stores a chapter or scripture
type Chapter struct {
	Book    string
	Chapter string
	Heading string
	Verses  map[string]Verse
}

// Verse stores a verse within a chapter
type Verse struct {
	Text      string
	Footnotes []Footnote
}

// Footnote stores a footnote within a verse
type Footnote struct {
	Position  int
	Reference string
}

// putFootnotes returns the Text field of the verse, but inserts a footnote
// identifier at each position indicated in the []Footnotes
func (v Verse) putFootnotes() string {
	var verseSlices []string
	var old = 0
	for i, footnote := range v.Footnotes {
		// add the slice of the text that goes from the last footnote to the current one
		verseSlices = append(verseSlices, v.Text[old:footnote.Position], "[", string(i+'a'), "]")
		old = footnote.Position
	}
	// add the slice from the position of the last footnote to the end of the verse
	verseSlices = append(verseSlices, v.Text[v.Footnotes[len(v.Footnotes)-1].Position:])
	return strings.Join(verseSlices, "")
}

// formatFootnotes returns formatted footnotes with their letter identifiers
func (v Verse) formatFootnotes() string {
	var footnotes string
	for _, note := range v.Footnotes {
		footnotes = footnotes + fmt.Sprintf("%s   ", note.Reference)
	}
	return footnotes
}

func buildChapterRef(book, chapter string) string {
	// lib location should be in the same working directory as the executable
	path := "./lib/"
	path = path + book + "/" + chapter + ".json.tar.gz"

	log.WithFields(logrus.Fields{"where": "buildChapterRef", "path": path}).Debug("Built a path to the chapter resource")
	return path
}

// ReadChapter opens a filepath, unzips and reads the json, and returns the chapter
func ReadChapter(path string) (*Chapter, error) {
	// open the file
	log.WithFields(logrus.Fields{"where": "ReadChapter", "path": path}).Debug("Reading the file...")
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("ReadChapter--failed to open path:\n%v", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.WithFields(logrus.Fields{"where": "anonymous closer", "error": err}).Errorf("Failed to close file")
		}
	}()

	// get a new gzip reader with the file
	gz, err := gzip.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("ReadChapter--failed to create gunzipper:\n%v", err)
	}
	defer func() {
		if err := gz.Close(); err != nil {
			log.WithFields(logrus.Fields{"where": "anonymous closer", "error": err}).Errorf("Failed to close gReader")
		}
	}()

	// get a new tar reader from the gzip reader
	tr := tar.NewReader(gz)
	_, err = tr.Next()
	if err != nil {
		return nil, fmt.Errorf("Failed to get the next tar file:\n%v", err)
	}

	// get data by reading all from the tar reader
	data, err := ioutil.ReadAll(tr)
	if err != nil {
		return nil, fmt.Errorf("ReadChapter--failed to read all from gzip:\n%v", err)
	}
	log.WithFields(logrus.Fields{"where": "ReadChapter", "data": string(data)}).Debug("Read these data")

	// unmarshal the data into a new chapter struct
	var ch Chapter
	err = json.Unmarshal(data, &ch)
	if err != nil {
		return nil, fmt.Errorf("ReadChapter--failed to unmarshal json:\n%v", err)
	}

	return &ch, nil
}
