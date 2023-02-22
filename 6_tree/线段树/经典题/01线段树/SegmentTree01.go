// 01线段树 (golang)

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"strings"
)

// https://leetcode.cn/problems/handling-sum-queries-after-update/
func handleQuery(nums1 []int, nums2 []int, queries [][]int) []int64 {
	sum := 0
	for _, num := range nums2 {
		sum += num
	}
	seg01 := NewSegmentTree01(nums1)
	res := make([]int64, 0, len(queries))
	for _, query := range queries {
		op, a, b := query[0], query[1], query[2]
		if op == 1 {
			seg01.Flip(a+1, b+1)
		} else if op == 2 {
			ones := seg01.OnesCount(a+1, b+1)
			sum += ones * a
		} else {
			res = append(res, int64(sum))
		}
	}
	return res
}

type SegmentTree01 struct {
	n        int
	ones     []int
	lazyFlip []bool
}

// 01线段树，支持 flip/indexOf/lastIndexOf/onesCount/kth，可用于模拟Bitset
func NewSegmentTree01(nums []int) *SegmentTree01 {
	if len(nums) == 0 {
		panic("len(bits) == 0")
	}
	n := len(nums)
	log := int(bits.Len(uint(n - 1)))
	size := 1 << log
	ones := make([]int, 2*size)
	lazyFlip := make([]bool, size)
	res := &SegmentTree01{n: n, ones: ones, lazyFlip: lazyFlip}
	res.build(1, 1, n, nums)
	return res
}

// 1 <= left <= right <= n
func (tree *SegmentTree01) Flip(left, right int) {
	tree.flip(1, left, right, 1, tree.n)
}

// 1 <= left <= right <= n
func (tree *SegmentTree01) OnesCount(left, right int) int {
	return tree.onesCount(1, left, right, 1, tree.n)
}

// 1 <= left <= right <= n
//  digit: 0 or 1
//  position: 搜索的起点, 1 <= position <= n
func (tree *SegmentTree01) IndexOf(digit, position int) int {
	if position > tree.n {
		return -1
	}
	if digit == 0 {
		return tree.indexOfZero(1, position, 1, tree.n)
	}
	return tree.indexOfOne(1, position, 1, tree.n)
}

// 1 <= left <= right <= n
//  digit: 0 or 1
//  position: 搜索的起点, 1 <= position <= n
func (tree *SegmentTree01) LastIndexOf(digit, position int) int {
	if position < 1 {
		return -1
	}
	if digit == 0 {
		return tree.lastIndexOfZero(1, position, 1, tree.n)
	}
	return tree.lastIndexOfOne(1, position, 1, tree.n)
}

// 树上二分查询第k个0/1的位置.如果不存在第k个0/1，返回-1.
//  1 <= k <= n
func (tree *SegmentTree01) Kth(digit, k int) int {
	if digit == 0 {
		if k > tree.n-tree.ones[1] {
			return -1
		}
		return tree.kthZero(1, k, 1, tree.n)
	}
	if k > tree.ones[1] {
		return -1
	}
	return tree.kthOne(1, k, 1, tree.n)
}

func (tree *SegmentTree01) String() string {
	var sb []string
	tree.toString(1, 1, tree.n, &sb)
	return strings.Join(sb, "")
}

func (tree *SegmentTree01) flip(root, left, right, l, r int) {
	if left <= l && r <= right {
		tree.propagateFlip(root, l, r)
		return
	}
	mid := (l + r) >> 1
	tree.pushDown(root, l, r)
	if left <= mid {
		tree.flip(root<<1, left, right, l, mid)
	}
	if right > mid {
		tree.flip(root<<1|1, left, right, mid+1, r)
	}
	tree.pushUp(root)
}

