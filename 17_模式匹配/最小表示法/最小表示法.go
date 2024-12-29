package main

// 返回字符串最小表示法/最大表示法.
// 返回循环位移后的开始下标.
func findIsomorphic(n int, s func(i int) int, isMin bool) int {
	if n <= 1 {
		return 0
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

	i1, i2, same := 0, 1, 0

	for i1 < n && i2 < n && same < n {
		diff := compare(s((i1+same)%n), s((i2+same)%n))

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
	return res
}

// 返回字典序最小的/最大的后缀开始下标.
// !注意这里不能循环位移.
func findSuffix(n int, s func(i int) int, isMin bool) int {
	if n <= 1 {
		return 0
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

	i1, i2, same := 0, 1, 0
	for i2+same < n {
		diff := compare(s(i1+same), s(i2+same))

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

	return i1
}
