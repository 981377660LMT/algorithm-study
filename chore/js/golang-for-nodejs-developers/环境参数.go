package main

import (
	"fmt"
	"os"
)

func main() {
	key := os.Getenv("API_KEY")

	fmt.Println(key)

	fmt.Println(os.Getenv("HOME"))
}
