package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	input := flag.String("input", ".", "path to Go source file or directory")
	flag.Parse()

	structs, err := Parse(*input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse error: %v\n", err)
		os.Exit(1)
	}

	// write the output to the same directory as the input
	outputPath := *input
	err = Generate(outputPath, structs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "generate error: %v\n", err)
		os.Exit(1)
	}
}
