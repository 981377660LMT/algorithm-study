package main

import "fmt"

func main() {
	fmt.Println(string(Lowercase('A')))
	fmt.Println(string(Uppercase('a')))
	fmt.Println(string(Swapcase('a')))
	fmt.Println(string(Swapcase('A')))
}

func Lowercase(b byte) byte {
	if 'A' <= b && b <= 'Z' {
		return b | 32
	}
	return b
}

func Uppercase(b byte) byte {
	if 'a' <= b && b <= 'z' {
		return b &^ 32
	}
	return b
}

func Swapcase(b byte) byte {
	if ('A' <= b && b <= 'Z') || ('a' <= b && b <= 'z') {
		return b ^ 32
	}
	return b
}
