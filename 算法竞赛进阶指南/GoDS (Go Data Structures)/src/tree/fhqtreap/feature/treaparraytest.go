package main

import (
	"fmt"
	"strings"
	"time"
)

// TODO
// !FIXME
// https://nyaannyaan.github.io/library/rbst/treap.hpp
// https://baobaobear.github.io/post/20191215-fhq-treap/

// TODO
// template<typename T>
// struct FHQ_Treap
// {
//     struct data
//     {
//         T v;
//         data(int _v = 0) :v(_v) {}
//         data operator + (const data& d) const
//         {
//             data r;
//             r.v = v + d.v;
//             return r;
//         }
//         data operator * (int t) const
//         {
//             data r;
//             r.v = v * t;
//             return r;
//         }
//         operator bool() const { return v != 0; }
//         operator T() const { return v; }
//     };
//     struct node
//     {
//         int ch[2], sz;
//         unsigned k;
//         data d, sum, lz_add;
//         node(int z = 1) :sz(z), k(rnd()) { ch[0] = ch[1] = 0; }
//         static unsigned rnd()
//         {
//             static unsigned r = 0x123;
//             r = r * 69069 + 1;
//             return r;
//         }
//     };
//     vector<node> nodes;
//     int root;
//     int recyc;
//     int reserve_size;
//     void reserve()
//     {
//         if (size() >= reserve_size)
//             nodes.reserve((reserve_size *= 2) + 1);
//     }
//     inline int& ch(int tp, int r) { return nodes[tp].ch[r]; }
//     int new_node(const data& d)
//     {
//         int id = (int)nodes.size();
//         if (recyc)
//         {
//             id = recyc;
//             if (ch(recyc, 0) && ch(recyc, 1))
//                 recyc = merge(ch(recyc, 0), ch(recyc, 1));
//             else
//                 recyc = ch(recyc, 0) ? ch(recyc, 0) : ch(recyc, 1);
//             nodes[id] = node();
//         }
//         else nodes.push_back(node());
//         nodes[id].d = d;
//         nodes[id].sum = d;
//         return id;
//     }
//     int update(int tp)
//     {
//         node& n = nodes[tp];
//         n.sz = 1 + nodes[n.ch[0]].sz + nodes[n.ch[1]].sz;
//         n.sum = n.d + nodes[n.ch[0]].sum + nodes[n.ch[1]].sum;
//         return tp;
//     }
//     void add(int tp, const data& d)
//     {
//         node& n = nodes[tp];
//         n.lz_add = n.lz_add + d;
//         n.d = n.d + d;
//         n.sum = n.sum + d * n.sz;
//     }
//     void pushdown(int tp)
//     {
//         node& n = nodes[tp];
//         if (n.lz_add)
//         {
//             add(n.ch[0], n.lz_add); add(n.ch[1], n.lz_add);
//             n.lz_add = 0;
//         }
//     }
//     int merge(int tl, int tr)
//     {
//         if (!tl) return tr;
//         else if (!tr) return tl;
//         if (nodes[tl].k < nodes[tr].k)
//         {
//             pushdown(tl);
//             ch(tl, 1) = merge(ch(tl, 1), tr);
//             return update(tl);
//         }
//         else
//         {
//             pushdown(tr);
//             ch(tr, 0) = merge(tl, ch(tr, 0));
//             return update(tr);
//         }
//     }
//     void split(int tp, int k, int &x, int &y)
//     {
//         if (!tp) { x = y = 0; return; }
//         pushdown(tp);
//         if (k <= nodes[ch(tp, 0)].sz)
//         {
//             y = tp;
//             split(ch(tp, 0), k, x, ch(tp, 0));
//             update(y);
//         }
//         else
//         {
//             x = tp;
//             split(ch(tp, 1), k - nodes[ch(tp, 0)].sz - 1, ch(tp, 1), y);
//             update(x);
//         }
//     }
//     void remove(int& tp)
//     {
//         if (recyc == 0) recyc = tp;
//         else recyc = merge(recyc, tp);
//         tp = 0;
//     }
//     // interface
//     void init(int size)
//     {
//         nodes.clear();
//         nodes.reserve((size = max(size, 15)) + 1);
//         nodes.push_back(node(0));
//         root = 0;
//         recyc = 0; reserve_size = size + 1;
//     }
//     T get(int id) { return nodes[id].d; }
//     int size() { return nodes[root].sz; }
//     int kth(int k)
//     {
//         int x, y, z;
//         split(root, k, y, z); split(y, k - 1, x, y);
//         int id = y;
//         root = merge(merge(x, y), z);
//         return id;
//     }
//     void insert(int k, data v)
//     {
//         int l, r;
//         split(root, k - 1, l, r);
//         int tp = new_node(v);
//         root = merge(merge(l, tp), r);
//     }
//     void erase(int l, int r)
//     {
//         int x, y, z;
//         split(root, r, y, z); split(y, l - 1, x, y);
//         remove(y);
//         root = merge(x, z);
//     }
//     void range_add(int l, int r, data v)
//     {
//         int x, y, z;
//         split(root, r, y, z); split(y, l - 1, x, y);
//         add(y, v);
//         root = merge(merge(x, y), z);
//     }
//     T getsum(int l, int r)
//     {
//         int x, y, z;
//         split(root, r, y, z); split(y, l - 1, x, y);
//         T ret = nodes[y].sum;
//         root = merge(merge(x, y), z);
//         return ret;
//     }
// }

