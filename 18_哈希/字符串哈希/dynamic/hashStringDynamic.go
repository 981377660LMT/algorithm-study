// dynamicHashString/hashStringDynamic
// 动态哈希.
//
// api:
//  1. Insert(i, c): 在第i个位置插入字符c.
//  2. Pop(i): 删除第i个位置的字符.
//  3. Set(i, c): 将第i个位置的字符设置为c.
//  4. Get(start, end): 获取[start, end)的哈希值.
//  5. GetAll(): 获取所有字符的哈希值.

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"runtime/debug"
	"time"
)

func main() {
	// demo()
	test()
	testTime()
	abc331_f()
}

func init() {
	debug.SetGCPercent(-1)
}

func demo() {
	s := "asezfvgbadpihoamgkcmco"
	n := int32(len(s))
	base := NewDynamicHashStringBase(n, 37)
	hs := NewDynamicHashString(n, func(i int32) uint64 { return uint64(s[i]) }, base)
	fmt.Println(hs.Get(0, 1))
	fmt.Println(hs.Get(1, 2))
	fmt.Println(hs.Get(2, 3))
	hs.Set(0, 1)
	fmt.Println(hs.Get(0, 1))
	fmt.Println(hs.Get(0, 1))
	fmt.Println(hs.Get(1, 2))
	fmt.Println(hs.Get(2, 3))
}

// F - Palindrome Query
// https://atcoder.jp/contests/abc331/tasks/abc331_f
func abc331_f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	var s string
	fmt.Fscan(in, &s)

	m := int32(len(s))
	base := NewDynamicHashStringBase(m, 0)
	hasher1 := NewDynamicHashString(m, func(i int32) uint64 { return uint64(s[i]) }, base)
	hasher2 := NewDynamicHashString(m, func(i int32) uint64 { return uint64(s[n-i-1]) }, base)

	isPalindrome := func(start, end int32) bool {
		return hasher1.Get(start, end) == hasher2.Get(n-end, n-start)
	}

	for i := int32(0); i < q; i++ {
		var op int32
		fmt.Fscan(in, &op)
		if op == 1 {
			var pos int32
			var c string
			fmt.Fscan(in, &pos, &c)
			pos--
			hasher1.Set(pos, uint64(c[0]))
			hasher2.Set(n-pos-1, uint64(c[0]))
		} else {
			var l, r int32
			fmt.Fscan(in, &l, &r)
			l--
			if isPalindrome(l, r) {
				fmt.Fprintln(out, "Yes")
			} else {
				fmt.Fprintln(out, "No")
			}
		}
	}
}

// https://leetcode.cn/problems/sum-of-scores-of-built-strings/description/
func sumScores(s string) int64 {
	n := int32(len(s))
	base := NewDynamicHashStringBase(n, 0)
	hasher := NewDynamicHashString(n, func(i int32) uint64 { return uint64(s[i]) }, base)
	countPre := func(curLen, start int32) int32 {
		left, right := int32(1), curLen
		for left <= right {
			mid := (left + right) >> 1
			hash1 := hasher.Get(start, start+mid)
			hash2 := hasher.Get(0, mid)
			if hash1 == hash2 {
				left = mid + 1
			} else {
				right = mid - 1
			}
		}

		return right
	}

	res := 0
	for i := int32(1); i < n+1; i++ {
		if s[0] != s[n-i] {
			continue
		}
		count := countPre(i, n-i)
		res += int(count)
	}

	return int64(res)

}

// ModFast/fastMod/mod61
const (
	hashStringMod    uint64 = (1 << 61) - 1
	hashStringModi64 int64  = (1 << 61) - 1
	hashStringMask30 uint64 = (1 << 30) - 1
	hashStringMask31 uint64 = (1 << 31) - 1
	hashStringMASK61 uint64 = hashStringMod
)

type DynamicHashStringBase struct {
	n    int32
	powb []uint64
}

