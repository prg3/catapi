package main

import (
	"fmt"
	"net/http"
)

func catHandler( *http.Request  ) string {
	fmt.Fprintln("This is in the function")
	return "You asked for a cat from a function\n"
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Sorry, the path %s is invalid\n", r.URL.Path)
	})

	http.HandleFunc("/cat", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, catHandler( r ))
	})

	http.HandleFunc("/history", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "You asked for history\n")
	})

	http.ListenAndServe(":80", nil)
}
