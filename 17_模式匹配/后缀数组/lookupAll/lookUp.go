package main

import (
	"fmt"
	"index/suffixarray"
)

func main() {
	s := "banana"
	sa := suffixarray.New([]byte(s))
	fmt.Println(sa.Lookup([]byte("a"), -1)) // [5 3 1]
}
