// #pragma once
// // Node 型を別に定義して使う
// template <typename Node>
// struct SplayTree {
//   Node *pool;
//   const int NODES;
//   int pid;
//   using np = Node *;
//   using X = typename Node::value_type;
//   using A = typename Node::operator_type;
//   vc<np> FREE;

//   SplayTree(int NODES) : NODES(NODES), pid(0) { pool = new Node[NODES]; }
//   ~SplayTree() { delete[] pool; }

//   void free_subtree(np c) {
//     auto dfs = [&](auto &dfs, np c) -> void {
//       if (c->l) dfs(dfs, c->l);
//       if (c->r) dfs(dfs, c->r);
//       FREE.eb(c);
//     };
//     dfs(dfs, c);
//   }

//   void reset() {
//     pid = 0;
//     FREE.clear();
//   }

//   np new_root() { return nullptr; }

//   np new_node(const X &x) {
//     assert(!FREE.empty() || pid < NODES);
//     np n = (FREE.empty() ? &(pool[pid++]) : POP(FREE));
//     Node::new_node(n, x);
//     return n;
//   }

//   np new_node(const vc<X> &dat) {
//     auto dfs = [&](auto &dfs, int l, int r) -> np {
//       if (l == r) return nullptr;
//       if (r == l + 1) return new_node(dat[l]);
//       int m = (l + r) / 2;
//       np l_root = dfs(dfs, l, m);
//       np r_root = dfs(dfs, m + 1, r);
//       np root = new_node(dat[m]);
//       root->l = l_root, root->r = r_root;
//       if (l_root) l_root->p = root;
//       if (r_root) r_root->p = root;
//       root->nodeUpdate();
//       return root;
//     };
//     return dfs(dfs, 0, len(dat));
//   }

//   u32 get_size(np root) { return (root ? root->size : 0); }

//   np merge(np l_root, np r_root) {
//     if (!l_root) return r_root;
//     if (!r_root) return l_root;
//     assert((!l_root->p) && (!r_root->p));
//     splay_kth(r_root, 0); // splay したので nodeProp 済
//     r_root->l = l_root;
//     l_root->p = r_root;
//     r_root->nodeUpdate();
//     return r_root;
//   }
//   np merge3(np a, np b, np c) { return merge(merge(a, b), c); }
//   np merge4(np a, np b, np c, np d) { return merge(merge(merge(a, b), c), d); }

//   pair<np, np> split(np root, u32 k) {
//     assert(!root || !root->p);
//     if (k == 0) return {nullptr, root};
//     if (k == (root->size)) return {root, nullptr};
//     splay_kth(root, k - 1);
//     np right = root->r;
//     root->r = nullptr, right->p = nullptr;
//     root->nodeUpdate();
//     return {root, right};
//   }
//   tuple<np, np, np> split3(np root, u32 l, u32 r) {
//     np nm, nr;
//     tie(root, nr) = split(root, r);
//     tie(root, nm) = split(root, l);
//     return {root, nm, nr};
//   }
//   tuple<np, np, np, np> split4(np root, u32 i, u32 j, u32 k) {
//     np d;
//     tie(root, d) = split(root, k);
//     auto [a, b, c] = split3(root, i, j);
//     return {a, b, c, d};
//   }

//   // 部分木が区間 [l,r) に対応するようなノードを作って返す
//   // そのノードが root になるわけではないので、
//   // このノードを参照した後にすぐに splay して根に持ち上げること
//   void goto_between(np &root, u32 l, u32 r) {
//     if (l == 0 && r == root->size) return;
//     if (l == 0) {
//       splay_kth(root, r);
//       root = root->l;
//       return;
//     }
//     if (r == root->size) {
//       splay_kth(root, l - 1);
//       root = root->r;
//       return;
//     }
//     splay_kth(root, r);
//     np rp = root;
//     root = rp->l;
//     root->p = nullptr;
//     splay_kth(root, l - 1);
//     root->p = rp;
//     rp->l = root;
//     rp->nodeUpdate();
//     root = root->r;
//   }

//   vc<X> get_all(const np &root) {
//     vc<X> res;
//     auto dfs = [&](auto &dfs, np root) -> void {
//       if (!root) return;
//       root->nodeProp();
//       dfs(dfs, root->l);
//       res.eb(root->nodeGet());
//       dfs(dfs, root->r);
//     };
//     dfs(dfs, root);
//     return res;
//   }

