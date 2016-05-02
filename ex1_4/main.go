// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 10.
//!+

// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	lineRefs := make(map[string][]string)
	files := os.Args[1:]
	if len(files) == 0 {
		lookForDuplicateLines(os.Stdin, "stdin", lineRefs)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			lookForDuplicateLines(f, arg, lineRefs)
			f.Close()
		}
	}

	for key, value := range lineRefs {
		if len(value) > 1 {
			fmt.Println(key, value)
		}
	}
}

func lookForDuplicateLines(f *os.File, filename string, lineRefs map[string][]string) {
	input := bufio.NewScanner(f)
	line := 1
	for input.Scan() {
		lineRefs[input.Text()] = append(lineRefs[input.Text()], fmt.Sprintf("%s:%d", filename, line))
		line++
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!-
