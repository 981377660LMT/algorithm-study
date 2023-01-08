// golang 的正则表达式可以保证线性
package main

import (
	"fmt"
	"regexp"
)

func main() {
	s := "aaaaaaaaaaaaaaaaaaab"
	p := "a*a*a*a*a*a*a*a*a*a*a*a*a*a*a*"
	fmt.Println(isMatch(s, p))
}

func isMatch(s, p string) bool {
	return regexp.MustCompile("^" + p + "$").MatchString(s)
}
