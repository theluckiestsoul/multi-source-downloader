// Package main defines the entry point of the program.
package main

import (
	"bufio" // For reading lines
	"fmt"   // For printing to standard output

	// For logging errors
	"os" // For accessing the file system
)

// readFile opens a file and prints its contents line by line.
func ReadFile(filename string) error {
	// Open the file for reading
	file, err := os.Open(filename)
	// If an error occurs, return it to the caller
	if err != nil {
		return err
	}
	// Close the file when the function returns
	defer file.Close()

	// Create a new scanner to read the file line by line
	scanner := bufio.NewScanner(file)
	// Loop through all the lines in the file
	for scanner.Scan() {
		// Print each line to the standard output
		fmt.Println(scanner.Text())
	}

	// Check for errors that occurred during scanning
	if err := scanner.Err(); err != nil {
		return err
	}

	// Return nil if no errors occurred
	return nil
}
