// Program acat copies its standard input to an output file.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/creachadair/atomicfile"
)

var fileMode = flag.String("mode", "0600", "Output file mode")

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %[1]s <output-file>

Copy standard input to the specified file through a temporary file.
In case of error, the original contents of the file, if any, are not
modified; otherwise, the file is replaced in one step by renaming the
temporary file.

Options:
`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatalf("Usage: %s <output-file>", filepath.Base(os.Args[0]))
	}
	mode, err := strconv.ParseInt(*fileMode, 0, 32)
	if err != nil {
		log.Fatalf("Invalid mode %q: %v", *fileMode, err)
	}
	if _, err := atomicfile.WriteAll(flag.Arg(0), os.Stdin, os.FileMode(mode)); err != nil {
		log.Fatalf("Error: %v", err)
	}
}