// base: 0 表示随机生成
func NewDynamicHashStringBase(n int32, base uint64) *DynamicHashStringBase {
	res := &DynamicHashStringBase{}
	if base == 0 {
		base = uint64(37 + rand.Intn(1e9))
	}
	powb := make([]uint64, n+1)
	powb[0] = 1
	for i := int32(1); i <= n; i++ {
		powb[i] = res.Mul(powb[i-1], base)
	}
	res.n = n
	res.powb = powb
	return res
}

// h1 <- h2. len(h2) == k.
func (hsb *DynamicHashStringBase) Concat(h1, h2, h2Len uint64) uint64 {
	return hsb.Mod(hsb.Mul(h1, hsb.powb[h2Len]) + h2)
}

// a*b % (2^61-1)
func (hsb *DynamicHashStringBase) Mul(a, b uint64) uint64 {
	au := a >> 31
	ad := a & hashStringMask31
	bu := b >> 31
	bd := b & hashStringMask31
	mid := ad*bu + au*bd
	midu := mid >> 30
	midd := mid & hashStringMask30
	return hsb.Mod(au*bu<<1 + midu + (midd << 31) + ad*bd)
}

// x % (2^61-1)
func (hsb *DynamicHashStringBase) Mod(x uint64) uint64 {
	xu := x >> 61
	xd := x & hashStringMASK61
	res := xu + xd
	if res >= hashStringMod {
		res -= hashStringMod
	}
	return res
}

func (hsb *DynamicHashStringBase) Inv(x uint64) uint64 {
	a, b, u, v, t := int64(x), hashStringModi64, int64(1), int64(0), int64(0)
	for b > 0 {
		t = a / b
		a -= t * b
		a, b = b, a
		u -= t * v
		u, v = v, u
	}
	u %= hashStringModi64
	if u < 0 {
		u += hashStringModi64
	}
	return uint64(u)
}

func (hsb *DynamicHashStringBase) Add(a, b uint64) uint64 {
	res := a + b
	if res >= hashStringMod {
		res -= hashStringMod
	}
	return res
}
func (hsb *DynamicHashStringBase) Sub(a, b uint64) uint64 {
	tmp := a - b
	if tmp >= hashStringMod {
		return hsb.Add(tmp, hashStringMod)
	}
	return tmp
}

type DynamicHashString struct {
	n    int32
	base *DynamicHashStringBase
	avl  *avlTreeWithSum
}

func NewDynamicHashString(n int32, f func(i int32) uint64, base *DynamicHashStringBase) *DynamicHashString {
	res := &DynamicHashString{n: n, base: base}
	res.avl = newAVLTreeWithSum(
		n, func(i int32) E { return E{hash: f(i), len: 1} },
		func() E { return E{} },
		func(a, b E) E { return E{hash: base.Concat(a.hash, b.hash, uint64(b.len)), len: a.len + b.len} },
	)
	return res
}

func (hs *DynamicHashString) Insert(index int32, c uint64) {
	hs.avl.Insert(index, E{hash: c, len: 1})
	hs.n++
}

func (hs *DynamicHashString) Pop(index int32) uint64 {
	hs.n--
	return hs.avl.Pop(index).hash
}

func (hs *DynamicHashString) Get(start, end int32) uint64 {
	if start < 0 {
		start = 0
	}
	if end > hs.n {
		end = hs.n
	}
	if start >= end {
		return 0
	}
	if start == 0 && end == hs.n {
		return hs.GetAll()
	}
	return hs.avl.Query(start, end).hash
}

func (hs *DynamicHashString) GetAll() uint64 {
	return hs.avl.QueryAll().hash
}

func (hs *DynamicHashString) Set(index int32, c uint64) {
	if index < 0 || index >= hs.n {
		return
	}
	hs.avl.Set(index, E{hash: c, len: 1})
}

func (hs *DynamicHashString) Len() int32 { return hs.n }

type E = struct {
	hash uint64
	len  int32
}

type aNode struct {
	key         E
	data        E
	left, right *aNode
	height      int8
	size        int32
}

func _newANode(key E) *aNode {
	return &aNode{key: key, data: key, height: 1, size: 1}
}

