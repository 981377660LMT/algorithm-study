package main

import (
	"fmt"
	"math/rand"
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

// !谢谢你 baobaobear
// 需要更多test
func main() {
	treapArray := NewFHQTreap(100)
	fmt.Println(treapArray.Size())
	for i := 0; i < 10; i++ {
		treapArray.Append(rand.Intn(100))
	}

	fmt.Println(treapArray)
	treapArray.Erase(1, 3)
	fmt.Println(treapArray)
	fmt.Println(treapArray.Query(1, 3), treapArray.At(1), treapArray.Size())

	treapArray.Insert(0, 110)
	fmt.Println(treapArray)
	treapArray.Insert(9, 110) // 插入到index位置
	fmt.Println(treapArray)
}

type Data = interface{}
type Lazy = interface{}

type Node struct {
	left, right int
	size        int
	priority    uint

	data int

	// 线段树部分
	sum     int
	lazyAdd int
}

// https://nyaannyaan.github.io/library/rbst/treap.hpp
// https://baobaobear.github.io/post/20191215-fhq-treap/
type FHQTreap struct {
	seed  uint
	root  int
	nodes []Node
}

// 要实现线段树功能 需要改写pushUp和pushDown 函数 并在node中加入不同的lazy数据
func NewFHQTreap(initCapacity int) *FHQTreap {
	treap := &FHQTreap{
		seed:  uint(time.Now().UnixNano()/2 + 1),
		nodes: make([]Node, 0, max(initCapacity, 16)),
	}

	// 0号表示dummy结点
	dummy := &Node{
		size:     0,
		priority: treap.fastRand(),
	}
	treap.nodes = append(treap.nodes, *dummy)
	return treap
}

func (t *FHQTreap) Size() int {
	return t.nodes[t.root].size
}

func (t *FHQTreap) At(index int) int {
	var x, y, z int
	t.split(t.root, index, &y, &z)
	t.split(y, index-1, &x, &y)
	res := &t.nodes[y].data
	t.root = t.merge(t.merge(x, y), z)
	return *res
}

func (t *FHQTreap) Append(value int) {
	t.Insert(t.Size(), value)
}

func (t *FHQTreap) Insert(index int, data int) {
	var x, y, z int
	t.split(t.root, index-1, &x, &y)
	z = t.newNode(data)
	t.root = t.merge(t.merge(x, z), y)
}

// 区间删除
func (t *FHQTreap) Erase(left, right int) {
	var x, y, z int
	t.split(t.root, right, &y, &z)
	t.split(y, left-1, &x, &y)
	t.root = t.merge(x, z)
}

// 区间加上delta
func (t *FHQTreap) Update(left, right int, delta int) {
	var x, y, z int
	t.split(t.root, right, &y, &z)
	t.split(y, left-1, &x, &y)
	t._pushDown(y, delta)
	t.root = t.merge(t.merge(x, y), z)
}

// 区间查询
func (t *FHQTreap) Query(left, right int) int {
	var x, y, z int
	t.split(t.root, right, &y, &z)
	t.split(y, left-1, &x, &y)
	res := t.nodes[y].sum
	t.root = t.merge(t.merge(x, y), z)
	return res
}

// 添加新结点，返回新结点的nodeId
func (t *FHQTreap) newNode(data int) int {
	node := &Node{
		size:     1,
		priority: t.fastRand(),
		data:     data,
		sum:      data,
	}
	t.nodes = append(t.nodes, *node)
	return len(t.nodes) - 1
}

func (t *FHQTreap) pushUp(root int) {
	node := &t.nodes[root]
	node.size = t.nodes[node.left].size + t.nodes[node.right].size + 1
	node.sum = t.nodes[node.left].sum + t.nodes[node.right].sum + node.data
}

func (t *FHQTreap) _pushDown(root int, delta int) {
	node := &t.nodes[root]
	node.data += delta
	node.sum += delta * node.size
	node.lazyAdd += delta
}

func (t *FHQTreap) pushDown(root int) {
	node := &t.nodes[root]
	if node.lazyAdd != 0 {
		delta := node.lazyAdd
		t._pushDown(node.left, delta)
		t._pushDown(node.right, delta)
		node.lazyAdd = 0
	}
}

// 根据排名 k 分裂 root。
//  x 是返回的左子树, y 是返回的右子树。
//  如果划分点在左子树，那么y一定是根，反之划分点在右子树，那么x一定是根。
func (t *FHQTreap) split(root, k int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	t.pushDown(root)
	if k <= t.nodes[t.nodes[root].left].size {
		*y = root
		t.split(t.nodes[root].left, k, x, &t.nodes[root].left)
		t.pushUp(*y)
	} else {
		*x = root
		t.split(t.nodes[root].right, k-t.nodes[t.nodes[root].left].size-1, &t.nodes[root].right, y)
		t.pushUp(*x)
	}
}

// 不能直接把右子树直接接在左子树的最后一个元素后，这样会导致树高度太大。
// 在Treap里面，引入了一个随机值，来决定谁来做根节点，
// 所以，我们就对比这个值，如果左子树的小，
// 那么就让左子树的右儿子与右子树merge，否则就让右子树的左儿子与左子树merge。
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

func (t *FHQTreap) String() string {
	sb := []string{"TreapArray{"}
	values := []string{}
	for i := 0; i < t.Size(); i++ {
		values = append(values, fmt.Sprintf("%d", t.At(i)))
	}
	sb = append(sb, strings.Join(values, ","), "}")
	return strings.Join(sb, "")

}

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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
