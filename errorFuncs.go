/*
	errorFuncs.go
		Provides custom error-handling functions for the API
*/

package main

import (
	"fmt"
	"log"
	"net/http"
)

// Return an error without killing the program
func printErrorMessage(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "The following error occurred: %v", err)
}

// Something really bad happened and the program needs to end
func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
