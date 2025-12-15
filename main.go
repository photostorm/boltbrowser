package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/br0xen/boltbrowser/pkg/boltbrowser"
)

func main() {
	args := boltbrowser.DefaultArgs()
	files := parseArgs(&args)
	err := boltbrowser.Main(args, files)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

func parseArgs(args *boltbrowser.Args) (databaseFiles []string) {
	var err error
	if len(os.Args) == 1 {
		printUsage(nil)
	}
	parms := os.Args[1:]
	for i := range parms {
		// All 'option' arguments start with "-"
		if !strings.HasPrefix(parms[i], "-") {
			databaseFiles = append(databaseFiles, parms[i])
			continue
		}
		if strings.Contains(parms[i], "=") {
			// Key/Value pair Arguments
			pts := strings.Split(parms[i], "=")
			key, val := pts[0], pts[1]
			switch key {
			case "-timeout":
				args.DBOpenTimeout, err = time.ParseDuration(val)
				if err != nil {
					// See if we can successfully parse by adding a 's'
					args.DBOpenTimeout, err = time.ParseDuration(val + "s")
				}
				// If err is still not nil, print usage
				if err != nil {
					printUsage(err)
				}
			case "-readonly", "-ro":
				if val == "true" {
					args.ReadOnly = true
				}
			case "-no-value":
				if val == "true" {
					args.NoValue = true
				}
			case "-help":
				printUsage(nil)
			default:
				printUsage(errors.New("Invalid option"))
			}
		} else {
			// Single-word arguments
			switch parms[i] {
			case "-readonly", "-ro":
				args.ReadOnly = true
			case "-no-value":
				args.NoValue = true
			case "-help":
				printUsage(nil)
			default:
				printUsage(errors.New("Invalid option"))
			}
		}
	}
	return
}

func printUsage(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
	fmt.Fprintf(os.Stderr, "Usage: boltbrowser [OPTIONS] <filename(s)>\nOptions:\n")
	fmt.Fprintf(os.Stderr, "  -timeout=duration\n        DB file open timeout (default 1s)\n")
	fmt.Fprintf(os.Stderr, "  -ro, -readonly   \n        Open the DB in read-only mode\n")
	fmt.Fprintf(os.Stderr, "  -no-value        \n        Do not display a value in left pane\n")
}
