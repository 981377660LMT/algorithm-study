package main

import (
	"sort"
)

const MOD int = 1e9 + 7

type Enumeration struct {
	fac, ifac, inv []int
	mod            int
}

func NewEnumeration(initSize, mod int) *Enumeration {
	res := &Enumeration{
		fac:  make([]int, 1, initSize+1),
		ifac: make([]int, 1, initSize+1),
		inv:  make([]int, 1, initSize+1),
		mod:  mod,
	}
	res.fac[0] = 1
	res.ifac[0] = 1
	res.inv[0] = 1
	res.expand(initSize)
	return res
}

func (e *Enumeration) Fac(k int) int {
	e.expand(k)
	return e.fac[k]
}

func (e *Enumeration) Ifac(k int) int {
	e.expand(k)
	return e.ifac[k]
}

func (e *Enumeration) Inv(k int) int {
	e.expand(k)
	return e.inv[k]
}

func (e *Enumeration) C(n, k int) int {
	if n < 0 || k < 0 || n < k {
		return 0
	}
	return e.Fac(n) * e.Ifac(k) % e.mod * e.Ifac(n-k) % e.mod
}

func (e *Enumeration) expand(size int) {
	if size < 0 {
		return
	}
	if len(e.fac) > size {
		return
	}
	for i := len(e.fac); i <= size; i++ {
		e.fac = append(e.fac, e.fac[i-1]*i%e.mod)
	}
	if len(e.ifac) <= size {
		e.ifac = make([]int, size+1)
		e.ifac[size] = Pow(e.fac[size], e.mod-2, e.mod)
		for i := size - 1; i >= 0; i-- {
			e.ifac[i] = e.ifac[i+1] * (i + 1) % e.mod
		}
	}
	if len(e.inv) <= size {
		e.inv = make([]int, size+1)
		e.inv[0] = 1
		for i := 1; i <= size; i++ {
			e.inv[i] = e.ifac[i] * e.fac[i-1] % e.mod
		}
	}
}

func Pow(base, exp, mod int) int {
	base %= mod
	if base < 0 {
		base += mod
	}
	res := 1
	for exp > 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

func subsequencesWithMiddleMode(nums []int) int {
	n := len(nums)
	if n < 5 {
		return 0
	}

	elementSet := make(map[int]struct{})
	for _, num := range nums {
		elementSet[num] = struct{}{}
	}
	uniqueElements := make([]int, 0, len(elementSet))
	for num := range elementSet {
		uniqueElements = append(uniqueElements, num)
	}
	sort.Ints(uniqueElements)
	k := len(uniqueElements)
	elementToID := make(map[int]int)
	for idx, num := range uniqueElements {
		elementToID[num] = idx
	}

	prefixCount := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		prefixCount[i] = make([]int, k)
	}
	for i := 1; i <= n; i++ {
		copy(prefixCount[i], prefixCount[i-1])
		id := elementToID[nums[i-1]]
		prefixCount[i][id]++
	}

	suffixCount := make([][]int, n+1)
	for i := 0; i <= n; i++ {
		suffixCount[i] = make([]int, k)
	}
	for i := n - 1; i >= 0; i-- {
		copy(suffixCount[i], suffixCount[i+1])
		id := elementToID[nums[i]]
		suffixCount[i][id]++
	}

	overlapNonM := make([]int, n)
	for i := 0; i < n; i++ {
		m := nums[i]
		id_m := elementToID[m]
		count := 0
		for y := 0; y < k; y++ {
			if y == id_m {
				continue
			}
			if prefixCount[i][y] > 0 && suffixCount[i+1][y] > 0 {
				count++
			}
		}
		overlapNonM[i] = count
	}

	U_left := make([]int, n)
	U_right := make([]int, n)
	for i := 0; i < n; i++ {
		m := nums[i]
		id_m := elementToID[m]
		count_unique_left := 0
		for y := 0; y < k; y++ {
			if y == id_m {
				continue
			}
			if prefixCount[i][y] > 0 {
				count_unique_left++
			}
		}
		U_left[i] = count_unique_left

		count_unique_right := 0
		for y := 0; y < k; y++ {
			if y == id_m {
				continue
			}
			if suffixCount[i+1][y] > 0 {
				count_unique_right++
			}
		}
		U_right[i] = count_unique_right
	}

	E := NewEnumeration(n, MOD)
	ans := 0

	for i := 0; i < n; i++ {
		m := nums[i]
		id_m := elementToID[m]

		cnt_m_left := prefixCount[i][id_m]
		cnt_m_right := suffixCount[i+1][id_m]

		l_non_m_unique := U_left[i]
		r_non_m_unique := U_right[i]

		sum_overlap := overlapNonM[i]

		ways_x1_part1 := (E.C(cnt_m_left, 1) * E.C(l_non_m_unique, 1) % MOD * E.C(r_non_m_unique, 2) % MOD) - (E.C(cnt_m_left, 1) * sum_overlap % MOD)
		if ways_x1_part1 < 0 {
			ways_x1_part1 += MOD
		}

		ways_x1_part2 := (E.C(cnt_m_right, 1) * E.C(r_non_m_unique, 1) % MOD * E.C(l_non_m_unique, 2) % MOD) - (E.C(cnt_m_right, 1) * sum_overlap % MOD)
		if ways_x1_part2 < 0 {
			ways_x1_part2 += MOD
		}

		ways_x1 := (ways_x1_part1 + ways_x1_part2) % MOD

		part1_inner := (E.C(l_non_m_unique, 1) * E.C(r_non_m_unique, 1) % MOD) - sum_overlap
		if part1_inner < 0 {
			part1_inner += MOD
		}
		ways_x2_part1 := E.C(cnt_m_left, 1) * E.C(cnt_m_right, 1) % MOD * part1_inner % MOD

		ways_x2_part2 := E.C(cnt_m_left, 2) * E.C(r_non_m_unique, 2) % MOD

		ways_x2_part3 := E.C(cnt_m_right, 2) * E.C(l_non_m_unique, 2) % MOD

		ways_x2 := (ways_x2_part1 + ways_x2_part2 + ways_x2_part3) % MOD

		total := (ways_x1 + ways_x2) % MOD
		ans = (ans + total) % MOD
	}

	return ans
}
