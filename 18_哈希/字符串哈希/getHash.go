package main

import "fmt"

func main() {
	s := "abc"
	hash1, hash2 := getHashString(s)
	fmt.Println(hash1, hash2)
}

// 131/13331/1713302033171(回文素数)
const BASE1 uint = 13331
const BASE2 uint = 1713302033171

// 字符串的哈希值.
func getHashString(s string) (hash1, hash2 uint) {
	if len(s) == 0 {
		return
	}
	for i := 0; i < len(s); i++ {
		hash1 = hash1*BASE1 + uint(s[i])
		hash2 = hash2*BASE2 + uint(s[i])
	}
	return
}

// 正整数切片的哈希值.
func getHashIntSlice(nums []int) (hash1, hash2 uint) {
	if len(nums) == 0 {
		return
	}
	for i := 0; i < len(nums); i++ {
		hash1 = hash1*BASE1 + uint(nums[i])
		hash2 = hash2*BASE2 + uint(nums[i])
	}
	return
}