func main() {
	treapArray := NewFHQTreap(100)
	nums := []int{}
	fmt.Println(treapArray.Size())
	for i := 0; i < 10; i++ {
		nums = append(nums, i+12)
	}
	for i := 0; i < 10; i++ {
		treapArray.Append(nums[i])
	}
	fmt.Println("num:", nums)
	fmt.Println(treapArray)
	fmt.Println(treapArray.Size(), treapArray.At(-1))
	fmt.Println(treapArray.Size(), treapArray.At(treapArray.Size()-1))
	treapArray.Insert(-1000, 100)
	fmt.Println(treapArray)
	fmt.Println(treapArray.Pop(0))
	fmt.Println(treapArray)
	treapArray.Erase(0, 1)
	fmt.Println(treapArray)
	fmt.Println(treapArray.At(0))
}

// An arraylist impleted by FHQTreap.
//
// Author:
// https://github.com/981377660LMT/algorithm-study
//
// Reference:
// https://baobaobear.github.io/post/20191215-fhq-treap/
// https://nyaannyaan.github.io/library/rbst/treap.hpp
// https://github.com/EndlessCheng/codeforces-go/blob/f9d97465d8b351af7536b5b6dac30b220ba1b913/copypasta/treap.go

type Raw = interface{}
type Data = interface{}
type Lazy = interface{}

type Node struct {
	// !Raw value
	raw int

	// !Data and lazy tag maintained by segment tree
	sum     int
	lazyAdd int

	// FHQTreap inner attributes
	left, right int
	size        int
	priority    uint
}

// !Segment tree op.Need to be modified according to the actual situation.
func (t *FHQTreap) pushUp(root int) {
	node := t.nodes[root]
	node.size = t.nodes[node.left].size + t.nodes[node.right].size + 1
	node.sum = t.nodes[node.left].sum + t.nodes[node.right].sum + node.raw
}

// !Segment tree functions.Need to be modified according to the actual situation.
func (t *FHQTreap) _pushDown(root int, delta int) {
	node := t.nodes[root]
	node.raw += delta
	node.sum += delta * node.size
	node.lazyAdd += delta
}

// !Segment tree functions.Need to be modified according to the actual situation.
func (t *FHQTreap) pushDown(root int) {
	node := t.nodes[root]
	if node.lazyAdd != 0 {
		delta := node.lazyAdd
		t._pushDown(node.left, delta)
		t._pushDown(node.right, delta)
		node.lazyAdd = 0
	}
}

type FHQTreap struct {
	seed  uint
	root  int
	nodes []*Node
}

// Need to be modified according to the actual situation to implement a segment tree.
func NewFHQTreap(initCapacity int) *FHQTreap {
	treap := &FHQTreap{
		seed:  uint(time.Now().UnixNano()/2 + 1),
		nodes: make([]*Node, 0, max(initCapacity, 16)),
	}

	// 0号表示dummy结点
	dummy := &Node{
		size:     0,
		priority: treap.fastRand(),
	}
	treap.nodes = append(treap.nodes, dummy)
	return treap
}

// Return the number of items in the list.
func (t *FHQTreap) Size() int {
	return t.nodes[t.root].size
}

// Return the value at the k-th position (0-indexed).
func (t *FHQTreap) At(index int) int {
	n := t.Size()
	if index < 0 {
		index += n
	}

	if index < 0 || index >= n {
		panic(fmt.Sprintf("index %d out of range [0,%d]", index, n-1))
	}

	index += 1 // dummy offset
	var x, y, z int
	t.splitByRank(t.root, index, &y, &z)
	t.splitByRank(y, index-1, &x, &y)
	res := &t.nodes[y].raw
	t.root = t.merge(t.merge(x, y), z)
	return *res
}

