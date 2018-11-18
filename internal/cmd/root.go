/**
 * Package cmd - handling commands for scripturetool
 * root.go - the root command
 */

package cmd

import (
	"fmt"

	"github.com/es3649/scripturetool/internal/parse"
	"github.com/spf13/cobra"
)

var flags = struct {
	Verbosity    int
	Context      int
	Footnotes    bool
	JST          bool
	Link         bool
	Headings     bool
	HeadingsOnly bool
	Refs         bool
	Paragraphs   bool
}{}

// Execute runs the command
func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:   "scripturetool",
	Short: "A command line tool for the Latter-day Saint standard works",
	Long: `Search and display scriptures from the command line.
The scriptures used in these libraries are owned by Intellectual Reserve Inc. And are available for free at https://lds.org/scriptures.
The Church of Jesus Christ of Latter-day Saints is not affiliated with, nor endorses this project.

Use 'scripturetool help' for explanation of arguments and abbreviations`,
	Args: cobra.MinimumNArgs(1),
	PreRun: func(_ *cobra.Command, _ []string) {
		// set the verbosity of the logger(s)
		parse.SetVerbosity(flags.Verbosity)
	},
	Run: run,
}

func init() {
	// add subcommands
	rootCmd.AddCommand(helpCmd)

	// add flags
	rootCmd.Flags().IntVarP(&flags.Verbosity, "verbose", "v", 0, "varying levels of verbosity")
	rootCmd.Flags().IntVarP(&flags.Context, "context", "c", 0, "show verses before and after the selected verse(s)")
	rootCmd.Flags().BoolVarP(&flags.Footnotes, "footnotes", "f", false, "display footnotes")
	rootCmd.Flags().BoolVarP(&flags.JST, "jst", "j", false, "show only JST footnotes")
	rootCmd.Flags().BoolVarP(&flags.Link, "link", "l", false, "show the verses referenced in footnotes of selected verse(s)")
	rootCmd.Flags().BoolVarP(&flags.Headings, "headings", "H", false, "Show headings for selected chapter(s)")
	rootCmd.Flags().BoolVar(&flags.HeadingsOnly, "headings-only", false, "Show only the headings for selected chapter(s)")
	rootCmd.Flags().BoolVarP(&flags.Refs, "refs", "r", true, "show full refereces with each verse displayed")
	rootCmd.Flags().BoolVarP(&flags.Paragraphs, "pars", "p", false, "print text in paragraphs")
}

// run parses the references from the arguments and sends them for retreival and display
func run(cmd *cobra.Command, args []string) {

	err := parse.Parse(args)

	if err != nil {
		fmt.Printf("Got error `%v'\n", err)
	}
	// 	ch, err := parse.ReadChapter("lib/moro/5.json.tar.gz")
	// 	if ch == nil {
	// 		ch = &parse.Chapter{}
	// 	}

	// 	fmt.Printf(`Chapter:
	// %v

	// Error: %v
	// 		`, *ch, err)

	// 	for i, verse := range ch.Verses {
	// 		v := verse.PutFootnotes()
	// 		fmt.Printf("Verse %s: %s\n\n", i, v)
	// 	}
}
