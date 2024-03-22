// https://raw.githubusercontent.com/EndlessCheng/codeforces-go/master/copypasta/bst.go

package copypasta

import (
	. "fmt"
	"strings"
)

// 二叉树常用函数
// M 表示多重集.

type bstNode struct {
	lr    [2]*bstNode
	msz   int // 子树和(可重集)
	key   int // 用于比较、定位的key
	count int // 次数(可重集)
}

// 设置如下返回值是为了方便使用 bstNode 中的 lr 数组
func (o *bstNode) cmp(a int) int {
	b := o.key
	if a == b {
		return -1
	}
	if a < b {
		return 0 // 左儿子
	}
	return 1 // 右儿子
}

func (o *bstNode) mSize() int {
	if o != nil {
		return o.msz
	}
	return 0
}

func (o *bstNode) maintain() {
	o.msz = int(o.count) + o.lr[0].mSize() + o.lr[1].mSize()
}

func (o *bstNode) pushDown() {
	// ...
}

type bst struct {
	root *bstNode
}

func newBST() *bst {
	return &bst{}
}

func buildBST(a []int) *bstNode {
	if len(a) == 0 {
		return nil
	}
	m := len(a) / 2
	o := &bstNode{key: a[m]} // 也可以事先预分配一个 []bstNode，指针指向数组元素
	o.lr[0] = buildBST(a[:m])
	o.lr[1] = buildBST(a[m+1:])
	o.maintain()
	return o
}

// a 需要是有序的，这样我们可以把它当成中序遍历来构造 BST
func newBSTWithArray(a []int) *bst {
	return &bst{buildBST(a)}
}

func (t *bst) empty() bool { return t.root == nil }

func (t *bst) find(key int) *bstNode {
	for o := t.root; o != nil; {
		if c := o.cmp(key); c >= 0 {
			o = o.lr[c]
		} else {
			return o
		}
	}
	return nil
}

// max <= key
// return nil if not found
func (t *bst) floor(key int) (floor *bstNode) {
	for o := t.root; o != nil; {
		switch c := o.cmp(key); {
		case c == 0:
			o = o.lr[0]
		case c > 0:
			floor = o
			o = o.lr[1]
		default:
			return o
		}
	}
	return
}

// min >= key
// return nil if not found
func (t *bst) lowerBound(key int) (lb *bstNode) {
	for o := t.root; o != nil; {
		switch c := o.cmp(key); {
		case c == 0:
			lb = o
			o = o.lr[0]
		case c > 0:
			o = o.lr[1]
		default:
			return o
		}
	}
	return
}

// 前驱（小于 key，且最大的数）
// 等价于 floor(key-1)
func (t *bst) prev(key int) (prev *bstNode) {
	// 另一种写法，适用于含有 lazy delete 的 BST，如替罪羊树等
	// rk, _ := t.mRank(key)
	// return t.mSelect(rk - 1)
	for o := t.root; o != nil; {
		if o.cmp(key) <= 0 {
			o = o.lr[0]
		} else {
			prev = o
			o = o.lr[1]
		}
	}
	return
}

// 后继（大于 key，且最小的数)
// 等价于 lowerBound(key+1)
func (t *bst) next(key int) (next *bstNode) {
	// 另一种写法，适用于含有 lazy delete 的 BST，如替罪羊树等
	// rk, o := t.mRank(key)
	// if o != nil {
	// 	 rk += o.value
	// }
	// return t.mSelect(rk)
	for o := t.root; o != nil; {
		if o.cmp(key) != 0 {
			o = o.lr[1]
		} else {
			next = o
			o = o.lr[0]
		}
	}
	return
}

func (*bst) put(int, int) {}
func (*bst) delete(int)   {}

// < key 的元素个数
func (t *bst) mRank(key int) (cnt int) {
	for o := t.root; o != nil; {
		switch c := o.cmp(key); {
		case c == 0:
			o = o.lr[0]
		case c > 0:
			cnt += o.count + o.lr[0].mSize()
			o = o.lr[1]
		default:
			cnt += o.lr[0].mSize()
			// 额外加上 1 或 o.dupCnt 就是 <= key 的元素个数
			return
		}
	}
	return
}

// >= key 的元素个数
// 等价于 t.root.size() - t.mRank(key)
func (t *bst) lowerCount(key int) (cnt int) {
	for o := t.root; o != nil; {
		switch c := o.cmp(key); {
		case c == 0:
			cnt += o.count + o.lr[1].mSize()
			o = o.lr[0]
		case c > 0:
			o = o.lr[1]
		default:
			cnt += o.count + o.lr[1].mSize()
			return
		}
	}
	return
}

