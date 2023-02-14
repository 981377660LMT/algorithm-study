package main

// 返回字符串最小表示法/最大表示法.
func findIsomorphic(s []int, isMin bool) []int {
	if len(s) <= 1 {
		return s
	}

	compare := func(s1, s2 int) int {
		if s1 == s2 {
			return 0
		}
		if isMin {
			if s1 > s2 {
				return 1
			}
			return -1
		} else {
			if s1 < s2 {
				return 1
			}
			return -1
		}
	}

	n := len(s)
	i1, i2, same := 0, 1, 0

	for i1 < n && i2 < n && same < n {
		diff := compare(s[(i1+same)%n], s[(i2+same)%n])

		if diff == 0 {
			same++
			continue
		} else if diff > 0 {
			i1 += same + 1
		} else {
			i2 += same + 1
		}

		if i1 == i2 {
			i2++
		}

		same = 0
	}

	res := i1
	if i2 < i1 {
		res = i2
	}

	return append(s[res:], s[:res]...)
}
