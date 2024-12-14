// Returning Pointers Made Easy with Generics.

package main

import (
	"fmt"
	"time"
)

func main() {
	timePtr := Box(time.Now())
	intPtr := Box(42)
	stringPtr := Box("hello")
	fmt.Println(*timePtr, *intPtr, *stringPtr)
}

func Box[T any](v T) *T { return &v }
