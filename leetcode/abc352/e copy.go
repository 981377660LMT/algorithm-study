package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

func main() {
	for i := 0; i < 10000; i++ {
		N := rand.Intn(4) + 1
		K := rand.Intn(4) + 1
		if N < K {
			N, K = K, N
		}
		P := make([]int, N)
		for i := 0; i < N; i++ {
			P[i] = i + 1
		}
		rand.Shuffle(N, func(i, j int) { P[i], P[j] = P[j], P[i] })

		res1 := Solve1(N, K, P)
		res2 := Solve2(N, K, P)
		if res1 != res2 {
			fmt.Println(N, K, P)
			fmt.Println(res1, res2)
			panic("err")
		}
	}
	fmt.Println("ok")
}
func Solve1(N, K int, P []int) int {
	P = append([]int{}, P...)
	for i := 0; i < N; i++ {
		P[i]--
	}
	if N == 1 {
		return 0
	}

	if K == 1 {
		return 0
	}
	check := func(mid int) bool {
		finder := NewFinder(N)
		for i := 0; i < N; i++ {
			finder.Insert(i)
		}
		for right := 0; right < N; right++ {
			finder.Erase(P[right])
			if right >= mid {
				finder.Insert(P[right-mid])
			}
			prev := finder.Prev(P[right]) + 1
			next := finder.Next(P[right]) - 1
			if next-prev+1 >= K {
				return true
			}
		}
		return false
	}
	left, right := K, N
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}

	return left - 1
}

// bruteForce
func Solve2(N, M int, P []int) int {
	if M == 1 {
		return 0
	}
	minDiff := N
	for i := 0; i < 1<<uint(N); i++ {
		choose := []int{}
		for j := 0; j < N; j++ {
			if i>>uint(j)&1 == 1 {
				choose = append(choose, j)
			}
		}
		// 是否存在连续k个数

		check := func() bool {
			nums := []int{}
			for _, v := range choose {
				nums = append(nums, P[v])
			}
			sort.Ints(nums)
			dp := 0
			pre := -1
			for _, v := range nums {
				if v-pre == 1 {
					dp++
					if dp >= M {
						return true
					}
				} else {
					dp = 1
				}
				pre = v
			}
			return false
		}

		if check() {
			min_, max_ := choose[0], choose[0]
			for _, v := range choose {
				min_ = min(min_, v)
				max_ = max(max_, v)
			}
			minDiff = min(minDiff, max_-min_)
		}
	}
	return minDiff
}

type Finder struct {
	n, lg int
	seg   [][]int
}

func NewFinder(n int) *Finder {
	res := &Finder{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	for i := 0; i < res.n; i++ {
		res.Insert(i)
	}
	return res
}

func (fs *Finder) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *Finder) Insert(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
}

func (fs *Finder) Erase(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] &= ^(1 << (i & 63))
		if fs.seg[h][i>>6] != 0 {
			break
		}
		i >>= 6
	}
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *Finder) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		if i>>6 == len(fs.seg[h]) {
			break
		}
		d := fs.seg[h][i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *Finder) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *Finder) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = fs.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (fs *Finder) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("Finder{%v}", strings.Join(res, ", "))
}

func (*Finder) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*Finder) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}
