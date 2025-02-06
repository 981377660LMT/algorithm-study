package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Hello World")
		w.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":8080", nil)

	http.Post("http://localhost:8080", "application/json", nil)
}
