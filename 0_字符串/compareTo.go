// 字符串字典序比较

package main

import "fmt"

func main() {
	fmt.Println(CompareTo("abc", "abcdd"))
}

func CompareTo(s1, s2 string) int {
	len1 := len(s1)
	len2 := len(s2)
	lim := min(len1, len2)
	ords1 := []rune(s1)
	ords2 := []rune(s2)
	for i := 0; i < lim; i++ {
		a, b := ords1[i], ords2[i]
		if a != b {
			return int(a) - int(b)
		}
	}
	return len1 - len2
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
