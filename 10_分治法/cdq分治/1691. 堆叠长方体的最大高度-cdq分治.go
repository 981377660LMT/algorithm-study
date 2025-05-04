// 1691. 堆叠长方体的最大高度
// O(nlog^2n)解法 二维树状数组 cdq分治
// https://leetcode.cn/problems/maximum-height-by-stacking-cuboids/solutions/2014599/by-vclip-gmet/
// 如果 widthi <= widthj 且 lengthi <= lengthj 且 heighti <= heightj ，
// 你就可以将长方体 i 堆叠在长方体 j 上。你可以通过旋转把长方体的长宽高重新排列，以将它放在另一个长方体上。
// 请你从 cuboids 选出一个 子集 ，并将它们堆叠起来。
// 返回 堆叠长方体 cuboids 可以得到的 最大高度 。

package main

import (
	"fmt"
	"sort"
)

func main() {
	cuboids1 := [][]int{{50, 45, 20}, {95, 37, 53}, {45, 23, 12}}
	fmt.Println(maxHeight(cuboids1))

	cuboids2 := [][]int{{38, 25, 45}, {76, 35, 3}}
	fmt.Println(maxHeight(cuboids2))

	cuboids3 := [][]int{{7, 11, 17}, {7, 17, 11}, {11, 7, 17}, {11, 17, 7}, {17, 7, 11}, {17, 11, 7}}
	fmt.Println(maxHeight(cuboids3))
}

type Item struct {
	w, l, h0 int
	h, s     int
}

func cdq(items []*Item, left, right int, buf []*Item, tree *BIT) int {
	n := right - left
	if n == 0 {
		return 0
	}
	if n == 1 {
		return items[left].s
	}

	mid := left + n/2

	res := cdq(items, left, mid, buf, tree)

	sort.Slice(items[mid:right], func(i, j int) bool {
		return items[mid+i].l < items[mid+j].l
	})

	it1 := left
	for it2 := mid; it2 < right; it2++ {
		for ; it1 < mid; it1++ {
			if items[it2].l < items[it1].l {
				break
			}
			tree.update(items[it1].h0, items[it1].s)
		}
		newS := items[it2].h + tree.query(items[it2].h0+1)
		if newS > items[it2].s {
			items[it2].s = newS
		}
	}

	for it := left; it < it1; it++ {
		tree.clear(items[it].h0)
	}

	sort.Slice(items[mid:right], func(i, j int) bool {
		if items[mid+i].w != items[mid+j].w {
			return items[mid+i].w < items[mid+j].w
		}
		if items[mid+i].l != items[mid+j].l {
			return items[mid+i].l < items[mid+j].l
		}
		return items[mid+i].h0 < items[mid+j].h0
	})

	rightRes := cdq(items, mid, right, buf, tree)
	if rightRes > res {
		res = rightRes
	}

	merge(items[left:mid], items[mid:right], buf[:n], func(a, b *Item) bool {
		return a.l < b.l
	})

	copy(items[left:right], buf[:n])

	return res
}

func merge(a, b, result []*Item, less func(*Item, *Item) bool) {
	i, j, k := 0, 0, 0
	for i < len(a) && j < len(b) {
		if less(a[i], b[j]) {
			result[k] = a[i]
			i++
		} else {
			result[k] = b[j]
			j++
		}
		k++
	}

	for ; i < len(a); i++ {
		result[k] = a[i]
		k++
	}

	for ; j < len(b); j++ {
		result[k] = b[j]
		k++
	}
}

func maxHeight(cuboids [][]int) int {
	n := len(cuboids)
	items := make([]*Item, n)
	heightSet := make([]int, 0, n)

	for _, e := range cuboids {
		sort.Ints(e)
		heightSet = append(heightSet, e[2])
	}

	sort.Ints(heightSet)
	heightMap := make(map[int]int)
	uniqueHeights := 0
	for _, h := range heightSet {
		if _, exists := heightMap[h]; !exists {
			heightMap[h] = uniqueHeights
			uniqueHeights++
		}
	}

	for i, e := range cuboids {
		h0 := heightMap[e[2]]
		items[i] = &Item{e[0], e[1], h0, e[2], e[2]}
	}

	sort.Slice(items, func(i, j int) bool {
		if items[i].w != items[j].w {
			return items[i].w < items[j].w
		}
		if items[i].l != items[j].l {
			return items[i].l < items[j].l
		}
		return items[i].h0 < items[j].h0
	})

	buf := make([]*Item, n)
	tree := NewBIT(uniqueHeights)

	return cdq(items, 0, n, buf, tree)
}

type BIT struct {
	n    int
	tree []int
}

func NewBIT(n int) *BIT {
	return &BIT{
		n:    n,
		tree: make([]int, n),
	}
}

func (b *BIT) clear(p int) {
	for i := p + 1; i <= b.n; i += i & -i {
		if b.tree[i-1] == 0 {
			break
		}
		b.tree[i-1] = 0
	}
}

func (b *BIT) update(p, l int) {
	for i := p + 1; i <= b.n; i += i & -i {
		if b.tree[i-1] < l {
			b.tree[i-1] = l
		}
	}
}

func (b *BIT) query(p int) int {
	res := 0
	for i := p; i > 0; i -= i & -i {
		if b.tree[i-1] > res {
			res = b.tree[i-1]
		}
	}
	return res
}
