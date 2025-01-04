package main

import (
	"strings"
)

func hasMatch(s string, p string) bool {
	index := strings.IndexByte(p, '*')
	if index == -1 {
		return false
	}

	n := len(s)
	left, right := p[:index], p[index+1:]
	for start := 0; start < n; start++ {
		for end := start; end < n; end++ {
			sub := s[start : end+1]
			if len(sub) < len(left)+len(right) {
				continue
			}
			if strings.HasPrefix(sub, left) && strings.HasSuffix(sub, right) {
				return true
			}
		}
	}

	return false
}
