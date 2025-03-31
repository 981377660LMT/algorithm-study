package main

func lcp(s1, s2 string) int {
	n := min(len(s1), len(s2))
	for i := range n {
		if s1[i] != s2[i] {
			return i
		}
	}
	return n
}
