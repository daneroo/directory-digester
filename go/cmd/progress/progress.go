package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// Write a message to stdout as
	fmt.Println("This is a message on stdout")

	// Wait for 2 seconds
	time.Sleep(2 * time.Second)

	// Move the cursor up one line and clear to the end of the line on stderr
	fmt.Fprintf(os.Stderr, "\033[1A\033[K")

	// Write an error message to stderr
	fmt.Fprintf(os.Stderr, "This is an error message on stderr")

	// Wait for 2 seconds
	time.Sleep(2 * time.Second)

	// Move the cursor up one line and clear to the end of the line on stdout
	fmt.Fprintf(os.Stdout, "\033[1A\033[K")

	// Write another message to stdout
	fmt.Println("This is another message on stdout")

	// Wait for 2 seconds
	time.Sleep(2 * time.Second)

	// Move the cursor up one line and clear to the end of the line on stderr
	fmt.Fprintf(os.Stderr, "\033[1A\033[K")

	// Write another error message to stderr
	fmt.Fprintf(os.Stderr, "This is another error message on stderr")

	// Wait for 2 seconds
	time.Sleep(2 * time.Second)
}
