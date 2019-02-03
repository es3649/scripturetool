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

// Execute runs the command
func Execute() error {
	return rootCmd.Execute()
}

var verbosity int

var rootCmd = &cobra.Command{
	Use:   "scripturetool",
	Short: "A command line tool for the Latter-day Saint standard works",
	Long: `Search and display scriptures from the command line.
The scriptures used in these libraries are owned by Intellectual Reserve Inc. And are available for free at https://lds.org/scriptures.
The Church of Jesus Christ of Latter-day Saints is not affiliated with, nor endorses this project.

Use 'scripturetool info' for explanation of arguments and abbreviations`,
	Args: cobra.MinimumNArgs(1),
	PreRun: func(_ *cobra.Command, _ []string) {
		// set the verbosity of the logger(s)
		parse.SetVerbosity(verbosity)
	},
	Run: run,
}

func init() {
	// add subcommands
	rootCmd.AddCommand(helpCmd)

	// add flags
	rootCmd.Flags().IntVarP(&verbosity, "verbose", "v", 0, "varying levels of verbosity")
	rootCmd.Flags().IntVarP(&parse.Flags.Context, "context", "c", 0, "show verses before and after the selected verse(s)")
	rootCmd.Flags().BoolVarP(&parse.Flags.Footnotes, "footnotes", "f", false, "display footnotes")
	rootCmd.Flags().BoolVarP(&parse.Flags.JST, "jst", "j", false, "show only JST footnotes")
	rootCmd.Flags().BoolVarP(&parse.Flags.Link, "link", "L", false, "show the verses referenced in footnotes of selected verse(s)")
	rootCmd.Flags().BoolVarP(&parse.Flags.Headings, "headings", "H", false, "Show headings for selected chapter(s)")
	rootCmd.Flags().BoolVar(&parse.Flags.HeadingsOnly, "headings-only", false, "Show only the headings for selected chapter(s)")
	rootCmd.Flags().BoolVarP(&parse.Flags.Refs, "refs", "r", false, "hide chapter and verse references with each verse displayed")
	rootCmd.Flags().BoolVarP(&parse.Flags.RefsFull, "refs-full", "R", false, "show full references (book and chapter) with each verse displayed")
	rootCmd.Flags().BoolVarP(&parse.Flags.Less, "less", "l", false, "print to stdout instead of opening in less")
	rootCmd.Flags().BoolVarP(&parse.Flags.Paragraphs, "pars", "p", false, "print text in paragraphs")
	rootCmd.Flags().StringVarP(&parse.Flags.Language, "lang", "L", "eng", "language to display in")
}

// run parses the references from the arguments and sends them for retreival and display
func run(cmd *cobra.Command, args []string) {
	// invert the refs flag
	parse.Flags.Refs = !parse.Flags.Refs
	parse.Flags.DetermineWriter()
	err := parse.Parse(args)

	if err != nil {
		fmt.Printf("Got error `%v'\n", err)
	}
}
