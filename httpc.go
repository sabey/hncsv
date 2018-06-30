package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	ERR_HTTP_GET_FAILED  = fmt.Errorf("HTTP GET Failed")
	ERR_HTTP_BODY_FAILED = fmt.Errorf("HTTP Failed To Read Body!")
)

func GetURL(
	url string,
) (
	[]byte,
	error,
) {
	// I'm going to just use this as a wrapper for the built in HTTP library
	// I'll be returning the body read into a byte slice or an error
	r, err := http.Get(url)
	if err != nil {
		// I'm currently shadowing the error since I may only be interested if the GET failed
		return nil, ERR_HTTP_GET_FAILED
	}
	// defer close
	defer r.Body.Close()
	// read the body
	bs, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// failed to read body
		return nil, ERR_HTTP_BODY_FAILED
	}
	// ready body
	return bs, nil
}