// kth: 排名为 k 的节点 o（即有 k 个键小于 o.key）
// 维护子树和的写法见 https://codeforces.com/contest/1398/submission/119651187
func (t *bst) mSelect(k int) (o *bstNode) {
	//if k < 0 {
	//	return
	//}
	for o = t.root; o != nil; {
		if ls := o.lr[0].mSize(); k < ls {
			o = o.lr[0]
		} else {
			k -= ls + o.count
			if k < 0 {
				return
			}
			o = o.lr[1]
		}
	}
	return
}

func (t *bst) min() (min *bstNode) {
	for o := t.root; o != nil; o = o.lr[0] {
		min = o
	}
	return
}

func (t *bst) max() (max *bstNode) {
	for o := t.root; o != nil; o = o.lr[1] {
		max = o
	}
	return
}

// 中序遍历，返回所有键值
func (t *bst) keys() []int {
	keys := make([]int, 0) // TODO:mSize()
	var f func(*bstNode)
	f = func(o *bstNode) {
		if o == nil {
			return
		}
		o.pushDown()
		f(o.lr[0])
		keys = append(keys, o.key)
		// 如果是多重集则需要多次插入
		//for i := 0; i < o.value; i++ {
		//	keys = append(keys, o.key)
		//}
		f(o.lr[1])
	}
	f(t.root)
	return keys
}

// 中序遍历（如果是多重集请用下面的 foreachM）
func (t *bst) foreach(do func(o *bstNode) (Break bool)) {
	var f func(*bstNode) bool
	f = func(o *bstNode) bool {
		if o == nil {
			return false
		}
		o.pushDown()
		return f(o.lr[0]) || do(o) || f(o.lr[1])
	}
	f(t.root)
}

// 中序遍历，适用于多重集
func (t *bst) foreachM(do func(o *bstNode) (Break bool)) {
	var f func(*bstNode) bool
	f = func(o *bstNode) bool {
		if o == nil {
			return false
		}
		o.pushDown()
		if f(o.lr[0]) {
			return true
		}
		for i := 0; i < o.count; i++ {
			if do(o) {
				return true
			}
		}
		return f(o.lr[1])
	}
	f(t.root)
}

//

func (o *bstNode) String() (s string) {
	//return strconv.Itoa(o.key)
	if o.count == 1 {
		s = Sprintf("%v", o.key)
	} else {
		s = Sprintf("%v(%v)", o.key, o.count)
	}
	s += Sprintf("[msz:%d]", o.msz)
	return
}

/*
	逆时针旋转 90° 打印这棵树：根节点在最左侧，右子树在上侧，左子树在下侧

效果如下（只打印 key）

Root
│           ┌── 95
│       ┌── 94
│   ┌── 90
│   │   │           ┌── 89
│   │   │       ┌── 88
│   │   │       │   └── 87
│   │   │       │       └── 81
│   │   │   ┌── 74
│   │   └── 66
└── 62

	│           ┌── 59
	│       ┌── 58
	│       │   └── 56
	│       │       └── 47
	│   ┌── 45
	└── 40
	    │       ┌── 37
	    │   ┌── 28
	    └── 25
	        │           ┌── 18
	        │       ┌── 15
	        │   ┌── 11
	        └── 6
	            └── 0
*/
func (o *bstNode) draw(treeSB, prefixSB *strings.Builder, isTail bool) {
	prefix := prefixSB.String()
	if o.lr[1] != nil {
		newPrefixSB := &strings.Builder{}
		newPrefixSB.WriteString(prefix)
		if isTail {
			newPrefixSB.WriteString("│   ")
		} else {
			newPrefixSB.WriteString("    ")
		}
		o.lr[1].draw(treeSB, newPrefixSB, false)
	}
	treeSB.WriteString(prefix)
	if isTail {
		treeSB.WriteString("└── ")
	} else {
		treeSB.WriteString("┌── ")
	}
	treeSB.WriteString(o.String())
	treeSB.WriteByte('\n')
	if o.lr[0] != nil {
		newPrefixSB := &strings.Builder{}
		newPrefixSB.WriteString(prefix)
		if isTail {
			newPrefixSB.WriteString("    ")
		} else {
			newPrefixSB.WriteString("│   ")
		}
		o.lr[0].draw(treeSB, newPrefixSB, true)
	}
}

func (t *bst) String() string {
	if t.root == nil {
		return "Empty\n"
	}
	treeSB := &strings.Builder{}
	treeSB.WriteString("Root\n")
	t.root.draw(treeSB, &strings.Builder{}, true)
	return treeSB.String()
}
