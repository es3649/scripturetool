/**
 * scripturetool - a command line tool for using the scriptures
 *
 *
 */

package main

import (
	"fmt"
	"os"

	"github.com/es3649/scripturetool/internal/cmd"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Printf("Encountered an error: %v\n", err)
		os.Exit(1)
	}
}
