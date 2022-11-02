package main

import "fmt"

type Metadata struct {
	age   int
	speed int
	pos   float64
}

type Fish struct {
	id   int
	name string
	Metadata
}

func main() {
	fish := make([]Fish, 5)
	fmt.Println(fish) // [{0  {0 0 0}} {0  {0 0 0}} {0  {0 0 0}} {0  {0 0 0}} {0  {0 0 0}}]
}