//   X nodeGet(np &root, u32 k) {
//     assert(root == nullptr || !root->p);
//     splay_kth(root, k);
//     return root->nodeGet();
//   }

//   void nodeSet(np &root, u32 k, const X &x) {
//     assert(root != nullptr && !root->p);
//     splay_kth(root, k);
//     root->nodeSet(x);
//   }

//   void multiply(np &root, u32 k, const X &x) {
//     assert(root != nullptr && !root->p);
//     splay_kth(root, k);
//     root->multiply(x);
//   }

//   X prod(np &root, u32 l, u32 r) {
//     assert(root == nullptr || !root->p);
//     using Mono = typename Node::Monoid_X;
//     if (l == r) return Mono::unit();
//     assert(0 <= l && l < r && r <= root->size);
//     goto_between(root, l, r);
//     X res = root->prod;
//     splay(root, true);
//     return res;
//   }

//   X prod(np &root) {
//     assert(root == nullptr || !root->p);
//     using Mono = typename Node::Monoid_X;
//     return (root ? root->prod : Mono::unit());
//   }

//   void apply(np &root, u32 l, u32 r, const A &a) {
//     if (l == r) return;
//     assert(0 <= l && l < r && r <= root->size);
//     goto_between(root, l, r);
//     root->apply(a);
//     splay(root, true);
//   }
//   void apply(np &root, const A &a) {
//     if (!root) return;
//     root->apply(a);
//   }

//   void nodeReverse(np &root, u32 l, u32 r) {
//     assert(root == nullptr || !root->p);
//     if (l == r) return;
//     assert(0 <= l && l < r && r <= root->size);
//     goto_between(root, l, r);
//     root->nodeReverse();
//     splay(root, true);
//   }
//   void nodeReverse(np root) {
//     if (!root) return;
//     root->nodeReverse();
//   }

//   void rotate(Node *n) {
//     // n を根に近づける。prop, nodeUpdate は rotate の外で行う。
//     Node *pp, *p, *c;
//     p = n->p;
//     pp = p->p;
//     if (p->l == n) {
//       c = n->r;
//       n->r = p;
//       p->l = c;
//     } else {
//       c = n->l;
//       n->l = p;
//       p->r = c;
//     }
//     if (pp && pp->l == p) pp->l = n;
//     if (pp && pp->r == p) pp->r = n;
//     n->p = pp;
//     p->p = n;
//     if (c) c->p = p;
//   }

//   void prop_from_root(np c) {
//     if (!c->p) {
//       c->nodeProp();
//       return;
//     }
//     prop_from_root(c->p);
//     c->nodeProp();
//   }

//   void splay(Node *me, bool prop_from_root_done) {
//     // これを呼ぶ時点で、me の祖先（me を除く）は既に nodeProp 済であることを仮定
//     // 特に、splay 終了時点で me は upd / nodeProp 済である
//     if (!prop_from_root_done) prop_from_root(me);
//     me->nodeProp();
//     while (me->p) {
//       np p = me->p;
//       np pp = p->p;
//       if (!pp) {
//         rotate(me);
//         p->nodeUpdate();
//         break;
//       }
//       bool same = (p->l == me && pp->l == p) || (p->r == me && pp->r == p);
//       if (same) rotate(p), rotate(me);
//       if (!same) rotate(me), rotate(me);
//       pp->nodeUpdate(), p->nodeUpdate();
//     }
//     // me の nodeUpdate は最後だけでよい
//     me->nodeUpdate();
//   }

//   void splay_kth(np &root, u32 k) {
//     assert(0 <= k && k < (root->size));
//     while (1) {
//       root->nodeProp();
//       u32 sl = (root->l ? root->l->size : 0);
//       if (k == sl) break;
//       if (k < sl)
//         root = root->l;
//       else {
//         k -= sl + 1;
//         root = root->r;
//       }
//     }
//     splay(root, true);
//   }

//   // check(x), 左側のノード全体が check を満たすように切る
//   template <typename F>
//   pair<np, np> split_max_right(np root, F check) {
//     if (!root) return {nullptr, nullptr};
//     assert(!root->p);
//     np c = find_max_right(root, check);
//     if (!c) {
//       splay(root, true);
//       return {nullptr, root};
//     }
//     splay(c, true);
//     np right = c->r;
//     if (!right) return {c, nullptr};
//     right->p = nullptr;
//     c->r = nullptr;
//     c->nodeUpdate();
//     return {c, right};
//   }

