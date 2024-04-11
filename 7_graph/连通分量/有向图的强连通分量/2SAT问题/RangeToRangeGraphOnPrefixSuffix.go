package main

func main() {

}

// 4 ← 5 ← 6 ← 7  (前缀)
// ↓	 ↓   ↓   ↓
// 0   1   2   3
// ↑	 ↑   ↑   ↑
// 8 → 9 → 10 → 11 (后缀)

// 前后缀优化建图.
type RangeToRangeGraphOnPrefixSuffix struct {
	n        int32
	maxSize  int32
	allocPtr int32
}

// 新建一个区间图，n 为原图的节点数，rangeToRangeOpCount 为区间到区间的最大操作次数.
func NewRangeToRangeGraphOnPrefixSuffix(n int32, rangeToRangeOpCount int32) *RangeToRangeGraphOnPrefixSuffix {
	return &RangeToRangeGraphOnPrefixSuffix{n: n, maxSize: 3*n + rangeToRangeOpCount, allocPtr: 3 * n}
}

// 新图的结点数.前n个节点为原图的节点.
func (g *RangeToRangeGraphOnPrefixSuffix) Size() int32 { return g.maxSize }

func (g *RangeToRangeGraphOnPrefixSuffix) Init(f func(from, to int32)) {
}

// 添加有向边 from -> to.
func (g *RangeToRangeGraphOnPrefixSuffix) Add(from, to int32, f func(from, to int32)) {
	f(from, to)
}

func (g *RangeToRangeGraphOnPrefixSuffix) AddFromPrefix(prefixEnd int32, to int32, f func(from, to int32)) {
}

func (g *RangeToRangeGraphOnPrefixSuffix) AddFromSuffix(suffixStart int32, to int32, f func(from, to int32)) {
}

func (g *RangeToRangeGraphOnPrefixSuffix) AddToPrefix(from int32, prefixEnd int32, f func(from, to int32)) {
}

func (g *RangeToRangeGraphOnPrefixSuffix) AddToSuffix(from int32, suffixStart int32, f func(from, to int32)) {
}

func (g *RangeToRangeGraphOnPrefixSuffix) AddPrefixToSuffix(prefixEnd, suffixStart int32, f func(from, to int32)) {
	newNode := g.allocPtr
	g.allocPtr++
	g.AddFromPrefix(prefixEnd, newNode, f)
	g.AddToSuffix(newNode, suffixStart, f)
}

func (g *RangeToRangeGraphOnPrefixSuffix) AddSuffixToPrefix(suffixStart, prefixEnd int32, f func(from, to int32)) {
	newNode := g.allocPtr
	g.allocPtr++
	g.AddFromSuffix(suffixStart, newNode, f)
	g.AddToPrefix(newNode, prefixEnd, f)
}
