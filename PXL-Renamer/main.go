package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
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

	1) Regex that matches strings starting with PXL_00000000_ with any 8
	digits in place of the zeroes

	2) Go function that lists all files in the current folder, detects which
	filenames match the regex above, and for each matching file, renames the
	file with the PXL_ prefix removed and in its place the 8 digits with a
	space separating the first and last four

	Human finishing touches:

	1) Add this section

	2) Modify function (and doc comment) to change final matched underscore
	to a space

	3) Use/test locally as one of the steps for cleaning up a bunch of photos
	downloaded using google takeout on my google photos. Debug a bit.

	4) Add dry run by default and more informative print statements.

*/

// renamePXLFiles lists files in the current directory, finds those matching
// ^PXL_\d{8}_, and renames them by removing the "PXL_" prefix and inserting
// spaces between the first and last four digits and after the last four.
//
// Example: PXL_20241023_photo.jpg -> 2024 1023 photo.jpg
func renamePXLFiles(dryRun bool) error {
	re := regexp.MustCompile(`^PXL_(\d{8})_`)

	entries, err := os.ReadDir(".")
	if err != nil {
		return err
	}

	if dryRun {
		fmt.Print("Would rename:\n")
	} else {
		fmt.Print("Renaming:\n")
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		matches := re.FindStringSubmatch(name)
		if matches == nil {
			continue
		}

		digits := matches[1] // the captured 8 digits
		newPrefix := fmt.Sprintf("%s %s ", digits[:4], digits[4:8])
		// Everything after "PXL_" + 8 digits (starts with the underscore)
		rest := name[len("PXL_")+9:]
		newName := newPrefix + rest

		if dryRun {
			fmt.Printf("  %s -> %s\n", name, newName)
		} else {
			fmt.Printf("  %s -> %s\n", name, newName)
			if err := os.Rename(name, newName); err != nil {
				fmt.Printf("    ERROR ON %s: %s\n", name, err.Error())
				fmt.Print("    continue...\n")
			}
		}
	}

	fmt.Print("Done.\n")

	return nil
}

func main() {
	exec := flag.Bool("exec", false, "actually perform the renames (default is dry-run)")
	flag.Parse()

	if err := renamePXLFiles(!*exec); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
