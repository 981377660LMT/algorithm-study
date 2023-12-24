// TODO1: ts版本字符串哈希
// TODO2: 优化哈希类型(type S = string)
// TODO4: dp[i][j]表示两个字符串是否想等

package main

const BASE uint = 13131
const MOD uint = 1000000007

type S = string

func GetHash(s S, base uint) uint {
	if len(s) == 0 {
		return 0
	}
	res := uint(0)
	for i := 0; i < len(s); i++ {
		res = (res*base + uint(s[i]))
	}
	return res
}

// TODO: 优化哈希类型
type RollingHash struct {
	base  uint
	power []uint
}

// 131/13331/1713302033171(回文素数)
func NewRollingHash(base uint) *RollingHash {
	return &RollingHash{
		base:  base,
		power: []uint{1},
	}
}

func (r *RollingHash) Build(s S) (hashTable []uint) {
	sz := len(s)
	hashTable = make([]uint, sz+1)
	for i := 0; i < sz; i++ {
		hashTable[i+1] = hashTable[i]*r.base + uint(s[i])
	}
	return hashTable
}

func (r *RollingHash) Query(sTable []uint, start, end int) uint {
	r.expand(end - start)
	return sTable[end] - sTable[start]*r.power[end-start]
}

func (r *RollingHash) Combine(h1, h2 uint, h2len int) uint {
	r.expand(h2len)
	return h1*r.power[h2len] + h2
}

func (r *RollingHash) AddChar(hash uint, c byte) uint {
	return hash*r.base + uint(c)
}

func (r *RollingHash) expand(sz int) {
	if len(r.power) < sz+1 {
		preSz := len(r.power)
		r.power = append(r.power, make([]uint, sz+1-preSz)...)
		for i := preSz - 1; i < sz; i++ {
			r.power[i+1] = r.power[i] * r.base
		}
	}
}

type V = uint
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}

func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}

func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}

func (d *Dictionary) HasValue(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}

func (d *Dictionary) Size() int {
	return len(d._idToValue)
}

func min2(a, b int) int {
	if a < b {
		return a
	}
	return b
}