func (tree *SegmentTree01) onesCount(root, left, right, l, r int) int {
	if left <= l && r <= right {
		return tree.ones[root]
	}
	mid := (l + r) >> 1
	tree.pushDown(root, l, r)
	var res int
	if left <= mid {
		res += tree.onesCount(root<<1, left, right, l, mid)
	}
	if right > mid {
		res += tree.onesCount(root<<1|1, left, right, mid+1, r)
	}
	return res
}

func (tree *SegmentTree01) indexOfZero(root, position, l, r int) int {
	if l == r {
		if tree.ones[root] == 0 {
			return l
		}
		return -1
	}
	mid := (l + r) >> 1
	tree.pushDown(root, l, r)
	if position <= mid && tree.ones[root<<1] < mid-l+1 {
		leftPos := tree.indexOfZero(root<<1, position, l, mid)
		if leftPos != -1 {
			return leftPos
		}
	}
	return tree.indexOfZero(root<<1|1, position, mid+1, r)
}

func (tree *SegmentTree01) indexOfOne(root, position, l, r int) int {
	if l == r {
		if tree.ones[root] > 0 {
			return l
		}
		return -1
	}
	mid := (l + r) >> 1
	tree.pushDown(root, l, r)
	if position <= mid && tree.ones[root<<1] > 0 {
		leftPos := tree.indexOfOne(root<<1, position, l, mid)
		if leftPos != -1 {
			return leftPos
		}
	}
	return tree.indexOfOne(root<<1|1, position, mid+1, r)
}

func (tree *SegmentTree01) lastIndexOfZero(root, position, l, r int) int {
	if l == r {
		if tree.ones[root] == 0 {
			return l
		}
		return -1
	}
	mid := (l + r) >> 1
	tree.pushDown(root, l, r)
	if position > mid && tree.ones[root<<1|1] < r-mid {
		rightPos := tree.lastIndexOfZero(root<<1|1, position, mid+1, r)
		if rightPos != -1 {
			return rightPos
		}
	}
	return tree.lastIndexOfZero(root<<1, position, l, mid)
}

func (tree *SegmentTree01) lastIndexOfOne(root, position, l, r int) int {
	if l == r {
		if tree.ones[root] > 0 {
			return l
		}
		return -1
	}
	mid := (l + r) >> 1
	tree.pushDown(root, l, r)
	if position > mid && tree.ones[root<<1|1] > 0 {
		rightPos := tree.lastIndexOfOne(root<<1|1, position, mid+1, r)
		if rightPos != -1 {
			return rightPos
		}
	}
	return tree.lastIndexOfOne(root<<1, position, l, mid)
}

func (tree *SegmentTree01) kthZero(root, k, l, r int) int {
	if l == r {
		return l
	}
	mid := (l + r) >> 1
	tree.pushDown(root, l, r)
	leftZeros := mid - l + 1 - tree.ones[root<<1]
	if k <= leftZeros {
		return tree.kthZero(root<<1, k, l, mid)
	}
	return tree.kthZero(root<<1|1, k-leftZeros, mid+1, r)
}

func (tree *SegmentTree01) kthOne(root, k, l, r int) int {
	if l == r {
		return l
	}
	mid := (l + r) >> 1
	tree.pushDown(root, l, r)
	if k <= tree.ones[root<<1] {
		return tree.kthOne(root<<1, k, l, mid)
	}
	return tree.kthOne(root<<1|1, k-tree.ones[root<<1], mid+1, r)
}

func (tree *SegmentTree01) toString(root, l, r int, sb *[]string) {
	if l == r {
		if tree.ones[root] == 1 {
			*sb = append(*sb, "1")
		} else {
			*sb = append(*sb, "0")
		}
		return
	}
	mid := (l + r) >> 1
	tree.pushDown(root, l, r)
	tree.toString(root<<1, l, mid, sb)
	tree.toString(root<<1|1, mid+1, r, sb)
}

// build
func (tree *SegmentTree01) build(root, l, r int, leaves []int) {
	if l == r {
		tree.ones[root] = leaves[l-1]
		return
	}
	mid := (l + r) >> 1
	tree.build(root<<1, l, mid, leaves)
	tree.build(root<<1|1, mid+1, r, leaves)
	tree.pushUp(root)
}

