package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there!")
}

func main() {
	fmt.Println("Listening on 8081")
    http.HandleFunc("/", handler)
    http.ListenAndServeTLS(":8081", "../../certs/cert.pem", "../../certs/key.pem", nil)
}