package main

import "fmt"

func main() {
	vs := NewVectorSpaceArray([]int32{7, 2, 3, 4})
	fmt.Println(vs.Has(7))
}

type VS [20]int32

func NewVectorSpaceArray(nums []int32) *VS {
	res := VS{}
	for _, num := range nums {
		res.Add(num)
	}
	return &res
}

func (lb *VS) Add(num int32) {
	if num != 0 {
		for i := len(lb) - 1; i >= 0; i-- {
			if num>>i&1 == 1 {
				if lb[i] == 0 {
					lb[i] = num
					break
				}
				num ^= lb[i]
			}
		}
	}
}

// 求xor与所有向量异或的最大值.
func (lb *VS) Max(xor int32) int32 {
	res := xor
	for i := len(lb) - 1; i >= 0; i-- {
		if tmp := res ^ lb[i]; tmp > res {
			res = tmp
		}
	}
	return res
}

// 求xor与所有向量异或的最小值.
func (lb *VS) Min(xorVal int32) int32 {
	res := xorVal
	for i := len(lb) - 1; i >= 0; i-- {
		if tmp := res ^ lb[i]; tmp < res {
			res = tmp
		}
	}
	return res
}

func (lb *VS) Copy() *VS {
	res := &VS{}
	copy(res[:], lb[:])
	return res
}

func (lb *VS) Clear() {
	for i := range lb {
		lb[i] = 0
	}
}

func (lb *VS) Len() int {
	return len(lb)
}

func (lb *VS) ForEach(f func(base int32)) {
	for _, base := range lb {
		f(base)
	}
}

func (lb *VS) Has(v int32) bool {
	for i := len(lb) - 1; i >= 0; i-- {
		if v == 0 {
			break
		}
		v = min32(v, v^lb[i])
	}
	return v == 0
}

// Merge.
func (lb *VS) Or(other *VS) *VS {
	v1, v2 := lb, other
	if v1.Len() < v2.Len() {
		v1, v2 = v2, v1
	}
	res := v1.Copy()
	for _, base := range v2 {
		res.Add(base)
	}
	return res
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a >= b {
		return a
	}
	return b
}