//   // check(x, cnt), 左側のノード全体が check を満たすように切る
//   template <typename F>
//   pair<np, np> split_max_right_cnt(np root, F check) {
//     if (!root) return {nullptr, nullptr};
//     assert(!root->p);
//     np c = find_max_right_cnt(root, check);
//     if (!c) {
//       splay(root, true);
//       return {nullptr, root};
//     }
//     splay(c, true);
//     np right = c->r;
//     if (!right) return {c, nullptr};
//     right->p = nullptr;
//     c->r = nullptr;
//     c->nodeUpdate();
//     return {c, right};
//   }

//   // 左側のノード全体の prod が check を満たすように切る
//   template <typename F>
//   pair<np, np> split_max_right_prod(np root, F check) {
//     if (!root) return {nullptr, nullptr};
//     assert(!root->p);
//     np c = find_max_right_prod(root, check);
//     if (!c) {
//       splay(root, true);
//       return {nullptr, root};
//     }
//     splay(c, true);
//     np right = c->r;
//     if (!right) return {c, nullptr};
//     right->p = nullptr;
//     c->r = nullptr;
//     c->nodeUpdate();
//     return {c, right};
//   }

//   template <typename F>
//   np find_max_right(np root, const F &check) {
//     // 最後に見つけた ok の点、最後に探索した点
//     np last_ok = nullptr, last = nullptr;
//     while (root) {
//       last = root;
//       root->nodeProp();
//       if (check(root->x)) {
//         last_ok = root;
//         root = root->r;
//       } else {
//         root = root->l;
//       }
//     }
//     splay(last, true);
//     return last_ok;
//   }

//   template <typename F>
//   np find_max_right_cnt(np root, const F &check) {
//     // 最後に見つけた ok の点、最後に探索した点
//     np last_ok = nullptr, last = nullptr;
//     ll n = 0;
//     while (root) {
//       last = root;
//       root->nodeProp();
//       ll ns = (root->l ? root->l->size : 0);
//       if (check(root->x, n + ns + 1)) {
//         last_ok = root;
//         n += ns + 1;
//         root = root->r;
//       } else {
//         root = root->l;
//       }
//     }
//     splay(last, true);
//     return last_ok;
//   }

//   template <typename F>
//   np find_max_right_prod(np root, const F &check) {
//     using Mono = typename Node::Monoid_X;
//     X prod = Mono::unit();
//     // 最後に見つけた ok の点、最後に探索した点
//     np last_ok = nullptr, last = nullptr;
//     while (root) {
//       last = root;
//       root->nodeProp();
//       X lprod = prod;
//       if (root->l) lprod = Mono::op(lprod, root->l->prod);
//       lprod = Mono::op(lprod, root->x);
//       if (check(lprod)) {
//         prod = lprod;
//         last_ok = root;
//         root = root->r;
//       } else {
//         root = root->l;
//       }
//     }
//     splay(last, true);
//     return last_ok;
//   }
// };

package main

type E = float64

func e() E {
	return 0
}

func NewSpalyTreeBasic() *SplayTreeBasic {
	return &SplayTreeBasic{}
}

type SplayNode struct {
	rev     bool
	size    int32
	x       E
	p, l, r *SplayNode
}

type SplayTreeBasic struct{}

func (st *SplayTreeBasic) NewRoot() *SplayNode {
	return nil
}

func (st *SplayTreeBasic) Build(n int32, f func(i int32) E) *SplayNode {
	var dfs func(l, r int32) *SplayNode
	dfs = func(l, r int32) *SplayNode {
		if l == r {
			return nil
		}
		if r == l+1 {
			return st.newNode(f(l))
		}
		m := (l + r) >> 1
		lRoot, rRoot := dfs(l, m), dfs(m+1, r)
		root := st.newNode(f(m))
		root.l, root.r = lRoot, rRoot
		if lRoot != nil {
			lRoot.p = root
		}
		if rRoot != nil {
			rRoot.p = root
		}
		st.nodeUpdate(root)
		return root
	}
	return dfs(0, n)
}

func (st *SplayTreeBasic) Size(n *SplayNode) int32 {
	if n == nil {
		return 0
	}
	return n.size
}

func (st *SplayTreeBasic) Merge(l, r *SplayNode) *SplayNode {
	if l == nil {
		return r
	}
	if r == nil {
		return l
	}
	st.splayKth(r, 0)
	r.l = l
	l.p = r
	st.nodeUpdate(r)
	return r
}

