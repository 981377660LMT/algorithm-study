// Map 实现动态开点线段树.

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	seg := NewDynamicSegmentTreeMap(1e9)
	seg.Set(3, 5)
	seg.Set(400000000, 6)
	fmt.Println(seg.QueryAll())
}

type E = int

func (*DynamicSegmentTreeMap) e() E { return 0 }
func (*DynamicSegmentTreeMap) op(a, b E) E {
	return a + b
}

type DynamicSegmentTreeMap struct {
	n, size, log int
	data         map[int]E
}

func NewDynamicSegmentTreeMap(n int) *DynamicSegmentTreeMap {
	res := &DynamicSegmentTreeMap{n: n, data: make(map[int]E)}
	log := bits.Len(uint(n))
	size := 1 << log
	res.size, res.log = size, log
	return res
}

func (seg *DynamicSegmentTreeMap) Get(index int) E {
	return seg._get(index + seg.size)
}

func (seg *DynamicSegmentTreeMap) Set(index int, v E) {
	index += seg.size
	seg.data[index] = v
	for i := 0; i < seg.log; i++ {
		index >>= 1
		seg.data[index] = seg.op(seg._get(index<<1), seg._get(index<<1|1))
	}
}

func (seg *DynamicSegmentTreeMap) Query(start, end int) E {
	start += seg.size
	end += seg.size
	lres, rres := seg.e(), seg.e()
	for start < end {
		if start&1 == 1 {
			lres = seg.op(lres, seg._get(start))
			start++
		}
		if end&1 == 1 {
			end--
			rres = seg.op(seg._get(end), rres)
		}
		start >>= 1
		end >>= 1
	}
	return seg.op(lres, rres)
}

func (seg *DynamicSegmentTreeMap) QueryAll() E {
	return seg._get(1)
}

func (seg *DynamicSegmentTreeMap) MaxRight(start int, f func(E) bool) int {
	if start >= seg.n {
		return seg.n
	}
	start += seg.size
	s := seg.e()
	for {
		for start&1 == 0 {
			start >>= 1
		}
		if !f(seg.op(s, seg._get(start))) {
			for start < seg.size {
				start <<= 1
				if f(seg.op(s, seg._get(start))) {
					s = seg.op(s, seg._get(start))
					start |= 1
				}
			}
			return start - seg.size
		}
		s = seg.op(s, seg._get(start))
		start++
		if start&(-start) == start {
			break
		}
	}
	return seg.n
}

func (seg *DynamicSegmentTreeMap) MinLeft(end int, f func(E) bool) int {
	if end <= 0 {
		return 0
	}
	end += seg.size
	s := seg.e()
	for {
		end--
		for end > 1 && end&1 == 1 {
			end >>= 1
		}
		if !f(seg.op(seg._get(end), s)) {
			for end < seg.size {
				end = end<<1 | 1
				if f(seg.op(seg._get(end), s)) {
					s = seg.op(seg._get(end), s)
					end ^= 1
				}
			}
			return end + 1 - seg.size
		}
		s = seg.op(seg._get(end), s)
		if end&(-end) == end {
			break
		}
	}
	return 0
}

func (seg *DynamicSegmentTreeMap) _get(index int) E {
	if v, ok := seg.data[index]; ok {
		return v
	} else {
		return seg.e()
	}
}