func (n *aNode) Balance() int8 {
	if n.left == nil {
		if n.right == nil {
			return 0
		}
		return -n.right.height
	}
	if n.right == nil {
		return n.left.height
	}
	return n.left.height - n.right.height
}

var (
	tmpPath = make([]*aNode, 0, 128)
)

type avlTreeWithSum struct {
	root *aNode
	e    func() E
	op   func(E, E) E
}

func newAVLTreeWithSum(n int32, f func(int32) E, e func() E, op func(E, E) E) *avlTreeWithSum {
	res := &avlTreeWithSum{e: e, op: op}
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func (t *avlTreeWithSum) Merge(other *avlTreeWithSum) {
	t.root = t._mergeNode(t.root, other.root)
}

func (t *avlTreeWithSum) Insert(k int32, key E) {
	n := t.Size()
	if k < 0 {
		k += n
	}
	if k < 0 {
		k = 0
	}
	if k > n {
		k = n
	}
	a, b := t._splitNode(t.root, k)
	t.root = t._mergeWithRoot(a, _newANode(key), b)
}

func (t *avlTreeWithSum) Split(k int32) (*avlTreeWithSum, *avlTreeWithSum) {
	a, b := t._splitNode(t.root, k)
	return _newWithRoot(a, t.e, t.op), _newWithRoot(b, t.e, t.op)
}

func (t *avlTreeWithSum) Pop(k int32) E {
	if k < 0 {
		k += t.Size()
	}
	a, b := t._splitNode(t.root, k+1)
	a, tmp := t._popRight(a)
	t.root = t._mergeNode(a, b)
	return tmp.key
}

func (t *avlTreeWithSum) Set(k int32, key E) {
	if k < 0 {
		k += t.Size()
	}
	node := t.root
	tmpPath = tmpPath[:0]
	path := tmpPath
	for {

		path = append(path, node)
		t := int32(0)
		if node.left != nil {
			t = node.left.size
		}
		if t == k {
			node.key = key
			break
		}
		if t < k {
			k -= t + 1
			node = node.right
		} else {
			node = node.left
		}
	}
	for i := len(path) - 1; i >= 0; i-- {
		t._update(path[i])
	}
}

func (t *avlTreeWithSum) Clear() { t.root = nil }

func (t *avlTreeWithSum) ToList() []E {
	node := t.root
	stack := make([]*aNode, 0)
	res := make([]E, 0, t.Size())
	for len(stack) > 0 || node != nil {
		if node != nil {

			stack = append(stack, node)
			node = node.left
		} else {
			node = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			res = append(res, node.key)
			node = node.right
		}
	}
	return res
}

func (t *avlTreeWithSum) Get(k int32) E {
	if k < 0 {
		k += t.Size()
	}
	node := t.root
	for {

		t := int32(0)
		if node.left != nil {
			t = node.left.size
		}
		if t == k {
			return node.key
		} else if t < k {
			k -= t + 1
			node = node.right
		} else {
			node = node.left
		}
	}
}

func (avl *avlTreeWithSum) Query(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if n := avl.Size(); end > n {
		end = n
	}
	if start >= end || avl.root == nil {
		return avl.e()
	}
	var dfs func(node *aNode, left, right int32) E
	dfs = func(node *aNode, left, right int32) E {
		if right <= start || end <= left {
			return avl.e()
		}
		if start <= left && right < end {
			return node.data
		}
		lsize := int32(0)
		if node.left != nil {
			lsize = node.left.size
		}
		res := avl.e()
		if node.left != nil {
			res = dfs(node.left, left, left+lsize)
		}
		if tmp := left + lsize; start <= tmp && tmp < end {
			res = avl.op(res, node.key)
		}
		if node.right != nil {
			res = avl.op(res, dfs(node.right, left+lsize+1, right))
		}
		return res
	}
	return dfs(avl.root, 0, avl.Size())
}

func (avl *avlTreeWithSum) QueryAll() E {
	if avl.root == nil {
		return avl.e()
	}
	return avl.root.data
}

func (t *avlTreeWithSum) Size() int32 {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func _newWithRoot(root *aNode, e func() E, op func(a, b E) E) *avlTreeWithSum {
	return &avlTreeWithSum{root: root, e: e, op: op}
}

func (t *avlTreeWithSum) _build(n int32, f func(int32) E) {
	var dfs func(l, r int32) *aNode
	dfs = func(l, r int32) *aNode {
		mid := (l + r) >> 1
		node := _newANode(f(mid))
		if l != mid {
			node.left = dfs(l, mid)
		}
		if mid+1 != r {
			node.right = dfs(mid+1, r)
		}
		t._update(node)
		return node
	}
	t.root = dfs(0, n)
}

func (t *avlTreeWithSum) _update(node *aNode) {
	node.size = 1
	node.data = node.key
	node.height = 1
	if node.left != nil {
		node.size += node.left.size
		node.data = t.op(node.left.data, node.data)
		node.height = max8(node.left.height+1, 1)
	}
	if node.right != nil {
		node.size += node.right.size
		node.data = t.op(node.data, node.right.data)
		node.height = max8(node.height, node.right.height+1)
	}
}

func (t *avlTreeWithSum) _rotateRight(node *aNode) *aNode {
	u := node.left
	node.left = u.right
	u.right = node
	t._update(node)
	t._update(u)
	return u
}

func (t *avlTreeWithSum) _rotateLeft(node *aNode) *aNode {
	u := node.right
	node.right = u.left
	u.left = node
	t._update(node)
	t._update(u)
	return u
}

func (t *avlTreeWithSum) _balanceLeft(node *aNode) *aNode {
	var u *aNode
	if node.left.left == nil || node.left.left.height+2 == node.left.height {
		u = node.left.right
		node.left.right = u.left
		u.left = node.left
		node.left = u.right
		u.right = node
		t._update(u.left)
	} else {
		u = node.left
		node.left = u.right
		u.right = node
	}
	t._update(u.right)
	t._update(u)
	return u
}

func (t *avlTreeWithSum) _balanceRight(node *aNode) *aNode {
	var u *aNode
	if node.right.right == nil || node.right.right.height+2 == node.right.height {
		u = node.right.left
		node.right.left = u.right
		u.right = node.right
		node.right = u.left
		u.left = node
		t._update(u.right)
	} else {
		u = node.right
		node.right = u.left
		u.left = node
	}
	t._update(u.left)
	t._update(u)
	return u
}

func (t *avlTreeWithSum) _mergeWithRoot(l, root, r *aNode) *aNode {
	diff := int8(0)
	if l == nil {
		if r != nil {
			diff = -r.height
		}
	} else {
		if r == nil {
			diff = l.height
		} else {
			diff = l.height - r.height
		}
	}
	if diff > 1 {
		l.right = t._mergeWithRoot(l.right, root, r)
		t._update(l)
		if l.left == nil {
			if l.right.height == 2 {
				return t._balanceRight(l)
			}
		} else {
			if l.left.height-l.right.height == -2 {
				return t._balanceRight(l)
			}
		}
		return l
	} else if diff < -1 {
		r.left = t._mergeWithRoot(l, root, r.left)
		t._update(r)
		if r.right == nil {
			if r.left.height == 2 {
				return t._balanceLeft(r)
			}
		} else {
			if r.left.height-r.right.height == 2 {
				return t._balanceLeft(r)
			}
		}
		return r
	} else {
		root.left = l
		root.right = r
		t._update(root)
		return root
	}
}

func (t *avlTreeWithSum) _mergeNode(l, r *aNode) *aNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	l, root := t._popRight(l)
	return t._mergeWithRoot(l, root, r)
}

func (t *avlTreeWithSum) _popRight(node *aNode) (*aNode, *aNode) {

	tmpPath = tmpPath[:0]
	path := tmpPath
	mx := node
	for node.right != nil {
		path = append(path, node)
		mx = node.right
		node = node.right

	}
	path = append(path, node.left)
	len_ := len(path)
	for i := 0; i < len_-1; i++ {
		node = path[len(path)-1]
		path = path[:len(path)-1]
		if node == nil {
			path[len(path)-1].right = nil
			t._update(path[len(path)-1])
			continue
		}
		b := node.Balance()
		if b == 2 {
			path[len(path)-1].right = t._balanceLeft(node)
		} else if b == -2 {
			path[len(path)-1].right = t._balanceRight(node)
		} else {
			path[len(path)-1].right = node
		}
		t._update(path[len(path)-1])
	}
	if path[0] != nil {
		b := path[0].Balance()
		if b == 2 {
			path[0] = t._balanceLeft(path[0])
		} else if b == -2 {
			path[0] = t._balanceRight(path[0])
		}
	}
	mx.left = nil
	t._update(mx)
	return path[0], mx
}

func (t *avlTreeWithSum) _splitNode(node *aNode, k int32) (*aNode, *aNode) {
	if node == nil {
		return nil, nil
	}
	tmp := k
	if node.left != nil {
		tmp -= node.left.size
	}
	if tmp == 0 {
		left := node.left
		right := t._mergeWithRoot(nil, node, node.right)
		return left, right
	}
	if tmp < 0 {
		left, right := t._splitNode(node.left, k)
		return left, t._mergeWithRoot(right, node, node.right)
	}
	left, right := t._splitNode(node.right, tmp-1)
	return t._mergeWithRoot(node.left, node, left), right
}

func max8(x, y int8) int8 {
	if x > y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

func max32(x, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func test() {
	for i := 0; i < 500; i++ {
		nums := make([]int, 500)
		for j := 0; j < 500; j++ {
			nums[j] = rand.Intn(100)
		}

		base := NewDynamicHashStringBase(int32(len(nums)), 0)
		hasher := NewDynamicHashString(int32(len(nums)), func(i int32) uint64 { return uint64(nums[i]) }, base)
		for j := 0; j < 500; j++ {
			// Set
			index := int32(rand.Intn(len(nums)))
			value := rand.Intn(100)
			hasher.Set(index, uint64(value))
			nums[index] = value

			// Get
			start, end := rand.Intn(len(nums)), rand.Intn(len(nums))
			if start > end {
				start, end = end, start
			}
			var sum uint64
			for k := start; k < end; k++ {
				sum = base.Concat(sum, uint64(nums[k]), 1)
			}
			if sum != hasher.Get(int32(start), int32(end)) {
				fmt.Println("Get error")
				panic("Get error")
			}

			// GetAll
			{
				sum = 0
				for k := 0; k < len(nums); k++ {
					sum = base.Concat(sum, uint64(nums[k]), 1)
				}
				if sum != hasher.GetAll() {
					fmt.Println("GetAll error")
					panic("GetAll error")
				}
			}

			// Pop
			index = int32(rand.Intn(len(nums)))
			if hasher.Pop(index) != uint64(nums[index]) {
				fmt.Println("Pop error")
				panic("Pop error")
			}
			nums = append(nums[:index], nums[index+1:]...)

			// Insert
			index = int32(rand.Intn(len(nums)))
			value = rand.Intn(100)
			hasher.Insert(index, uint64(value))
			nums = append(nums[:index], append([]int{value}, nums[index:]...)...)
		}
	}

	fmt.Println("pass")
}

func testTime() {
	n := int32(200000)
	nums := make([]int, n)
	for i := int32(0); i < n; i++ {
		nums[i] = rand.Intn(100)
	}
	time1 := time.Now()
	base := NewDynamicHashStringBase(int32(len(nums)), 0)
	hasher := NewDynamicHashString(int32(len(nums)), func(i int32) uint64 { return uint64(nums[i]) }, base)

	for i := int32(0); i < n; i++ {
		hasher.Set(i, uint64(nums[i]))
		hasher.Get(0, i+1)
		hasher.GetAll()
		hasher.Pop(i)
		hasher.Insert(i, uint64(nums[i]))
	}
	fmt.Println(time.Since(time1)) // 396.156ms
}
