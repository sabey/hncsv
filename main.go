package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
)

const (
	hn_yc       = `https://news.ycombinator.com/`
	MAX_RESULTS = 20
)

var (
	flag_csv = flag.String("csv", "", "Path to CSV File. Empty path defaults to stdout.")
)

func main() {
	// parse flags
	if !flag.Parsed() {
		flag.Parse()
	}
	// log prints to stderr, fmt prints to stdout
	// it's safe to print to log and pipe stdout to a file
	//
	// we're going to set a default writing interface
	// this way we don't have to repeat any code
	// we could use fmt for stdout but this is easier
	// we're going to use a WriteCloser just so that we can close our file if we open a csv file
	// we don't have to close our stdout by default, our runtime handles that
	var w io.WriteCloser
	// we're going to set w to os.Stdout by default
	w = os.Stdout
	// use a csv file???
	if *flag_csv != "" {
		var err error
		// overwrite our io writer if we manage to open our file for write access
		w, err = OpenFile(*flag_csv)
		if err != nil {
			// failed to open file!
			log.Fatalf("failed to open file for writing: \"%s\"\n", err)
			return
		}
		// close our file descriptor once we finish
		defer w.Close()
	}
	// get index
	body, err := GetURL(hn_yc)
	if err != nil {
		// failed to get URL
		log.Fatalf("failed to GET URL: \"%s\"\n", err)
		return
	}
	// parse
	results, err := Parse(body)
	if err != nil {
		// failed to get URL
		log.Fatalf("failed to parse Body: \"%s\"\n", err)
		return
	}
	// parsed!!!
	// create a csv writer
	// the reason for using a csv writer and not doing this ourselves is that this can easily handle quoted delimiters, like a,b,"1,2,3",c,d
	csv_w := csv.NewWriter(w)
	// write each individual row
	for i := 0; i < len(results) && i < MAX_RESULTS; i++ {
		if err := csv_w.Write(results[i].GetCSVFields()); err != nil {
			log.Fatalf("error writing row[%d] to csv: \"%s\"\n", i, err)
			return
		}
	}
	// Write any buffered data to the underlying writer
	csv_w.Flush()
	// check if any errors occurred during write or flush
	if err := csv_w.Error(); err != nil {
		log.Fatalf("failed to write to csv: \"%s\"\n", err)
		return
	}
	// success!
}
