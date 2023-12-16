package main

import "fmt"

func main() {
	s := "abc"
	fmt.Println(GetHash(s, BASE, MOD))
}

// 131/13331/1713302033171(回文素数)
const BASE uint = 13131
const MOD uint = 1000000007

type S = string

func GetHash(s S, base uint, mod uint) uint {
	if len(s) == 0 {
		return 0
	}
	res := uint(0)
	for i := 0; i < len(s); i++ {
		res = (res*base + uint(s[i])) % mod
	}
	return res
}
