// https://sotanishy.github.io/cp-library-cpp/data-structure/sqrt_tree.cpp
// 比st表稍快，且节省空间
// O(nloglogn)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	a := make([]int, n)
	for i := range a {
		fmt.Fscan(in, &a[i])
	}
	st := NewSqrtTree(a)
	for q > 0 {
		q--
		var l, r int
		fmt.Fscan(in, &l, &r)
		fmt.Fprintln(out, st.Query(l, r))
	}
}

type S = int

func (st *SqrtTree) op(a, b S) S {
	return min(a, b)
}

type SqrtTree struct {
	layerLg, onLayer []int
	nums             []S
	pref, suf, btwn  [][]S
	n                int
}

func NewSqrtTree(nums []S) *SqrtTree {
	res := &SqrtTree{}
	n := len(nums)
	lg := 0
	for 1<<lg < n {
		lg++
	}

	onLayer := make([]int, lg+1)
	layerLog := []int{}
	nLayer := 0
	for i := lg; i > 1; i = (i + 1) >> 1 {
		onLayer[i] = nLayer
		layerLog = append(layerLog, i)
		nLayer++
	}

	for i := lg - 1; i >= 0; i-- {
		onLayer[i] = max(onLayer[i], onLayer[i+1])
	}

	pref := make([][]S, nLayer)
	suf := make([][]S, nLayer)
	btwn := make([][]S, nLayer)
	for i := range pref {
		pref[i] = make([]S, n)
		suf[i] = make([]S, n)
		btwn[i] = make([]S, 1<<lg)
	}

	for layer := 0; layer < nLayer; layer++ {
		prevBSz := 1 << layerLog[layer]
		bSz := 1 << ((layerLog[layer] + 1) >> 1)
		bCnt := 1 << (layerLog[layer] >> 1)
		for l := 0; l < n; l += prevBSz {
			r := min(l+prevBSz, n)
			for a := l; a < r; a += bSz {
				b := min(a+bSz, r)
				pref[layer][a] = nums[a]
				for i := a + 1; i < b; i++ {
					pref[layer][i] = res.op(pref[layer][i-1], nums[i])
				}
				suf[layer][b-1] = nums[b-1]
				for i := b - 2; i >= a; i-- {
					suf[layer][i] = res.op(nums[i], suf[layer][i+1])
				}
			}
			for i := 0; i < bCnt && l+i*bSz < n; i++ {
				val := suf[layer][l+i*bSz]
				btwn[layer][l+i*bCnt+i] = val
				for j := i + 1; j < bCnt && l+j*bSz < n; j++ {
					val = res.op(val, suf[layer][l+j*bSz])
					btwn[layer][l+i*bCnt+j] = val
				}
			}
		}
	}

	res.layerLg = layerLog
	res.onLayer = onLayer
	res.nums = nums
	res.pref = pref
	res.suf = suf
	res.btwn = btwn
	res.n = n
	return res
}

// 0<=l<=r<=n
func (st *SqrtTree) Query(l, r int) S {
	r--
	if l > r {
		panic(fmt.Sprintf("l>r: %d %d", l, r))
	}
	if l == r {
		return st.nums[l]
	}
	if l+1 == r {
		return st.op(st.nums[l], st.nums[r])
	}
	layer := st.onLayer[32-bits.LeadingZeros32(uint32(l^r))]
	bSz := 1 << ((st.layerLg[layer] + 1) >> 1)
	bCnt := 1 << (st.layerLg[layer] >> 1)
	a := (l >> st.layerLg[layer]) << st.layerLg[layer]
	lBlock := (l-a)/bSz + 1
	rBlock := (r-a)/bSz - 1
	val := st.suf[layer][l]
	if lBlock <= rBlock {
		val = st.op(val, st.btwn[layer][a+lBlock*bCnt+rBlock])
	}
	val = st.op(val, st.pref[layer][r])
	return val
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
func (ds *SqrtTree) MaxRight(left int, check func(e S) bool) int {
	if left == ds.n {
		return ds.n
	}
	ok, ng := left, ds.n+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(ds.Query(left, mid)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
func (ds *SqrtTree) MinLeft(right int, check func(e S) bool) int {
	if right == 0 {
		return 0
	}
	ok, ng := right, -1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(ds.Query(mid, right)) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
