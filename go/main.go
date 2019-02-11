package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	log.Printf("main function")
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello")
}
