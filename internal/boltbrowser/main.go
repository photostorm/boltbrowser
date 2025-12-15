package boltbrowser

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
	"go.etcd.io/bbolt"
)

var ProgramName = "boltbrowser"
var VersionNum = 2.0

var AppArgs Args
var databaseFiles []string
var db *bbolt.DB
var memBolt *BoltDB

var currentFilename string

const DefaultDBOpenTimeout = time.Second

type Args struct {
	DBOpenTimeout time.Duration
	ReadOnly      bool
	NoValue       bool
}

func DefaultArgs() Args {
	return Args{
		DBOpenTimeout: DefaultDBOpenTimeout,
	}
}

func Main(args Args, files []string) error {
	// Set the global args. This is done to convert main package into a library with minimal changes.
	AppArgs, databaseFiles = args, files

	var err error
	err = termbox.Init()
	if err != nil {
		return err
	}
	defer termbox.Close()
	style := defaultStyle()
	termbox.SetOutputMode(termbox.Output256)

	for _, databaseFile := range databaseFiles {
		currentFilename = databaseFile
		db, err = bbolt.Open(databaseFile, 0600, &bbolt.Options{Timeout: AppArgs.DBOpenTimeout})
		if err == bbolt.ErrTimeout {
			termbox.Close()
			fmt.Printf("File %s is locked. Make sure it's not used by another app and try again\n", databaseFile)
			os.Exit(1)
		} else if err != nil {
			if len(databaseFiles) > 1 {
				mainLoop(nil, style)
				continue
			} else {
				termbox.Close()
				fmt.Printf("Error reading file: %q\n", err.Error())
				os.Exit(1)
			}
		}

		// First things first, load the database into memory
		memBolt.refreshDatabase()
		if AppArgs.ReadOnly {
			// If we're opening it in readonly mode, close it now
			db.Close()
		}

		// Kick off the UI loop
		mainLoop(memBolt, style)
		defer db.Close()
	}
	return nil
}
