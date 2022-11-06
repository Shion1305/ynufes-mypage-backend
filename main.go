package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Welcome to ynufes-mypage-backend")
	})
	http.HandleFunc("/hello/", func(w http.ResponseWriter, r *http.Request) {
		name := r.URL.Path[len("/hello/")]
		fmt.Fprintf(w, "Hello %s\n", name)
	})
	http.ListenAndServe(":1305", nil)
}
