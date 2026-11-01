package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"

	id3v2 "github.com/bogem/id3v2/v2"
)

/*
	Usage:
		- Move this script into the folder it should operate on
		- From inside that folder:
			# dry run (tells you what it would rename)
			go run main.go
			# actually rename the files
			go run main.go --exec


	Grok prompts:

	1) Go script to clear one specific field "Grouping" of metadata on a whole
	folder of mp3 files.

	Human finishing touches:

	Replaced majority with the iteration and flag logic from PXL-Renamer since
	I prefer that, then changed the actions and print statements.

*/

// stripGrouping lists files in the current directory, finds those mp3 files
// with "Grouping" metadata column (as seen in iTunes) and clears it. Ignores
// subfolders. If dry run, merely lists files which would be modified.
func stripGrouping(dryRun bool) error {
	entries, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Print("Would Clear `Grouping` Metadata:\n")
	} else {
		fmt.Print("Clearing `Grouping` Metadata:\n")
	}

	re := regexp.MustCompile(`(?i)\.mp3$`)
	for _, entry := range entries {
		if entry.IsDir() {
			continue // Ignore subfolders
		}
		filename := entry.Name()
		if !re.MatchString(filename) {
			continue // Ignore non-mp3 files
		}
		tag, err := id3v2.Open(filename, id3v2.Options{Parse: true})
		if err != nil {
			fmt.Printf("Open error on %s:\n", filename)
			return err
		}
		defer tag.Close()
		frameID := "GRP1" // Modern iTunes Grouping field. NOT "T1T1"
		if frameID == "" {
			return fmt.Errorf("CommonID not found")
		}
		if len(tag.GetFrames(frameID)) == 0 {
			continue // No grouping field, no print
		}
		fmt.Printf("  %s\n", filename)
		if !dryRun {
			tag.DeleteFrames(frameID)
			if err := tag.Save(); err != nil {
				fmt.Printf("Write error on %s:\n", filename)
				return err // Only possible error on write step
			}
		}
	}

	fmt.Print("Done.\n")

	return nil
}

func main() {
	exec := flag.Bool("exec", false, "actually perform the metadata stripping (default is dry-run)")
	flag.Parse()

	if err := stripGrouping(!*exec); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