func (st *SplayTreeBasic) Merge3(a, b, c *SplayNode) *SplayNode {
	return st.Merge(st.Merge(a, b), c)
}

func (st *SplayTreeBasic) Merge4(a, b, c, d *SplayNode) *SplayNode {
	return st.Merge(st.Merge(st.Merge(a, b), c), d)
}

func (st *SplayTreeBasic) Split(root *SplayNode, k int32) (*SplayNode, *SplayNode) {
	if k == 0 {
		return nil, root
	}
	if k == root.size {
		return root, nil
	}
	st.splayKth(root, k-1)
	right := root.r
	root.r = nil
	root.p = nil
	st.nodeUpdate(root)
	return root, right
}

func (st *SplayTreeBasic) Split3(root *SplayNode, l, r int32) (*SplayNode, *SplayNode, *SplayNode) {
	var nm, nr *SplayNode
	root, nr = st.Split(root, r)
	root, nm = st.Split(root, l)
	return root, nm, nr
}

func (st *SplayTreeBasic) Split4(root *SplayNode, i, j, k int32) (*SplayNode, *SplayNode, *SplayNode, *SplayNode) {
	var d *SplayNode
	root, d = st.Split(root, k)
	a, b, c := st.Split3(root, i, j)
	return a, b, c, d
}

func (st *SplayTreeBasic) gotoBetween(root *SplayNode, l, r int32) {
	if l == 0 && r == root.size {
		return
	}
	if l == 0 {
		st.splayKth(root, r)
		root = root.l
		return
	}
	if r == root.size {
		st.splayKth(root, l-1)
		root = root.r
		return
	}
	st.splayKth(root, r)
	rp := root
	root = rp.l
	root.p = nil
	st.splayKth(root, l-1)
	root.p = rp
	rp.l = root
	st.nodeUpdate(rp)
	root = root.r
}

func (st *SplayTreeBasic) EnumerateAll(root *SplayNode, f func(E)) {
	var dfs func(*SplayNode)
	dfs = func(root *SplayNode) {
		if root == nil {
			return
		}
		st.nodeProp(root)
		dfs(root.l)
		f(st.nodeGet(root))
		dfs(root.r)
	}
	dfs(root)
}

func (st *SplayTreeBasic) Get(root *SplayNode, k int32) E {
	st.splayKth(root, k)
	return st.nodeGet(root)
}

func (st *SplayTreeBasic) Set(root *SplayNode, k int32, x E) {
	st.splayKth(root, k)
	st.nodeSet(root, x)
}

func (st *SplayTreeBasic) Update(root *SplayNode, k int32, x E) {
	st.splayKth(root, k)
	st.nodeSet(root, x)
}

// func (st *SplayTreeBasic) Query(root *SplayNode, l, r int32) E {
// 	if l == r {
// 		return e()
// 	}
// 	if l < 0 {
// 		l = 0
// 	}
// 	if r > root.size {
// 		r = root.size
// 	}
// 	if l >= r {
// 		return e()
// 	}
// 	st.gotoBetween(root, l, r)
// 	res := root.sum
// 	st.splay(root, true)
// 	return res
// }

// func (st *SplayTreeBasic) QueryAll(root *SplayNode) E {
// 	if root == nil {
// 		return e()
// 	}
// 	return st.sum
// }

// 私有方法需要重写
func (st *SplayTreeBasic) newNode(x E) *SplayNode {
	return &SplayNode{x: x, size: 1}
}

func (st *SplayTreeBasic) nodeUpdate(n *SplayNode) {
	n.size = 1
	if n.l != nil {
		n.size += n.l.size
	}
	if n.r != nil {
		n.size += n.r.size
	}
}

func (st *SplayTreeBasic) nodeProp(n *SplayNode) {
	if n.rev {
		if left := n.l; left != nil {
			left.rev = !left.rev
			left.l, left.r = left.r, left.l
		}
		if right := n.r; right != nil {
			right.rev = !right.rev
			right.l, right.r = right.r, right.l
		}
		n.rev = false
	}
}

func (st *SplayTreeBasic) nodeGet(n *SplayNode) E {
	return n.x
}

func (st *SplayTreeBasic) nodeSet(n *SplayNode, x E) {
	n.x = x
	st.nodeUpdate(n)
}

func (st *SplayTreeBasic) nodeReverse(n *SplayNode) {
	n.l, n.r = n.r, n.l
	n.rev = !n.rev
}
