// 重链剖分专用树状数组

package main

// !Range Add Range Sum, 0-based.
type BITArray struct {
	n     int
	tree1 []int
	tree2 []int
}

func NewBITArray(n int) *BITArray {
	return &BITArray{
		n:     n,
		tree1: make([]int, n+1),
		tree2: make([]int, n+1),
	}
}

// 切片内[start, end)的每个元素加上delta.
//  0<=start<=end<=n
func (b *BITArray) Add(start, end, delta int) {
	end--
	b.add(start, delta)
	b.add(end+1, -delta)
}

// 求切片内[start, end)的和.
//  0<=start<=end<=n
func (b *BITArray) Query(start, end int) int {
	end--
	return b.query(end) - b.query(start-1)
}

func (b *BITArray) add(index, delta int) {
	index++
	rawIndex := index
	for index <= b.n {
		b.tree1[index] += delta
		b.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (b *BITArray) query(index int) (res int) {
	index++
	if index > b.n {
		index = b.n
	}
	rawIndex := index
	for index > 0 {
		res += rawIndex*b.tree1[index] - b.tree2[index]
		index -= index & -index
	}
	return
}
