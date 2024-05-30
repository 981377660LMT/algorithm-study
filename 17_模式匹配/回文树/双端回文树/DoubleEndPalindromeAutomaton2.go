package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo()
}

func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int32
	fmt.Fscan(in, &q)
	T := NewDoubleEndedPalindromicTree(q)
	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 0 {
			var c string
			fmt.Fscan(in, &c)
			T.PushFront(int32(c[0]))
		}
		if op == 1 {
			var c string
			fmt.Fscan(in, &c)
			T.PushBack(int32(c[0]))
		}
		if op == 2 {
			T.PopFront()
		}
		if op == 3 {
			T.PopBack()
		}
		a := T.CountDistinctPalindrome()
		b := T.MaximumPrefixPalindrome()
		c := T.MaximumSuffixPalindrome()
		fmt.Fprintln(out, a, b, c)
	}
}

const SIGMA int32 = 26
const OFFSET int32 = 'a'

type node struct {
	to         [SIGMA]int32
	parent     int32
	link       int32 // 最长回文后缀
	length     int32
	count      int32
	linkCount  int32 // 子节点个数
	directLink [SIGMA]int32
}

type bag struct {
	c           int32
	left, right int32
}

const (
	ODD  int32 = 0
	EVEN int32 = 1
)

type DoubleEndedPalindromicTree struct {
	nodes       []*node
	free        []int32
	mod, mask   int32
	data        []*bag
	left, right int32
	numNode     int32
}

func NewDoubleEndedPalindromicTree(maxSize int32) *DoubleEndedPalindromicTree {
	res := &DoubleEndedPalindromicTree{}
	res.newNode(-1, -1, -1, -1)
	res.newNode(-1, 0, 0, -1)
	res.nodes[ODD].count = 1<<31 - 1
	res.nodes[EVEN].count = 1<<31 - 1
	res.mod = 4
	for res.mod < maxSize {
		res.mod *= 2
	}
	res.data = make([]*bag, res.mod)
	for i := int32(0); i < res.mod; i++ {
		res.data[i] = &bag{}
	}
	res.mask = res.mod - 1
	return res
}

func (t *DoubleEndedPalindromicTree) PushBack(c int32) {
	c -= OFFSET
	v := t.suffixNode()
	t.data[t.right&t.mask].c = c
	getV := func() int32 {
		var dfs func(int32) int32
		dfs = func(v int32) int32 {
			w := &t.nodes[v].directLink[c]
			if *w != -1 {
				return *w
			}
			p := t.nodes[v].link
			j := t.right - 1 - t.nodes[p].length
			if t.left <= j && j <= t.right && t.data[j&t.mask].c == c {
				*w = p
				return *w
			}
			*w = dfs(p)
			return *w
		}
		j := t.right - 1 - t.nodes[v].length
		if !(t.left <= j && j <= t.right && t.data[j&t.mask].c == c) {
			v = dfs(v)
		}
		if tmp := t.nodes[v].to[c]; tmp != -1 {
			return tmp
		}
		link := EVEN
		if v != ODD {
			link = t.nodes[dfs(v)].to[c]
		}
		return t.newNode(v, link, t.nodes[v].length+2, c)
	}
	v = getV()
	t.data[t.right&t.mask].right = EVEN
	t.data[(t.right-t.nodes[v].length+1)&t.mask].right = v
	t.data[t.right&t.mask].left = v
	w := t.nodes[v].link
	k := t.right - t.nodes[v].length + t.nodes[w].length
	if t.nodes[w].length >= 1 && t.data[k&t.mask].left == w {
		t.data[k&t.mask].left = EVEN
	}
	t.right++
	t.nodes[v].count++
}

func (t *DoubleEndedPalindromicTree) PopBack() {
	v := t.suffixNode()
	w := t.nodes[v].link
	k := t.right - 1 - t.nodes[v].length + t.nodes[w].length
	if t.nodes[t.data[k&t.mask].left].length < t.nodes[w].length {
		t.data[k&t.mask].left = w
		t.data[(k-t.nodes[w].length+1)&t.mask].right = w
	} else {
		t.data[(k-t.nodes[w].length+1)&t.mask].right = EVEN
	}
	t.nodes[v].count--
	if t.nodes[v].linkCount == 0 && t.nodes[v].count == 0 {
		t.removeNode(v, t.data[(t.right-1)&t.mask].c)
	}
	t.right--
}

