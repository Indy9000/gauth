package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.Handle("/", http.FileServer(http.Dir("src/ui/")))
	fmt.Printf("Listening on 80")
	http.ListenAndServe(":80", nil)
}