// Append value to the end of the list.
func (t *FHQTreap) Append(value int) {
	t.Insert(t.Size(), value)
}

// Insert value before index.
func (t *FHQTreap) Insert(index int, value int) {
	n := t.Size()
	if index < 0 {
		index += n
	}

	index += 1 // dummy offset
	var x, y, z int
	t.splitByRank(t.root, index-1, &x, &y)
	z = t.newNode(value)
	t.root = t.merge(t.merge(x, z), y)
}

// Remove and return item at index.
func (t *FHQTreap) Pop(index int) int {
	n := t.Size()
	if index < 0 {
		index += n
	}

	index += 1 // dummy offset
	var x, y, z int
	t.splitByRank(t.root, index, &y, &z)
	t.splitByRank(y, index-1, &x, &y)
	res := &t.nodes[y].raw
	t.root = t.merge(x, z)
	return *res
}

// Remove [start, stop) from list.
func (t *FHQTreap) Erase(start, stop int) {
	var x, y, z int
	start++ // dummy offset
	t.splitByRank(t.root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	t.root = t.merge(x, z)
}

// Reverse [start, stop) in place.
func (t *FHQTreap) Reverse(start, stop int) {
	panic("not implemented")
}

// Update [start, stop) with value (defaults to range add).
//  0 <= start <= stop <= n
func (t *FHQTreap) Update(start, stop int, delta int) {
	var x, y, z int
	t.splitByRank(t.root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	t._pushDown(y, delta)
	t.root = t.merge(t.merge(x, y), z)
}

// Query [start, stop) (defaults to range sum).
//  0 <= start <= stop <= n
func (t *FHQTreap) Query(start, stop int) int {
	var x, y, z int
	t.splitByRank(t.root, stop, &y, &z)
	t.splitByRank(y, start-1, &x, &y)
	res := t.nodes[y].sum
	t.root = t.merge(t.merge(x, y), z)
	return res
}

// Query [0, n) (defaults to range sum).
func (t *FHQTreap) QueryAll() int {
	return t.nodes[t.root].sum
}

// Split by rank.
// Split the tree rooted at root into two trees, x and y, such that the size of x is k.
// x is the left subtree, y is the right subtree.
func (t *FHQTreap) splitByRank(root, k int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	t.pushDown(root)
	if k <= t.nodes[t.nodes[root].left].size {
		*y = root
		t.splitByRank(t.nodes[root].left, k, x, &t.nodes[root].left)
		t.pushUp(*y)
	} else {
		*x = root
		t.splitByRank(t.nodes[root].right, k-t.nodes[t.nodes[root].left].size-1, &t.nodes[root].right, y)
		t.pushUp(*x)
	}
}

// Make sure that the height of the resulting tree is at most O(log n).
// A random priority is introduced to determine who is the root after merge operation.
// If left subtree is smaller, merge right subtree with the right child of the left subtree.
// Otherwise, merge left subtree with the left child of the right subtree.
func (t *FHQTreap) merge(x, y int) int {
	if x == 0 {
		return y
	}
	if y == 0 {
		return x
	}

	if t.nodes[x].priority < t.nodes[y].priority {
		t.pushDown(x)
		t.nodes[x].right = t.merge(t.nodes[x].right, y)
		t.pushUp(x)
		return x
	} else {
		t.pushDown(y)
		t.nodes[y].left = t.merge(x, t.nodes[y].left)
		t.pushUp(y)
		return y
	}
}

// Add a new node and return its nodeId.
func (t *FHQTreap) newNode(data int) int {
	node := &Node{
		size:     1,
		priority: t.fastRand(),
		raw:      data,
		sum:      data,
	}
	t.nodes = append(t.nodes, node)
	return len(t.nodes) - 1
}

func (t *FHQTreap) String() string {
	sb := []string{"TreapArray{"}
	values := []string{}
	for i := 0; i < t.Size(); i++ {
		values = append(values, fmt.Sprintf("%d", t.At(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")

}

// https://github.com/EndlessCheng/codeforces-go/blob/f9d97465d8b351af7536b5b6dac30b220ba1b913/copypasta/treap.go#L31
func (t *FHQTreap) fastRand() uint {
	t.seed ^= t.seed << 13
	t.seed ^= t.seed >> 17
	t.seed ^= t.seed << 5
	return t.seed
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
