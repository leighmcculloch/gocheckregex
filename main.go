package main // import "4d63.com/gocheckregex"

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	flagPrintHelp := flag.Bool("h", false, "Print help")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: gocheckregex [path] [path] ...\n")
		flag.PrintDefaults()
	}
	flag.Parse()

	if *flagPrintHelp {
		flag.Usage()
		return
	}

	paths := flag.Args()
	if len(paths) == 0 {
		paths = []string{"./..."}
	}

	exitWithError := false

	for _, path := range paths {
		messages, err := checkRegex(path)
		for _, message := range messages {
			fmt.Fprintf(os.Stdout, "%s\n", message)
			exitWithError = true
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			exitWithError = true
		}
	}

	if exitWithError {
		os.Exit(1)
	}
}
