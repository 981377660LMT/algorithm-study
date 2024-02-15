package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo()
}

// https://judge.yosupo.jp/problem/lyndon_factorization
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	L := NewIncrementalLyndonFactorization()
	for _, c := range s {
		L.Add(int(c))
	}
	res := L.Factorize()
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

type IncrementalLyndonFactorization struct {
	i, j, k      int
	nums         []int
	minSuffixLen []int
}

func NewIncrementalLyndonFactorization() *IncrementalLyndonFactorization {
	return &IncrementalLyndonFactorization{minSuffixLen: []int{0}}
}

func (ilf *IncrementalLyndonFactorization) Add(c int) int {
	ilf.nums = append(ilf.nums, c)
	for ilf.i < len(ilf.nums) {
		if ilf.k == ilf.i {
			ilf.i++
		} else if ilf.nums[ilf.k] == ilf.nums[ilf.i] {
			ilf.k++
			ilf.i++
		} else if ilf.nums[ilf.k] < ilf.nums[ilf.i] {
			ilf.k = ilf.j
			ilf.i++
		} else {
			ilf.j += (ilf.i - ilf.j) / (ilf.i - ilf.k) * (ilf.i - ilf.k)
			ilf.i, ilf.k = ilf.k, ilf.j
		}
	}
	if (ilf.i-ilf.j)%(ilf.i-ilf.k) == 0 {
		ilf.minSuffixLen = append(ilf.minSuffixLen, ilf.i-ilf.k)
	} else {
		ilf.minSuffixLen = append(ilf.minSuffixLen, ilf.minSuffixLen[ilf.k])
	}
	return ilf.minSuffixLen[ilf.i]
}

func (ilf *IncrementalLyndonFactorization) Factorize() []int {
	i := len(ilf.nums)
	var res []int
	for i > 0 {
		res = append(res, i)
		i -= ilf.minSuffixLen[i]
	}
	res = append(res, 0)
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}
