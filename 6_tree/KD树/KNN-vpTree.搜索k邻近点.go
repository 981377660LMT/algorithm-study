//
// Vantage Point Tree (vp tree)
//
// Description
//   Vantage point tree is a metric tree.
//   Each tree node has a point, radius, and two childs.
//   The points of left descendants are contained in the ball B(p,r)
//   and the points of right descendants are exluded from the ball.
//
//   We can find k-nearest neighbors of a given point p efficiently
//   by pruning search.
//
//   The data structure is independently proposed by J. Uhlmann and
//   P. N. Yianilos.
//
// Complexity:
//   Construction: O(n log n).
//   Search: O(log n)
//
//   In my implementation, its construction is few times slower than kd tree
//   and its search is bit faster than kd tree.
//
// References
//   J. Uhlmann (1991):
//   Satisfying General Proximity/Similarity Queries with Metric Trees.
//   Information Processing Letters, vol. 40, no. 4, pp. 175--179.
//
//   Peter N. Yianilos (1993):
//   Data structures and algorithms for nearest neighbor search in general metric spaces.
//   in Proceedings of the 4th Annual ACM-SIAM Symposium on Discrete algorithms,
//   Society for Industrial and Applied Mathematics Philadelphia, PA, USA. pp. 311--321.

// 近邻搜索之制高点树（VP-Tree）
// https://www.cxyzjd.com/article/y459541195/102846739
// 以图搜图

package main

// typedef complex<double> point;
// namespace std {
// bool operator < (point p, point q) {
//   if (real(p) != real(q)) return real(p) < real(q);
//   return imag(p) < imag(q);
// }
// };
// struct vantage_point_tree {
//   struct node {
//     point p;
//     double th;
//     node *l, *r;
//   } *root;
//   vector<pair<double, point>> aux;
//   vantage_point_tree(vector<point> ps) {
//     for (int i = 0; i < ps.size(); ++i)
//       aux.push_back({0, ps[i]});
//     root = build(0, ps.size());
//   }
//   node *build(int l, int r) {
//     if (l == r) return 0;
//     swap(aux[l], aux[l + rand() % (r - l)]);
//     point p = aux[l++].snd;
//     if (l == r) return new node({p});
//     for (int i = l; i < r; ++i)
//       aux[i].fst = norm(p - aux[i].snd);
//     int m = (l + r) / 2;
//     nth_element(aux.begin()+l, aux.begin()+m, aux.begin()+r);
//     return new node({p, sqrt(aux[m].fst), build(l, m), build(m, r)});
//   }
//   priority_queue<pair<double, node*>> que;
//   void k_nn(node *t, point p, int k) {
//     if (!t) return;
//     double d = abs(p - t->p);
//     if (que.size() < k) que.push({d, t});
//     else if (que.top().fst > d) {
//       que.pop();
//       que.push({d, t});
//     }
//     if (!t->l && !t->r) return;
//     if (d < t->th) {
//       k_nn(t->l, p, k);
//       if (t->th - d <= que.top().fst) k_nn(t->r, p, k);
//     } else {
//       k_nn(t->r, p, k);
//       if (d - t->th <= que.top().fst) k_nn(t->l, p, k);
//     }
//   }
//   vector<point> k_nn(point p, int k) {
//     k_nn(root, p, k);
//     vector<point> ans;
//     for (; !que.empty(); que.pop())
//       ans.push_back(que.top().snd->p);
//     reverse(all(ans));
//     return ans;
//   }
// };

func main() {

}

type V int
type Point2D struct{ x, y V }
type VantagePointTree struct {
	root     *VNode
	aux      []VNode
	calDist2 func(p1, p2 Point2D) V // 计算两点距离的平方
}

func NewVantagePointTree(points []Point2D, calDist2 func(p1, p2 Point2D) V) *VantagePointTree {}

func (vp *VantagePointTree) Build(points []Point2D) {}

func (vp *VantagePointTree) KNN(p Point2D, k int) []Point2D {
	vp._knn(vp.root, p, k)
	var res []Point2D
}

func (vp *VantagePointTree) _knn(t *VNode, p Point2D, k int) {}

type VNode struct {
	point Point2D
	th    V
	l, r  *VNode
}
