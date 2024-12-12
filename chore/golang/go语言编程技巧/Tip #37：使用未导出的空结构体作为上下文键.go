package main

import (
	"context"
	"fmt"
)

func main() {
	s1 := symbol1{}
	s2 := symbol2{}
	ctx := context.WithValue(context.Background(), s1, "foo")

	fmt.Println(ctx.Value(s1)) // "foo"
	fmt.Println(ctx.Value(s2)) // nil
}

type symbol1 struct{}
type symbol2 struct{}
