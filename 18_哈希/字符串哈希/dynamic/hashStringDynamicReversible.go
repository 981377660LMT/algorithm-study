// dynamicHashString/hashStringDynamic
// 支持区间翻转的动态哈希.
//
// api:
//  1. Insert(i, c): 在第i个位置插入字符c.
//  2. Pop(i): 删除第i个位置的字符.
//  3. Set(i, c): 将第i个位置的字符设置为c.
//  4. Get(start, end): 获取[start, end)的哈希值.
//  5. GetAll(): 获取所有字符的哈希值.
//  6. Reverse(start, end): 将[start, end)的字符反转.
//  7. ReverseAll(): 将所有字符反转.

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
	base := NewDynamicHashStringBase(len(s), 37)
	hs := NewDynamicHashStringReverible(len(s), func(i int) uint { return uint(s[i]) }, base)
	fmt.Println(hs.Get(0, 1))
	fmt.Println(hs.Get(1, 2))
	fmt.Println(hs.Get(2, 3))
	hs.Set(0, 1)
	fmt.Println(hs.Get(0, 1))
	hs.Reverse(0, 3)
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

	var n, q int
	fmt.Fscan(in, &n, &q)
	var s string
	fmt.Fscan(in, &s)

	base := NewDynamicHashStringBase(len(s), 0)
	hasher1 := NewDynamicHashStringReverible(len(s), func(i int) uint { return uint(s[i]) }, base)

	isPalindrome := func(start, end int) bool {
		hash1 := hasher1.Get(start, end)
		hasher1.Reverse(start, end)
		hash2 := hasher1.Get(start, end)
		hasher1.Reverse(start, end)
		return hash1 == hash2
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var pos int
			var c string
			fmt.Fscan(in, &pos, &c)
			pos--
			hasher1.Set(pos, uint(c[0]))
		} else {
			var l, r int
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
	n := len(s)
	base := NewDynamicHashStringBase(n, 0)
	hasher := NewDynamicHashStringReverible(n, func(i int) uint { return uint(s[i]) }, base)
	countPre := func(curLen, start int) int {
		left, right := 1, curLen
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
	for i := 1; i < n+1; i++ {
		if s[0] != s[n-i] {
			continue
		}
		count := countPre(i, n-i)
		res += count
	}

	return int64(res)

}

const (
	hashStringMod    uint = (1 << 61) - 1
	hashStringMask30 uint = (1 << 30) - 1
	hashStringMask31 uint = (1 << 31) - 1
	hashStringMASK61 uint = hashStringMod
)

type DynamicHashStringBase struct {
	n    int
	powb []uint
}

// base: 0 表示随机生成
func NewDynamicHashStringBase(n int, base uint) *DynamicHashStringBase {
	res := &DynamicHashStringBase{}
	if base == 0 {
		base = uint(37 + rand.Intn(1e9))
	}
	powb := make([]uint, n+1)
	powb[0] = 1
	for i := 1; i <= n; i++ {
		powb[i] = res.Mul(powb[i-1], base)
	}
	res.n = n
	res.powb = powb
	return res
}

// h1 <- h2. len(h2) == k.
func (hsb *DynamicHashStringBase) Concat(h1, h2, h2Len uint) uint {
	return hsb.Mod(hsb.Mul(h1, hsb.powb[h2Len]) + h2)
}

// a*b % (2^61-1)
func (hsb *DynamicHashStringBase) Mul(a, b uint) uint {
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
func (hsb *DynamicHashStringBase) Mod(x uint) uint {
	xu := x >> 61
	xd := x & hashStringMASK61
	res := xu + xd
	if res >= hashStringMod {
		res -= hashStringMod
	}
	return res
}

type DynamicHashStringReversible struct {
	n    int
	base *DynamicHashStringBase
	avl  *avlTreeWithSumReversibleNoncommutative
}

func NewDynamicHashStringReverible(n int, f func(i int) uint, base *DynamicHashStringBase) *DynamicHashStringReversible {
	res := &DynamicHashStringReversible{n: n, base: base}
	res.avl = newAVLTreeWithSumReversibleNoncommutative(
		int32(n), func(i int32) E { return E{hash: f(int(i)), len: 1} },
		func() E { return E{} },
		func(a, b E) E { return E{hash: base.Concat(a.hash, b.hash, uint(b.len)), len: a.len + b.len} },
	)
	return res
}

func (hs *DynamicHashStringReversible) Insert(index int, c uint) {
	hs.avl.Insert(int32(index), E{hash: c, len: 1})
	hs.n++
}

func (hs *DynamicHashStringReversible) Pop(index int) uint {
	hs.n--
	return hs.avl.Pop(int32(index)).hash
}

func (hs *DynamicHashStringReversible) Reverse(start, end int) {
	if start < 0 {
		start = 0
	}
	if end > hs.n {
		end = hs.n
	}
	if start >= end {
		return
	}
	hs.avl.Reverse(int32(start), int32(end))
}

func (hs *DynamicHashStringReversible) ReverseAll() { hs.avl.ReverseAll() }

func (hs *DynamicHashStringReversible) Get(start, end int) uint {
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
	return hs.avl.Query(int32(start), int32(end)).hash
}

func (hs *DynamicHashStringReversible) GetAll() uint { return hs.avl.QueryAll().hash }

func (hs *DynamicHashStringReversible) Set(index int, c uint) {
	if index < 0 || index >= hs.n {
		return
	}
	hs.avl.Set(int32(index), E{hash: c, len: 1})
}

func (hs *DynamicHashStringReversible) Len() int { return hs.n }

type E = struct {
	hash uint
	len  int32
}

type aNode struct {
	key         E
	data        E
	revData     E
	left, right *aNode
	height      int8
	size        int32
	rev         bool
}

func _newANode(key E) *aNode {
	return &aNode{key: key, data: key, revData: key, height: 1, size: 1}
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

func (n *aNode) String() string {
	if n.left == nil && n.right == nil {
		return fmt.Sprintf("key=%v, height=%v, size=%v\n", n.key, n.height, n.size)
	}
	return fmt.Sprintf("key=%v, height=%v, size=%v,\n left:%v,\n right:%v\n", n.key, n.height, n.size, n.left, n.right)
}

var (
	tmpPath = make([]*aNode, 0, 128)
)

type avlTreeWithSumReversibleNoncommutative struct {
	root *aNode
	e    func() E
	op   func(E, E) E
}

func newAVLTreeWithSumReversibleNoncommutative(n int32, f func(int32) E, e func() E, op func(E, E) E) *avlTreeWithSumReversibleNoncommutative {
	res := &avlTreeWithSumReversibleNoncommutative{e: e, op: op}
	if n > 0 {
		res._build(n, f)
	}
	return res
}

func (t *avlTreeWithSumReversibleNoncommutative) Merge(other *avlTreeWithSumReversibleNoncommutative) {
	t.root = t._mergeNode(t.root, other.root)
}

func (t *avlTreeWithSumReversibleNoncommutative) Insert(k int32, key E) {
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

func (t *avlTreeWithSumReversibleNoncommutative) Split(k int32) (*avlTreeWithSumReversibleNoncommutative, *avlTreeWithSumReversibleNoncommutative) {
	a, b := t._splitNode(t.root, k)
	return _newWithRoot(a, t.e, t.op), _newWithRoot(b, t.e, t.op)
}

func (t *avlTreeWithSumReversibleNoncommutative) Pop(k int32) E {
	if k < 0 {
		k += t.Size()
	}
	a, b := t._splitNode(t.root, k+1)
	a, tmp := t._popRight(a)
	t.root = t._mergeNode(a, b)
	return tmp.key
}

func (t *avlTreeWithSumReversibleNoncommutative) Set(k int32, key E) {
	if k < 0 {
		k += t.Size()
	}
	node := t.root
	tmpPath = tmpPath[:0]
	path := tmpPath
	for {
		t._propagate(node)
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

func (t *avlTreeWithSumReversibleNoncommutative) Clear() { t.root = nil }

func (t *avlTreeWithSumReversibleNoncommutative) ToList() []E {
	node := t.root
	stack := make([]*aNode, 0)
	res := make([]E, 0, t.Size())
	for len(stack) > 0 || node != nil {
		if node != nil {
			t._propagate(node)
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

func (t *avlTreeWithSumReversibleNoncommutative) Get(k int32) E {
	if k < 0 {
		k += t.Size()
	}
	node := t.root
	for {
		t._propagate(node)
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
func (avl *avlTreeWithSumReversibleNoncommutative) Reverse(start, end int32) {
	if start < 0 {
		start = 0
	}
	if n := avl.Size(); end > n {
		end = n
	}
	if start >= end {
		return
	}
	if start == 0 && end == avl.Size() {
		avl.ReverseAll()
		return
	}
	s, t := avl._splitNode(avl.root, end)
	r, s := avl._splitNode(s, start)
	s.rev = !s.rev
	s.data, s.revData = s.revData, s.data
	avl.root = avl._mergeNode(avl._mergeNode(r, s), t)
}

func (avl *avlTreeWithSumReversibleNoncommutative) ReverseAll() {
	if avl.root == nil {
		return
	}
	avl.root.rev = !avl.root.rev
}

func (avl *avlTreeWithSumReversibleNoncommutative) Query(start, end int32) E {
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
		avl._propagate(node)
		if start <= left && right < end {
			if node.rev {
				return node.revData
			}
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

func (avl *avlTreeWithSumReversibleNoncommutative) QueryAll() E {
	if avl.root == nil {
		return avl.e()
	}
	return avl.root.data
}

func (t *avlTreeWithSumReversibleNoncommutative) Size() int32 {
	if t.root == nil {
		return 0
	}
	return t.root.size
}

func _newWithRoot(root *aNode, e func() E, op func(a, b E) E) *avlTreeWithSumReversibleNoncommutative {
	return &avlTreeWithSumReversibleNoncommutative{root: root, e: e, op: op}
}

func (t *avlTreeWithSumReversibleNoncommutative) _build(n int32, f func(int32) E) {
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

func (t *avlTreeWithSumReversibleNoncommutative) _update(node *aNode) {
	node.size = 1
	node.data = node.key
	node.revData = node.key
	node.height = 1
	if left := node.left; left != nil {
		node.size += left.size
		node.data = t.op(left.data, node.data)
		node.revData = t.op(node.revData, left.revData)
		node.height = max8(left.height+1, 1)
	}
	if right := node.right; right != nil {
		node.size += right.size
		node.data = t.op(node.data, right.data)
		node.revData = t.op(right.revData, node.revData)
		node.height = max8(node.height, right.height+1)
	}
}

func (t *avlTreeWithSumReversibleNoncommutative) _rotateRight(node *aNode) *aNode {
	u := node.left
	node.left = u.right
	u.right = node
	t._update(node)
	t._update(u)
	return u
}

func (t *avlTreeWithSumReversibleNoncommutative) _rotateLeft(node *aNode) *aNode {
	u := node.right
	node.right = u.left
	u.left = node
	t._update(node)
	t._update(u)
	return u
}

func (t *avlTreeWithSumReversibleNoncommutative) _balanceLeft(node *aNode) *aNode {
	t._propagate(node.left)
	var u *aNode
	if node.left.left == nil || node.left.left.height+2 == node.left.height {
		u = node.left.right
		t._propagate(u)
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

func (t *avlTreeWithSumReversibleNoncommutative) _balanceRight(node *aNode) *aNode {
	t._propagate(node.right)
	var u *aNode
	if node.right.right == nil || node.right.right.height+2 == node.right.height {
		u = node.right.left
		t._propagate(u)
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

func (t *avlTreeWithSumReversibleNoncommutative) _mergeWithRoot(l, root, r *aNode) *aNode {
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
		t._propagate(l)
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
		t._propagate(r)
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

func (t *avlTreeWithSumReversibleNoncommutative) _mergeNode(l, r *aNode) *aNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	l, root := t._popRight(l)
	return t._mergeWithRoot(l, root, r)
}

func (t *avlTreeWithSumReversibleNoncommutative) _popRight(node *aNode) (*aNode, *aNode) {
	t._propagate(node)
	tmpPath = tmpPath[:0]
	path := tmpPath
	mx := node
	for node.right != nil {
		path = append(path, node)
		mx = node.right
		node = node.right
		t._propagate(node)
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

func (t *avlTreeWithSumReversibleNoncommutative) _splitNode(node *aNode, k int32) (*aNode, *aNode) {
	if node == nil {
		return nil, nil
	}
	t._propagate(node)
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

func (avl *avlTreeWithSumReversibleNoncommutative) _propagate(node *aNode) {
	if node == nil {
		return
	}
	l, r := node.left, node.right
	if node.rev {
		node.left, node.right = r, l
		if l != nil {
			l.rev = !l.rev
			l.data, l.revData = l.revData, l.data
		}
		if r != nil {
			r.rev = !r.rev
			r.data, r.revData = r.revData, r.data
		}
		node.rev = false
	}
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
		base := NewDynamicHashStringBase(len(nums), 0)
		hasher := NewDynamicHashStringReverible(len(nums), func(i int) uint { return uint(nums[i]) }, base)
		for j := 0; j < 500; j++ {
			// Set
			index := rand.Intn(len(nums))
			value := rand.Intn(100)
			hasher.Set(index, uint(value))
			nums[index] = value

			// Get
			start, end := rand.Intn(len(nums)), rand.Intn(len(nums))
			if start > end {
				start, end = end, start
			}
			var sum uint
			for k := start; k < end; k++ {
				sum = base.Concat(sum, uint(nums[k]), 1)
			}
			if sum != hasher.Get(start, end) {
				fmt.Println("Get error")
				panic("Get error")
			}

			// GetAll
			{
				sum = 0
				for k := 0; k < len(nums); k++ {
					sum = base.Concat(sum, uint(nums[k]), 1)
				}
				if sum != hasher.GetAll() {
					fmt.Println("GetAll error")
					panic("GetAll error")
				}
			}

			// Reverse
			start, end = rand.Intn(len(nums)), rand.Intn(len(nums))
			if start > end {
				start, end = end, start
			}
			hasher.Reverse(start, end)
			for i, j := start, end-1; i < j; i, j = i+1, j-1 {
				nums[i], nums[j] = nums[j], nums[i]
			}

			// ReverseAll
			{
				hasher.ReverseAll()
				for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
					nums[i], nums[j] = nums[j], nums[i]
				}
			}

			// Pop
			index = rand.Intn(len(nums))
			if hasher.Pop(index) != uint(nums[index]) {
				fmt.Println("Pop error")
				panic("Pop error")
			}
			nums = append(nums[:index], nums[index+1:]...)

			// Insert
			index = rand.Intn(len(nums))
			value = rand.Intn(100)
			hasher.Insert(index, uint(value))
			nums = append(nums[:index], append([]int{value}, nums[index:]...)...)
		}
	}

	fmt.Println("pass")
}

func testTime() {
	n := 200000
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = rand.Intn(100)
	}
	time1 := time.Now()
	base := NewDynamicHashStringBase(len(nums), 0)
	hasher := NewDynamicHashStringReverible(len(nums), func(i int) uint { return uint(nums[i]) }, base)

	for i := 0; i < n; i++ {
		hasher.Set(i, uint(nums[i]))
		hasher.Get(0, i+1)
		hasher.GetAll()
		hasher.Reverse(0, i+1)
		hasher.ReverseAll()
		hasher.Pop(i)
		hasher.Insert(i, uint(nums[i]))
	}
	fmt.Println(time.Since(time1)) // 751.756708ms
}
