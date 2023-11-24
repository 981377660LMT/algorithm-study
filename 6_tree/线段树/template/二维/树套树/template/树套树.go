// 树套树：线段树套树/树状数组套树

package main

func main() {

}

type InnerTree = struct{}

type SegmentTreeDivideInterval struct {
	n               int
	smallN          bool
	offset          int // 线段树中一共offset+n个节点,offset+i对应原来的第i个节点.
	createInnerTree func() InnerTree
	innerTreeList   []InnerTree
	innerTreeMap    map[int]InnerTree
}

// 线段树套数据结构.
// n: 第一个维度(一般是序列)的长度.
// createInnerTree: 创建第二个维度(一般是线段树)的树.
// smallN: n较小时，会预先创建好所有的内层树; 否则会用map保存内层树，并在需要的时候创建.
func NewSegmentTreeDivideInterval(n int, createInnerTree func() InnerTree, smallN bool) *SegmentTreeDivideInterval {
	offset := 1
	for offset < n {
		offset <<= 1
	}
	var innerTreeList []InnerTree
	if smallN {
		innerTreeList = make([]InnerTree, offset+n)
		for i := range innerTreeList {
			innerTreeList[i] = createInnerTree()
		}
	}
	return &SegmentTreeDivideInterval{
		n:               n,
		smallN:          smallN,
		offset:          offset,
		createInnerTree: createInnerTree,
		innerTreeList:   innerTreeList,
		innerTreeMap:    map[int]InnerTree{},
	}
}

func (tree *SegmentTreeDivideInterval) EnumeratePoint(index int, f func(tree InnerTree)) {
	if index < 0 || index >= tree.n {
		return
	}
	index += tree.offset
	for index > 0 {
		f(tree.getTree(index))
		index >>= 1
	}
}

func (tree *SegmentTreeDivideInterval) EnumerateRange(start, end int, f func(tree InnerTree)) {
	if start < 0 {
		start = 0
	}
	if end > tree.n {
		end = tree.n
	}
	if start >= end {
		return
	}

	leftSegments := []InnerTree{}
	rightSegments := []InnerTree{}
	for start, end = start+tree.offset, end+tree.offset; start < end; start, end = start>>1, end>>1 {
		if start&1 == 1 {
			leftSegments = append(leftSegments, tree.getTree(start))
			start++
		}
		if end&1 == 1 {
			end--
			rightSegments = append(rightSegments, tree.getTree(end))
		}
	}

	for i := 0; i < len(leftSegments); i++ {
		f(leftSegments[i])
	}
	for i := len(rightSegments) - 1; i >= 0; i-- {
		f(rightSegments[i])
	}
}

func (tree *SegmentTreeDivideInterval) getTree(segmentId int) InnerTree {
	if tree.smallN {
		return tree.innerTreeList[segmentId]
	} else {
		if v, ok := tree.innerTreeMap[segmentId]; ok {
			return v
		} else {
			newTree := tree.createInnerTree()
			tree.innerTreeMap[segmentId] = newTree
			return newTree
		}
	}
}

// 树状数组套数据结构.
type FenwickTreeDivideInterval struct {
	n               int
	smallN          bool
	createInnerTree func() InnerTree
	innerTreeList   []InnerTree
	innerTreeMap    map[int]InnerTree
}

func NewFenwickTreeDivideInterval(n int, createInnerTree func() InnerTree, smallN bool) *FenwickTreeDivideInterval {
	var innerTreeList []InnerTree
	if smallN {
		innerTreeList = make([]InnerTree, n)
		for i := range innerTreeList {
			innerTreeList[i] = createInnerTree()
		}
	}
	return &FenwickTreeDivideInterval{
		n:               n,
		smallN:          smallN,
		createInnerTree: createInnerTree,
		innerTreeList:   innerTreeList,
		innerTreeMap:    map[int]InnerTree{},
	}
}

func (tree *FenwickTreeDivideInterval) Update(index int, f func(tree InnerTree)) {
	if index < 0 || index >= tree.n {
		return
	}
	for index++; index <= tree.n; index += index & -index {
		f(tree.getTree(index - 1))
	}
}

func (tree *FenwickTreeDivideInterval) QueryPrefix(end int, f func(tree InnerTree)) {
	if end > tree.n {
		end = tree.n
	}
	for end > 0 {
		f(tree.getTree(end - 1))
		end &= end - 1
	}
}

func (tree *FenwickTreeDivideInterval) getTree(segmentId int) InnerTree {
	if tree.smallN {
		return tree.innerTreeList[segmentId]
	} else {
		if v, ok := tree.innerTreeMap[segmentId]; ok {
			return v
		} else {
			newTree := tree.createInnerTree()
			tree.innerTreeMap[segmentId] = newTree
			return newTree
		}
	}
}
