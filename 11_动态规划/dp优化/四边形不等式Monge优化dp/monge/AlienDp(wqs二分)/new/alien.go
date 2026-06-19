// https://atcoder.jp/contests/abc462/tasks/abc462_f

package main

type AlienResult struct{ Val, Cnt int }

func AlienTrick(k int, solve func(penalty int) AlienResult) (res int, penalty int) {
	c0 := solve(0).Cnt
	var lo, hi int
	if c0 <= k {
		lo, hi = -1, 0
		for solve(lo).Cnt < k {
			d := hi - lo
			lo -= 2 * d
			hi -= d
		}
	} else {
		lo, hi = 0, 1
		for solve(hi).Cnt > k {
			d := hi - lo
			lo += d
			hi += 2 * d
		}
	}
	return alien(k, lo, hi, solve)
}

func alien(k int, lo int, hi int, solve func(penalty int) AlienResult) (res int, penalty int) {
	for lo+1 < hi {
		mid := lo + (hi-lo)/2
		res := solve(mid)
		if res.Cnt <= k {
			hi = mid
		} else {
			lo = mid
		}
	}
	tmp := solve(hi)
	res = tmp.Val - hi*k
	penalty = hi
	return
}
