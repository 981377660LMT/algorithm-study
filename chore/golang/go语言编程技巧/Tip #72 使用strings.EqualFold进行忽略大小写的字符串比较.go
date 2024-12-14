package main

import (
	"fmt"
	"strings"
)

func main() {
	s1, s2 := "Go", "go"
	fmt.Println(strings.EqualFold(s1, s2)) // true
}

// 快速路径：快速检查 ASCII 字符，逐个字符查看每个字符。
// - 慢速路径：如果在任何字符串中发现 Unicode 字符，则goto到详细的 Unicode 比较。
