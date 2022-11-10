// TODO
// !FIXME
// https://nyaannyaan.github.io/library/rbst/treap.hpp
// https://baobaobear.github.io/post/20191215-fhq-treap/
package treap

import (
	"fmt"
	"time"
)

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

func demo() {

}

type Data = interface{}
type Lazy = interface{}

type node struct {
	left, right int
	size        int
	priority    uint

	// alias : S
	data Data
	// alias : F
	lazy Lazy
}

// https://nyaannyaan.github.io/library/rbst/treap.hpp
type FHQTreap struct {
	seed       uint
	nodeId     int // 从1开始
	root       int
	recyc      int
	comparator func(a, b Data) int
	nodes      []node
}

// TODO Value类型不等于初始值类型???
func NewFHQTreap(comparator func(a, b Data) int, initCapacity int) *FHQTreap {
	return &FHQTreap{
		seed:       uint(time.Now().UnixNano()/2 + 1),
		comparator: comparator,
		nodes:      make([]node, max(initCapacity, 16)),
		nodeId:     1,
	}
}

// !public
func (t *FHQTreap) Size() int {
	return t.nodes[t.root].size
}

// 排名第k的元素
func (t *FHQTreap) At(index int) Data {
	n := t.Size()
	if index < 0 {
		index += n
	}
	if index < 0 || index >= n {
		panic(fmt.Sprintf("%d index out of range: [%d,%d]", index, 0, n-1))
	}
	return t.nodes[t.kthNode(index)].data
}

func (t *FHQTreap) Insert(index int, data Data) {
	t.resureCapacity(t.nodeId + 2)
	var x, y, z int
	t.split(t.root, index, &x, &y)
	z = t.newNode(data)
	t.root = t.merge(t.merge(x, z), y)
}

// 区间删除
func (t *FHQTreap) Erase(left, right int) {
	var x, y, z int
	t.split(t.root, right, &y, &z)
	t.split(y, left-1, &x, &y)
	t.remove(&y)
	t.root = t.merge(x, z)
}

// 区间修改
func (t *FHQTreap) Update(left, right int, lazy Lazy) {
	var x, y, z int
	t.split(t.root, right, &y, &z)
	t.split(y, left-1, &x, &y)
	t.add(y, lazy)
	t.root = t.merge(t.merge(x, y), z)
}

// 区间查询
func (t *FHQTreap) Query(left, right int) Data {
	var x, y, z int
	t.split(t.root, right, &y, &z)
	t.split(y, left-1, &x, &y)
	res := t.nodes[y].data
	t.root = t.merge(t.merge(x, y), z)
	return res
}

// !private
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

func (t *FHQTreap) remove(root *int) {
	if t.recyc == 0 {
		t.recyc = *root
	} else {
		t.recyc = t.merge(t.recyc, *root)
	}
	*root = 0
}

func (t *FHQTreap) add(root int, v Lazy) {
	///
}

// TODO 需要从外界传入pushDown逻辑
func (t *FHQTreap) pushDown(root int) {
	node := t.nodes[root]
	if node.lazy != 0 {
		// TODO
		t.add(node.left, node.lazy)
		t.add(node.right, node.lazy)
		node.lazy = 0
	}
}

// TODO 需要从外界传入pushUp逻辑
func (t *FHQTreap) pushUp(root int) {
	node := t.nodes[root]
	node.size = t.nodes[node.left].size + t.nodes[node.right].size + 1
	// TODO
	// node.sum= t.nodes[node.left].sum + t.nodes[node.right].sum + node.data
}

func (t *FHQTreap) updateData(root int, data *Data) {

}

// TODO recyc逻辑
func (t *FHQTreap) newNode(data Data) int {
	t.nodeId++
	index := t.nodeId
	t.nodes[index].data = data
	t.nodes[index].priority = t.fastRand()
	t.nodes[index].size = 1
	return index
}

func (t *FHQTreap) fastRand() uint {
	t.seed ^= t.seed << 13
	t.seed ^= t.seed >> 17
	t.seed ^= t.seed << 5
	return t.seed
}

func (t *FHQTreap) resureCapacity(needCap int) {
	if needCap > len(t.nodes) {
		t.resize(needCap)
	}
}

func (t *FHQTreap) resize(needCap int) {
	newCap := len(t.nodes) * 2
	if newCap < needCap {
		newCap = needCap
	}
	newNodes := make([]node, newCap)
	copy(newNodes, t.nodes)
	t.nodes = newNodes
}

func (t *FHQTreap) kthNode(k int) int {
	var x, y, z int
	t.split(t.root, k, &y, &z)
	t.split(y, k-1, &x, &y)
	res := y
	t.root = t.merge(t.merge(x, y), z)
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
