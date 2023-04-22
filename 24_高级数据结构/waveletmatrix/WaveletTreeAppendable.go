// WaveletTreeAppendable
// NewAppendableWaveletTree(maxVal int) *AppendableWaveletTree
// Append(v int)
// Rank(v int, left, right int) int
// RangeFreq(left, right int, lower, upper int) int

package main

import "fmt"

func main() {
	dic := NewAppendableRankDictionary()
	dic.Append(true)
	dic.Append(false)
	fmt.Println(dic.Rank(3))

	wt := NewAppendableWaveletTree(10)
	wt.Append(1)
	wt.Append(2)
	wt.Append(3)
	wt.Append(4)
	wt.Append(5)
	wt.Append(6)
	fmt.Println(wt.Rank(2, 0, 2))
	fmt.Println(wt.RangeFreq(0, 2, 3, 4))
}

type AppendableWaveletTree struct {
	dicts           [][]*AppendableRankDictionary
	maxVal          int
	bitSize, length int
}

// 支持append的WaveletTree.
func NewAppendableWaveletTree(maxVal int) *AppendableWaveletTree {
	res := &AppendableWaveletTree{}
	bitSize := 0
	for 1<<bitSize <= maxVal {
		bitSize++
	}
	dicts := make([][]*AppendableRankDictionary, bitSize)
	for bit := 0; bit < bitSize; bit++ {
		dicts[bit] = make([]*AppendableRankDictionary, 1<<bit)
		for i := 0; i < 1<<bit; i++ {
			dicts[bit][i] = NewAppendableRankDictionary()
		}
	}
	res.dicts = dicts
	res.maxVal = maxVal
	res.bitSize = bitSize
	return res
}

// 0<=v<=maxVal
func (t *AppendableWaveletTree) Append(v int) {
	node := 0
	for bit := 0; bit < t.bitSize; bit++ {
		dir := v >> (t.bitSize - bit - 1) & 1
		t.dicts[bit][node].Append(dir == 1)
		node = node*2 + dir
	}
	t.length++
}

// [left,right)内v是第几小(0-based)/比v小的数有几个.
func (t *AppendableWaveletTree) Rank(v int, left, right int) int {
	var outLt, outGt int
	t._rankAll(v, left, right, &outLt, &outGt)
	return outLt
}

// [left,right)内在[lower,upper)内的数有几个.
func (t *AppendableWaveletTree) RangeFreq(left, right int, lower, upper int) int {
	return t.Rank(upper, left, right) - t.Rank(lower, left, right)
}

func (t *AppendableWaveletTree) _rankAll(v int, left, right int, outLt, outGt *int) {
	if v > t.maxVal {
		*outLt = right - left
		*outGt = 0
		return
	}
	*outLt = 0
	*outGt = 0
	node := 0
	for bit := 0; ; bit++ {
		dir := v >> (t.bitSize - bit - 1) & 1
		rightCount := t.dicts[bit][node].RankBool(true, right)
		leftCount := t.dicts[bit][node].RankBool(true, left)
		count := rightCount - leftCount
		if dir == 1 {
			*outLt += (right - left) - count
		} else {
			*outGt += count
		}
		if bit == t.bitSize-1 {
			return
		}
		if dir == 1 {
			left = leftCount
			right = rightCount
		} else {
			left -= leftCount
			right -= rightCount
		}
		node = node*2 + dir
	}
}

type AppendableRankDictionary struct {
	length, blockslength, count int
	blocks                      []uint32
	ranktable                   []int
}

func NewAppendableRankDictionary() *AppendableRankDictionary {
	return &AppendableRankDictionary{
		blockslength: 1,
		blocks:       make([]uint32, 1),
		ranktable:    make([]int, 2),
	}
}

func (d *AppendableRankDictionary) Append(b bool) {
	if b {
		block := d.length >> 5
		d.blocks[block] |= 1 << d.length & 31
		d.ranktable[block+1]++
		d.count++
	}
	if d.length++; d.length&31 == 0 {
		d.blockslength++
		d.blocks = append(d.blocks, 0)
		d.ranktable = append(d.ranktable, d.ranktable[len(d.ranktable)-1])
	}
}

// [0..pos)中1的个数.
func (d *AppendableRankDictionary) Rank(pos int) int {
	blockId := pos >> 5
	return d.ranktable[blockId] + popcount32(d.blocks[blockId]&(1<<pos&31-1))
}

func (d *AppendableRankDictionary) RankBool(b bool, pos int) int {
	if b {
		return d.Rank(pos)
	}
	return pos - d.Rank(pos)
}

func (d *AppendableRankDictionary) RankBoolRange(b bool, left, right int) int {
	return d.RankBool(b, right) - d.RankBool(b, left)
}

func popcount32(x uint32) int {
	x = (x - (x>>1)&0x55555555)
	x = (x & 0x33333333) + ((x >> 2) & 0x33333333)
	return int(((x + (x>>4)&0xF0F0F0F) * 0x1010101) >> 24)
}
