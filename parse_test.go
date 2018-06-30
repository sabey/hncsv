package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

var (
	second_result = []string{
		"17424057",                          // ID
		"https://www.begin.re/",             // URL
		"Reverse Engineering for Beginners", // Title
		"219",          // Points
		"jacquesm",     // User
		"10 hours ago", // Age
		"12",           // Comments
	}
)

func TestParse(t *testing.T) {
	f, err := os.Open("unittest/source.html")
	if err != nil {
		t.Fatalf("failed to open unittest/source.html: \"%s\"\n", err)
		return
	}
	// read our file to memory
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		f.Close()
		t.Fatalf("failed to readall unittest/source.html: \"%s\"\n", err)
		return
	}
	// close file
	f.Close()
	results, err := Parse(bs)
	if err != nil {
		f.Close()
		t.Fatalf("failed to parse unittest/source.html: \"%s\"\n", err)
		return
	}
	// success!!!
	// print parsed as our unittest to log
	// we don't want to print our output to stdout as to not confuse the usecase of this library
	log.Println("parsed unittest/source.html - results:")

	bsr, _ := json.MarshalIndent(results, "", "\t")
	log.Printf("\"%s\"\n", bsr)

	// compare our second result
	if !reflect.DeepEqual(
		results[1].GetCSVFields(),
		second_result,
	) {
		// deep equal doesn't match
		t.Fatal("deep equal for our second result doesn't match")
	}
}