func (tree *SegmentTree01) propagateFlip(root, l, r int) {
	tree.ones[root] = r - l + 1 - tree.ones[root]
	if root < len(tree.lazyFlip) {
		tree.lazyFlip[root] = !tree.lazyFlip[root]
	}
}

func (tree *SegmentTree01) pushDown(root, l, r int) {
	if tree.lazyFlip[root] {
		tree.propagateFlip(root<<1, l, (l+r)>>1)
		tree.propagateFlip(root<<1|1, ((l+r)>>1)+1, r)
		tree.lazyFlip[root] = false
	}
}

func (tree *SegmentTree01) pushUp(root int) {
	tree.ones[root] = tree.ones[root<<1] + tree.ones[root<<1|1]
}

func main() {
	seg01 := NewSegmentTree01([]int{1, 0, 1, 0, 1, 0, 1, 0, 1, 0})

	// test flip/onesCount/lastIndexOf/kth/indexOf, generated by copilot
	for i := 0; i < 100; i++ {
		n := rand.Intn(100) + 1
		seg01 = NewSegmentTree01(make([]int, n))
		nums01 := make([]int, n)
		for i := 0; i < 1000; i++ {
			// flip
			left := rand.Intn(n) + 1
			right := rand.Intn(n) + 1
			if left > right {
				left, right = right, left
			}
			seg01.Flip(left, right)
			for j := left - 1; j < right; j++ {
				nums01[j] ^= 1
			}
			for j := 0; j < n; j++ {
				if nums01[j] != seg01.OnesCount(j+1, j+1) {
					panic("checkSame failed at flip")
				}
			}

			// onesCount
			left = rand.Intn(n) + 1
			right = rand.Intn(n) + 1
			if left > right {
				left, right = right, left
			}
			arrOnesCount := 0
			for j := left - 1; j < right; j++ {
				arrOnesCount += nums01[j]
			}
			if arrOnesCount != seg01.OnesCount(left, right) {
				panic("checkSame failed at onesCount")
			}

			// lastIndexOf
			digit := rand.Intn(2)
			position := rand.Intn(n) + 1
			arrLastIndexOf := -1
			for j := position - 1; j >= 0; j-- {
				if nums01[j] == digit {
					arrLastIndexOf = j
					break
				}
			}
			segLast := seg01.LastIndexOf(digit, position)
			if arrLastIndexOf == -1 {
				if segLast != -1 {
					panic("checkSame failed at lastIndexOf")
				}
			} else if arrLastIndexOf+1 != segLast {
				panic("checkSame failed at lastIndexOf")
			}

			// kth
			digit = rand.Intn(2)
			k := rand.Intn(n) + 1
			arrKth := -1
			count := 0
			for j := 0; j < n; j++ {
				if nums01[j] == digit {
					count++
				}
				if count == k {
					arrKth = j
					break
				}
			}
			segKth := seg01.Kth(digit, k)
			if arrKth == -1 {
				if segKth != -1 {
					panic("checkSame failed at kth")
				}
			} else if arrKth+1 != segKth {
				panic(fmt.Sprintf("checkSame failed at kth,  arrKth: %d, segKth: %d, digit: %d, k: %d", arrKth, segKth, digit, k))
			}

			// indexOf
			digit = rand.Intn(2)
			position = rand.Intn(n) + 1
			arrIndexOf := -1
			for j := position - 1; j < n; j++ {
				if nums01[j] == digit {
					arrIndexOf = j
					break
				}
			}
			segIndex := seg01.IndexOf(digit, position)
			if arrIndexOf == -1 {
				if segIndex != -1 {
					panic("checkSame failed at indexOf")
				}
			} else if arrIndexOf+1 != segIndex {
				panic("checkSame failed at indexOf")
			}
		}
	}

	fmt.Println("all tests passed!")
}
