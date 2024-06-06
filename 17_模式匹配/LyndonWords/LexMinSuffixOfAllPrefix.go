// LexMinSuffixOfAllPrefix/MinLexSuffixOfAllPrefix
// 所有前缀的字典序最小后缀(的长度).

package main

import "fmt"

func main() {
	s := "ababa"
	res := LexMinSuffixOfAllPrefix(int32(len(s)), func(i int32) int32 { return int32(s[i]) })
	fmt.Println(res)
}

// res[i] 表示前缀 s[0:i) 的字典序最小后缀的长度.
func LexMinSuffixOfAllPrefix(n int32, f func(i int32) int32) []int32 {
	L := NewIncrementalLyndonFactorization()
	for i := int32(0); i < n; i++ {
		L.Add(f(i))
	}
	return L.MinSuffixLen
}

type IncrementalLyndonFactorization struct {
	MinSuffixLen []int32
	i, j, k      int32
	nums         []int32
}

func NewIncrementalLyndonFactorization() *IncrementalLyndonFactorization {
	return &IncrementalLyndonFactorization{MinSuffixLen: []int32{0}}
}

func (ilf *IncrementalLyndonFactorization) Add(c int32) int32 {
	ilf.nums = append(ilf.nums, c)
	for ilf.i < int32(len(ilf.nums)) {
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
		ilf.MinSuffixLen = append(ilf.MinSuffixLen, ilf.i-ilf.k)
	} else {
		ilf.MinSuffixLen = append(ilf.MinSuffixLen, ilf.MinSuffixLen[ilf.k])
	}
	return ilf.MinSuffixLen[ilf.i]
}

func (ilf *IncrementalLyndonFactorization) Factorize() []int32 {
	i := int32(len(ilf.nums))
	var res []int32
	for i > 0 {
		res = append(res, i)
		i -= ilf.MinSuffixLen[i]
	}
	res = append(res, 0)
	for i, j := 0, len(res)-1; i < j; i, j = i+1, j-1 {
		res[i], res[j] = res[j], res[i]
	}
	return res
}
