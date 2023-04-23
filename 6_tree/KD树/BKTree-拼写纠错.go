// https://zhuanlan.zhihu.com/p/360925212
// https://www.twblogs.net/a/5c22542bbd9eee16b4a778a2
// https://yfsyfs.github.io/2019/06/25/%E7%BA%A0%E9%94%99%E5%88%A9%E5%99%A8-BK%E6%A0%91/
// BKTree实现拼写纠错

package main

func main() {

}

// #define fst first
// #define snd second
// typedef pair<int,int> PII;
// int dist(PII a, PII b) { return max(abs(a.fst-b.fst),abs(a.snd-b.snd)); }
// void process(PII a) { printf("%d %d\n", a.fst,a.snd); }

// template <class T>
// struct bk_tree {
//   typedef int dist_type;
//   struct node {
//     T p;
//     unordered_map<dist_type, node*> ch;
//   } *root;
//   bk_tree() : root(0) { }

//   node *insert(node *n, T p) {
//     if (!n) { n = new node(); n->p = p; return n; }
//     dist_type d = dist(n->p, p);
//     n->ch[d] = insert(n->ch[d], p);
//     return n;
//   }
//   void insert(T p) { root = insert(root, p); }
//   void traverse(node *n, T p, dist_type dmax) {
//     if (!n) return;
//     dist_type d = dist(n->p, p);
//     if (d < dmax) {
//       process(n->p); // write your process
//     }
//     for (auto i: n->ch)
//       if (-dmax <= i.fst - d && i.fst - d <= dmax)
//         traverse(i.snd, p, dmax);
//   }
//   void traverse(T p, dist_type dmax) { traverse(root, p, dmax); }
// };

type BKTree struct {
}

func NewBKTree() *BKTree {
	return &BKTree{}
}

type BNode struct {
}
