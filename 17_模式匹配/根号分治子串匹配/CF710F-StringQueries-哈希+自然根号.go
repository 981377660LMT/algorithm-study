// https://www.cnblogs.com/TianMeng-hyl/p/14989441.html
// https://www.cnblogs.com/Dfkuaid-210/p/bit_divide.html

package main

import (
	"bufio"
	"fmt"
	"os"
)

// String Set Queries
// https://www.luogu.com.cn/problem/CF710F
// 1 s : 在数据结构中插入 s
// 2 s : 在数据结构中删除 s
// 3 s : 查询集合中的所有字符串在给出的模板串中出现的次数
//
// 哈希+自然根号：
// !由于字符串总长度不超过3e5，因此最多有不超过根号种不同的长度。
// 所以遍历子串时只需要遍历根号种长度的子串即可。
// 维护根号个set即可.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)

	R := NewSimpleHash(BASE, MOD)
	hashGroupByLen := make(map[int]map[uint]int) // !每个长度对应一个counter

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op int
		var s string
		fmt.Fscan(in, &op, &s)
		if op == 1 {
			hash := GetHash(s, BASE, MOD)
			len_ := len(s)
			if _, ok := hashGroupByLen[len_]; !ok {
				hashGroupByLen[len_] = make(map[uint]int)
			}
			hashGroupByLen[len_][hash]++
		} else if op == 2 {
			hash := GetHash(s, BASE, MOD)
			len_ := len(s)
			if _, ok := hashGroupByLen[len_]; !ok {
				continue
			}
			inner := hashGroupByLen[len_]
			if _, ok := inner[hash]; !ok {
				continue
			}
			inner[hash]--
		} else {
			res := 0
			hashTable := R.Build(s)
			// !遍历每种长度的子串
			for len_, group := range hashGroupByLen {
				for start := 0; start+len_ <= len(s); start++ {
					hash := R.Query(hashTable, start, start+len_)
					res += group[hash]
				}
			}
			fmt.Fprintln(out, res)
			out.Flush() // 强制在线，需要刷新缓冲区
		}
	}
}

const BASE uint = 13331
const MOD uint = 1e9 + 7

type S = string

func GetHash(s S, base uint, mod uint) uint {
	if len(s) == 0 {
		return 0
	}
	res := uint(0)
	for i := 0; i < len(s); i++ {
		res = (res*base + uint(s[i])) % mod
	}
	return res
}

type SimpleHash struct {
	base  uint
	mod   uint
	power []uint
}

// 131/13331/1713302033171(回文素数)
func NewSimpleHash(base uint, mod uint) *SimpleHash {
	return &SimpleHash{
		base:  base,
		mod:   mod,
		power: []uint{1},
	}
}

func (r *SimpleHash) Build(s S) (hashTable []uint) {
	sz := len(s)
	hashTable = make([]uint, sz+1)
	for i := 0; i < sz; i++ {
		hashTable[i+1] = (hashTable[i]*r.base + uint(s[i])) % r.mod
	}
	return hashTable
}

func (r *SimpleHash) Query(sTable []uint, start, end int) uint {
	r.expand(end - start)
	return (r.mod + sTable[end] - sTable[start]*r.power[end-start]%r.mod) % r.mod
}

func (r *SimpleHash) expand(sz int) {
	if len(r.power) < sz+1 {
		preSz := len(r.power)
		r.power = append(r.power, make([]uint, sz+1-preSz)...)
		for i := preSz - 1; i < sz; i++ {
			r.power[i+1] = (r.power[i] * r.base) % r.mod
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
