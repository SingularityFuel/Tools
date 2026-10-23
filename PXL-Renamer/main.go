package main

import (
	"fmt"
	"os"
	"regexp"
)

// renamePXLFiles lists files in the current directory, finds those matching
// ^PXL_\d{8}_, and renames them by removing the "PXL_" prefix and inserting
// a space between the first and last four digits.
//
// Example: PXL_20241023_photo.jpg → 2024 1023_photo.jpg
func renamePXLFiles() error {
	re := regexp.MustCompile(`^PXL_(\d{8})_`)

	entries, err := os.ReadDir(".")
	if err != nil {
		return err
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
		newPrefix := digits[:4] + " " + digits[4:]
		// Everything after "PXL_" + 8 digits (starts with the underscore)
		rest := name[len("PXL_")+8:]
		newName := newPrefix + rest

		fmt.Printf("Renaming: %s → %s\n", name, newName)
		if err := os.Rename(name, newName); err != nil {
			return fmt.Errorf("rename %s: %w", name, err)
		}
	}

	return nil
}

func main() {
	if err := renamePXLFiles(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