func (t *DoubleEndedPalindromicTree) PushFront(c int32) {
	c -= OFFSET
	v := t.prefixNode()
	t.data[(t.left-1)&t.mask].c = c
	getV := func() int32 {
		var dfs func(int32) int32
		dfs = func(v int32) int32 {
			w := &t.nodes[v].directLink[c]
			if *w != -1 {
				return *w
			}
			p := t.nodes[v].link
			j := t.left + t.nodes[p].length
			if t.left-1 <= j && j <= t.right-1 && t.data[j&t.mask].c == c {
				*w = p
				return *w
			}
			*w = dfs(p)
			return *w
		}
		j := t.left + t.nodes[v].length
		if !(t.left-1 <= j && j <= t.right-1 && t.data[j&t.mask].c == c) {
			v = dfs(v)
		}
		if tmp := t.nodes[v].to[c]; tmp != -1 {
			return tmp
		}
		link := EVEN
		if v != ODD {
			link = t.nodes[dfs(v)].to[c]
		}
		return t.newNode(v, link, t.nodes[v].length+2, c)
	}
	v = getV()
	t.data[(t.left-1)&t.mask].left = EVEN
	t.data[(t.left-2+t.nodes[v].length)&t.mask].left = v
	t.data[(t.left-1)&t.mask].right = v
	w := t.nodes[v].link
	k := (t.left - 1 + t.nodes[v].length - t.nodes[w].length)
	if t.nodes[w].length >= 1 && t.data[k&t.mask].right == w {
		t.data[k&t.mask].right = EVEN
	}
	t.left--
	t.nodes[v].count++
}

func (t *DoubleEndedPalindromicTree) PopFront() {
	v := t.prefixNode()
	w := t.nodes[v].link
	k := t.left + t.nodes[v].length - t.nodes[w].length
	if t.nodes[t.data[k&t.mask].right].length < t.nodes[w].length {
		t.data[k&t.mask].right = w
		t.data[(k+t.nodes[w].length-1)&t.mask].left = w
	} else {
		t.data[(k+t.nodes[w].length-1)&t.mask].left = EVEN
	}
	t.nodes[v].count--
	if t.nodes[v].linkCount == 0 && t.nodes[v].count == 0 {
		t.removeNode(v, t.data[t.left&t.mask].c)
	}
	t.left++
}

func (t *DoubleEndedPalindromicTree) CountDistinctPalindrome() int32 {
	return t.numNode - 2
}

func (t *DoubleEndedPalindromicTree) MaximumPrefixPalindrome() int32 {
	return t.nodes[t.prefixNode()].length
}

func (t *DoubleEndedPalindromicTree) MaximumSuffixPalindrome() int32 {
	return t.nodes[t.suffixNode()].length
}

func (t *DoubleEndedPalindromicTree) newNode(par, link, length, c int32) int32 {
	t.numNode++
	n := &node{
		parent: par,
		link:   link,
		length: length,
	}
	for i := int32(0); i < SIGMA; i++ {
		n.to[i] = -1
		n.directLink[i] = -1
	}
	if link != -1 {
		t.nodes[link].linkCount++
	}
	p := int32(0)
	if len(t.free) == 0 {
		p = int32(len(t.nodes))
		t.nodes = append(t.nodes, n)
	} else {
		p = t.free[len(t.free)-1]
		t.free = t.free[:len(t.free)-1]
		t.nodes[p] = n
	}
	if par != -1 {
		t.nodes[par].to[c] = p
	}
	return p
}

func (t *DoubleEndedPalindromicTree) removeNode(nid, c int32) {
	t.numNode--
	pid := t.nodes[nid].parent
	t.nodes[pid].to[c] = -1
	k := t.nodes[nid].link
	t.nodes[k].linkCount--
	if t.nodes[k].linkCount == 0 {
		t.free = append(t.free, nid)
	}
}

func (t *DoubleEndedPalindromicTree) suffixNode() int32 {
	if t.left == t.right {
		return EVEN
	}
	return t.data[(t.right-1)&t.mask].left
}

func (t *DoubleEndedPalindromicTree) prefixNode() int32 {
	if t.left == t.right {
		return EVEN
	}
	return t.data[t.left&t.mask].right
}
