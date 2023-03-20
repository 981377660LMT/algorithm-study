package main

import "math/rand"

func sumScores(s string) int64 {
	n := len(s)
	hasher := NewRollingHashField()
	hasher.Build(s)
	countPre := func(curLen, start int) int {
		left, right := 1, curLen
		for left <= right {
			mid := (left + right) >> 1
			hash00 := hasher.Query(start, start+mid)
			hash10 := hasher.Query(0, mid)
			if hash00 == hash10 {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}

		return right
	}

	res := 0
	for i := 1; i < n+1; i++ {
		count := countPre(i, n-i)
		res += count
	}
	return int64(res)
}

// DIY
func (*RollingHashField) add(a, b uint) uint { return a + b }
func (*RollingHashField) sub(a, b uint) uint { return a - b }
func (*RollingHashField) mul(a, b uint) uint { return a * b }

type RollingHashField struct {
	base   uint
	power  []uint
	hashed []uint
}

func NewRollingHashField() *RollingHashField {
	return &RollingHashField{
		base:  uint(rand.Uint64()),
		power: []uint{1},
	}
}

func (r *RollingHashField) Build(s string) (hashed []uint) {
	sz := len(s)
	hashed = make([]uint, sz+1)
	for i := 0; i < sz; i++ {
		hashed[i+1] = r.add(r.mul(hashed[i], r.base), uint(s[i]))
	}
	r.hashed = hashed
	return hashed
}

// [start, end)
func (r *RollingHashField) Query(start, end int) uint {
	r.expand(end - start)
	return r.sub(r.hashed[end], r.mul(r.hashed[start], r.power[end-start]))
}

func (r *RollingHashField) QueryWithHashed(hashed []uint, start, end int) uint {
	r.expand(end - start)
	return r.sub(hashed[end], r.mul(hashed[start], r.power[end-start]))
}

func (r *RollingHashField) Combine(h1, h2 uint, h2len int) uint {
	r.expand(h2len)
	return r.add(r.mul(h1, r.power[h2len]), h2)
}

// 两个字符串的最长公共前缀长度
func (r *RollingHashField) LCP(hashed1 []uint, start1, end1 int, hashed2 []uint, start2, end2 int) int {
	len1 := end1 - start1
	len2 := end2 - start2
	len := min(len1, len2)
	low := 0
	high := len + 1
	for high-low > 1 {
		mid := (low + high) / 2
		if r.QueryWithHashed(hashed1, start1, start1+mid) == r.QueryWithHashed(hashed2, start2, start2+mid) {
			low = mid
		} else {
			high = mid
		}
	}
	return low
}

func (r *RollingHashField) expand(sz int) {
	if len(r.power) < sz+1 {
		preSz := len(r.power)
		r.power = append(r.power, make([]uint, sz+1-preSz)...)
		for i := preSz - 1; i < sz; i++ {
			r.power[i+1] = r.power[i] * r.base
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
