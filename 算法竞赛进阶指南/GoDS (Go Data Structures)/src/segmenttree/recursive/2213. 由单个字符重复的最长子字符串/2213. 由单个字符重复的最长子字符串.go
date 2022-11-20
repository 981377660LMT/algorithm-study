package main

import "math/bits"

func main() {

}

func longestRepeating(s string, queryCharacters string, queryIndices []int) []int {
	chars := []byte(s)
	q := len(queryCharacters)
	res := make([]int, q)
	tree := NewLazySegmentTree(chars)
	for i := 0; i < q; i++ {
		qv, qi := queryCharacters[i], queryIndices[i]
		tree.Update(qi+1, qi+1, qv)
		res[i] = tree.QueryAll()
	}
	return res
}

type LazySegmentTree struct {
	n                   int
	max, preMax, sufMax []int // data
	// 别的一些信息
	chars []byte
}

func NewLazySegmentTree(leaves []byte) *LazySegmentTree {
	cap := 1 << (bits.Len(uint(len(leaves)-1)) + 1)
	// !初始化data和lazy数组 然后建树
	tree := &LazySegmentTree{
		n:      len(leaves),
		max:    make([]int, cap),
		preMax: make([]int, cap),
		sufMax: make([]int, cap),
		chars:  leaves,
	}
	tree._build(1, 1, tree.n)
	return tree
}

// TODO
func (t *LazySegmentTree) _build(root, left, right int) {
	if left == right {
		// !初始化叶子结点信息 例如data和lazy的monoid
		t.max[root] = 1
		t.preMax[root] = 1
		t.sufMax[root] = 1
		return
	}
	mid := (left + right) >> 1
	t._build(root<<1, left, mid)
	t._build(root<<1|1, mid+1, right)
	t._pushUp(root, left, right)
}

func (t *LazySegmentTree) _pushUp(root, left, right int) {
	// !op操作更新root结点的data信息
	leftPre, rightPre := t.preMax[root<<1], t.preMax[root<<1|1]
	leftSuf, rightSuf := t.sufMax[root<<1], t.sufMax[root<<1|1]
	leftMax, rightMax := t.max[root<<1], t.max[root<<1|1]

	t.preMax[root] = leftPre
	t.sufMax[root] = rightSuf

	mid := (left + right) >> 1
	if t.chars[mid-1] == t.chars[mid] {
		t.max[root] = max(leftMax, rightMax, leftSuf+rightPre)
		if leftPre == mid-left+1 {
			t.preMax[root] += rightPre
		}
		if rightSuf == right-mid {
			t.sufMax[root] += leftSuf
		}

	} else {
		t.max[root] = max(leftMax, rightMax)
		t.max[root] = max(leftMax, rightMax)
	}
}

func (t *LazySegmentTree) _propagate(root, left, right int, lazy byte) {
	// !mapping + composition 来更新子节点data和lazy信息
	t.chars[left-1] = lazy
}

func (t *LazySegmentTree) _query(root, L, R, l, r int) int {
	if L <= l && r <= R {
		return t.max[root]
	}

	mid := (l + r) >> 1
	res := 0 // monoid
	if L <= mid {
		res += t._query(root<<1, L, R, l, mid) // op
	}
	if R > mid {
		res += t._query(root<<1|1, L, R, mid+1, r) // op
	}
	return res
}

func (t *LazySegmentTree) _update(root, L, R, l, r int, val byte) {
	if L <= l && r <= R {
		t._propagate(root, l, r, val)
		return
	}

	mid := (l + r) >> 1
	if L <= mid {
		t._update(root<<1, L, R, l, mid, val)
	}
	if R > mid {
		t._update(root<<1|1, L, R, mid+1, r, val)
	}
	t._pushUp(root, l, r)
}

// public api
func (t *LazySegmentTree) Query(left, right int) int         { return t._query(1, left, right, 1, t.n) }
func (t *LazySegmentTree) Update(left, right int, lazy byte) { t._update(1, left, right, 1, t.n, lazy) }
func (t *LazySegmentTree) QueryAll() int                     { return t.max[1] }

func max(nums ...int) int {
	res := nums[0]
	for _, v := range nums {
		if v > res {
			res = v
		}
	}
	return res
}
