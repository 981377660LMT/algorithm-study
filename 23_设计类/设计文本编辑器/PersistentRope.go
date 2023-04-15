// // Persistent Rope
// //
// // Description:
// //   Rope is a binary tree data structure to maintains a sequence.
// //

// #include <iostream>
// #include <vector>
// #include <cstdio>
// #include <sstream>
// #include <cstdlib>
// #include <map>
// #include <cmath>
// #include <cstring>
// #include <functional>
// #include <algorithm>
// #include <unordered_map>
// #include <unordered_set>

// using namespace std;

// #define fst first
// #define snd second
// #define all(c) ((c).begin()), ((c).end())

// struct rope {
//   struct node {
//     char v;
//     node *l, *r;
//     int s;
//     node(char v, node *l, node *r) : v(v), l(l), r(r) {
//       s = 1 + (l ? l->s : 0) + (r ? r->s : 0);
//     }
//   } *root;
//   node *join(node *a, node *b) {
//     auto R = [](int a, int b) { return rand() % (a + b) < a; };
//     if (!a || !b) return a ? a : b;
//     if (R(a->s, b->s)) return new node(a->v, a->l, join(a->r, b));
//     else               return new node(b->v, join(a, b->l), b->r);
//   }
//   pair<node*,node*> split(node *a, int s) {
//     if (!a || s <= 0) return {0, a};
//     if (a->s <= s)    return {a, 0};
//     if (a->l && s <= a->l->s) {
//       auto p = split(a->l, s);
//       return {p.fst, new node(a->v, p.snd, a->r)};
//     } else {
//       auto p = split(a->r, s - (a->l ? a->l->s : 0) - 1);
//       return {new node(a->v, a->l, p.fst), p.snd};
//     }
//   }
//   pair<node*, node*> cut(node *a, int l, int r) { // (sub, rest)
//     if (l >= r) return {0, a};
//     auto p = split(a, l), q = split(p.snd, r - l);
//     return {q.fst, join(p.fst, q.snd)};
//   }
//   rope(const char s[]) {
//     function<node*(int,int)> build = [&](int l, int r) {
//       if (l >= r) return (node*)0;
//       int m = (l + r) / 2;
//       return new node(s[m], build(l, m), build(m+1, r));
//     };
//     root = build(0, strlen(s));
//   }
//   rope(rope::node *r) : root(r) { }
//   int size() const { return root ? root->s : 0; }
//   rope insert(int k, const char s[]) {
//     auto p = split(root, k);
//     return {join(p.fst, join(rope(s).root, p.snd))};
//   }
//   rope substr(int l, int r) { return {cut(root, l, r).fst}; }
//   rope erase(int l, int r) { return {cut(root, l, r).snd}; }
//   char _at(int k) const {
//     function<char(node*)> rec = [&](node *a) {
//       int s = a->l ? a->l->s : 0;
//       if (k == s) return a->v;
//       if (k < s) return rec(a->l);
//       k -= s+1; return rec(a->r);
//     };
//     return rec(root);
//   }
//   string str() const {
//     stringstream ss;
//     function<void(node*)> rec = [&](node *a) {
//       if (!a) return;
//       rec(a->l); ss << a->v; rec(a->r);
//     }; rec(root);
//     return ss.str();
//   }
// };

// int main() {
//   rope a("abcde"), b("ABCDE");
//   for (int i = 0; i < 5; ++i) {
//     cout << a._at(i) << " ";
//   }
// }
package main

import (
	"fmt"
	"math/rand"
)

func main() {
	root := NewPersistentRope([]string{"a", "b", "c", "d", "e"})
	root = root.Insert(0, []string{"A", "B", "C", "D", "E"})
	fmt.Println(root)
	fmt.Println(root.Substring(0, 3))
}

type R = string
type PersistentRope struct {
	root *_RNode
}

func NewPersistentRope(bytes []R) *PersistentRope {
	res := &PersistentRope{}
	res.root = res._build(0, len(bytes), bytes)
	return res
}

func (rp *PersistentRope) Insert(i int, bytes []R) *PersistentRope {
	p1, p2 := rp._split(rp.root, i)
	newRoot := rp._join(p1, rp._join(NewPersistentRope(bytes).root, p2))
	return &PersistentRope{newRoot}
}

func (rp *PersistentRope) Substring(start, end int) *PersistentRope {
	first, _ := rp._cut(rp.root, start, end)
	return &PersistentRope{first}
}

func (rp *PersistentRope) Erase(start, end int) *PersistentRope {
	_, second := rp._cut(rp.root, start, end)
	return &PersistentRope{second}
}

func (rp *PersistentRope) At(i int) R {
	var rec func(node *_RNode) R
	rec = func(node *_RNode) R {
		s := 0
		if node.l != nil {
			s = node.l.s
		}
		if i == s {
			return node.v
		}
		if i < s {
			return rec(node.l)
		}
		i -= s + 1
		return rec(node.r)
	}
	return rec(rp.root)
}

func (rp *PersistentRope) Size() int {
	if rp.root == nil {
		return 0
	}
	return rp.root.s
}

func (rp *PersistentRope) String() string {
	sb := []string{}
	var dfs func(node *_RNode)
	dfs = func(node *_RNode) {
		if node == nil {
			return
		}
		dfs(node.l)
		sb = append(sb, fmt.Sprintf("%v", node.v))
		dfs(node.r)
	}
	dfs(rp.root)
	return fmt.Sprintf("%v", sb)
}

func (rp *PersistentRope) _build(l, r int, bytes []R) *_RNode {
	if l >= r {
		return nil
	}
	m := (l + r) >> 1
	return _NewNode(bytes[m], rp._build(l, m, bytes), rp._build(m+1, r, bytes))
}

func (rp *PersistentRope) _join(a, b *_RNode) *_RNode {
	if a == nil || b == nil {
		if a != nil {
			return a
		}
		return b
	}
	if r, s1, s2 := rand.Int(), a.s, b.s; r%(s1+s2) < s1 {
		return _NewNode(a.v, a.l, rp._join(a.r, b))
	}
	return _NewNode(b.v, rp._join(a, b.l), b.r)
}

func (rp *PersistentRope) _split(a *_RNode, s int) (*_RNode, *_RNode) {
	if a == nil || s <= 0 {
		return nil, a
	}
	if a.s <= s {
		return a, nil
	}
	if a.l != nil && s <= a.l.s {
		p1, p2 := rp._split(a.l, s)
		return p1, _NewNode(a.v, p2, a.r)
	}
	tmp := 0
	if a.l != nil {
		tmp = a.l.s
	}
	l, r := rp._split(a.r, s-tmp-1)
	return _NewNode(a.v, a.l, l), r
}

// (sub,rest)
func (rp *PersistentRope) _cut(a *_RNode, l, r int) (*_RNode, *_RNode) {
	if l >= r {
		return nil, a
	}
	p1, p2 := rp._split(a, l)
	q1, q2 := rp._split(p2, r-l)
	return q1, rp._join(p1, q2)
}

type _RNode struct {
	v    R
	l, r *_RNode
	s    int
}

func _NewNode(v R, l, r *_RNode) *_RNode {
	res := &_RNode{v: v, l: l, r: r, s: 1}
	if l != nil {
		res.s += l.s
	}
	if r != nil {
		res.s += r.s
	}
	return res
}
