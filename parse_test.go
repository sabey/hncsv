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
	// we're using the second result
	source_result = []string{
		"17424057",                          // ID
		"https://www.begin.re/",             // URL
		"Reverse Engineering for Beginners", // Title
		"219",          // Points
		"jacquesm",     // User
		"10 hours ago", // Age
		"12",           // Comments
		"false",        // IsJob
	}
	// we're using the 7th (job posting) result
	source2_result = []string{
		"17432428",                                                      // ID
		"http://www.irisonboard.com/careers/",                           // URL
		"Iris Automation Is Hiring a C++ Engineer – Self-Flying Drones", // Title
		"0",             // Points
		"",              // User
		"5 minutes ago", // Age
		"0",             // Comments
		"true",          // IsJob
	}
)

func TestParseSource1(t *testing.T) {
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
		source_result,
	) {
		// deep equal doesn't match
		t.Fatal("deep equal for our second result doesn't match")
	}
}
func TestParseSource2(t *testing.T) {
	f, err := os.Open("unittest/source2.html")
	if err != nil {
		t.Fatalf("failed to open unittest/source2.html: \"%s\"\n", err)
		return
	}
	// read our file to memory
	bs, err := ioutil.ReadAll(f)
	if err != nil {
		f.Close()
		t.Fatalf("failed to readall unittest/source2.html: \"%s\"\n", err)
		return
	}
	// close file
	f.Close()
	results, err := Parse(bs)
	if err != nil {
		f.Close()
		t.Fatalf("failed to parse unittest/source2.html: \"%s\"\n", err)
		return
	}
	// success!!!
	// print parsed as our unittest to log
	// we don't want to print our output to stdout as to not confuse the usecase of this library
	log.Println("parsed unittest/source2.html - results:")

	bsr, _ := json.MarshalIndent(results, "", "\t")
	log.Printf("\"%s\"\n", bsr)

	// compare our second result
	if !reflect.DeepEqual(
		results[6].GetCSVFields(),
		source2_result,
	) {
		// deep equal doesn't match
		t.Fatal("deep equal for our second result doesn't match")
	}
}
