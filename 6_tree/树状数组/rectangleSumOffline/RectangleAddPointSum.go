// #include "ds/fenwicktree/fenwicktree.hpp"

// template <typename AbelGroup, typename XY, bool SMALL_X = false>
// struct Rectangle_Add_Point_Sum {
//   using G = typename AbelGroup::value_type;
//   vector<tuple<XY, XY, XY, G>> rect;
//   vector<tuple<int, XY, XY>> point;

//   Rectangle_Add_Point_Sum() {}

//   void add_query(XY x1, XY x2, XY y1, XY y2, G g) {
//     rect.eb(y1, x1, x2, g), rect.eb(y2, x2, x1, g);
//   }
//   void sum_query(XY x, XY y) { point.eb(len(point), x, y); }

//   vector<G> calc() {
//     int N = rect.size(), Q = point.size();
//     if (N == 0 || Q == 0) return vector<G>(Q, AbelGroup::unit());
//     // X 方向の座圧
//     int NX = 0;
//     if (!SMALL_X) {
//       sort(all(point),
//            [&](auto &x, auto &y) -> bool { return get<1>(x) < get<1>(y); });
//       vc<XY> keyX;
//       keyX.reserve(Q);
//       for (auto &&[i, a, b]: point) {
//         if (len(keyX) == 0 || keyX.back() != a) { keyX.eb(a); }
//         a = len(keyX) - 1;
//       }
//       for (auto &&[y, x1, x2, g]: rect) x1 = LB(keyX, x1), x2 = LB(keyX, x2);
//       NX = len(keyX);
//     }
//     if (SMALL_X) {
//       XY mx = infty<XY>;
//       for (auto &&[i, x, y]: point) chmin(mx, x);
//       for (auto &&[i, x, y]: point) x -= mx, chmax(NX, x + 1);
//       for (auto &&[y, x1, x2, g]: rect) {
//         x1 -= mx, x2 -= mx;
//         x1 = max(0, min<int>(x1, NX)), x2 = max(0, min<int>(x2, NX));
//       }
//     }

//     sort(all(point),
//          [&](auto &x, auto &y) -> bool { return get<2>(x) < get<2>(y); });
//     sort(all(rect),
//          [&](auto &x, auto &y) -> bool { return get<0>(x) < get<0>(y); });
//     FenwickTree<AbelGroup> bit(NX);
//     vc<G> res(Q, AbelGroup::unit());
//     int j = 0;
//     FOR(i, Q) {
//       auto [q, x, y] = point[i];
//       while (j < N && get<0>(rect[j]) <= y) {
//         auto [yy, x1, x2, g] = rect[j++];
//         bit.add(x1, g), bit.add(x2, AbelGroup::inverse(g));
//       }
//       res[q] = bit.sum(x + 1);
//     }
//     return res;
//   }
// };

// 矩形区间修改单点查询(离线)

package main

import "sort"

func main() {
	e := func() int32 { return 0 }
	op := func(a, b int32) int32 { return a + b }
	inv := func(a int32) int32 { return -a }
	ps := NewPointAddRectangleSumOffline(e, op, inv, false)
	ps.AddRectangle(0, 2, 0, 2, 1)
	ps.AddRectangle(1, 3, 1, 3, 2)
	ps.AddQuery(1, 1)
	ps.AddQuery(2, 2)
	ps.AddQuery(3, 3)
	res := ps.Calc()
	for _, v := range res {
		println(v)
	}
}

const INF32 int32 = 1e9 + 10

type rect[E any] struct {
	y, x1, x2 int32
	w         E
}

type point struct {
	id   int32
	x, y int32
}

type RectangleAddPointSumOffline[E any] struct {
	smallX bool
	rects  []rect[E]
	points []point
	e      func() E
	op     func(e1, e2 E) E
	inv    func(e E) E
}

func NewPointAddRectangleSumOffline[E any](
	e func() E, op func(e1, e2 E) E, inv func(e E) E,
	smallX bool,
) *RectangleAddPointSumOffline[E] {
	return &RectangleAddPointSumOffline[E]{smallX: smallX, e: e, op: op, inv: inv}
}

func (ps *RectangleAddPointSumOffline[E]) AddRectangle(xl, xr, yl, yr int32, w E) {
	ps.rects = append(ps.rects, rect[E]{y: yl, x1: xl, x2: xr, w: w})
	ps.rects = append(ps.rects, rect[E]{y: yr, x1: xr, x2: xl, w: w})
}

func (ps *RectangleAddPointSumOffline[E]) AddQuery(x, y int32) {
	ps.points = append(ps.points, point{id: int32(len(ps.points)), x: x, y: y})
}

func (ps *RectangleAddPointSumOffline[E]) Calc() []E {
	n, q := int32(len(ps.rects)), int32(len(ps.points))
	if n == 0 || q == 0 {
		res := make([]E, q)
		for i := int32(0); i < q; i++ {
			res[i] = ps.e()
		}
		return res
	}

	rects, points := ps.rects, ps.points
	// compress x
	nx := int32(0)
	if !ps.smallX {
		sort.Slice(points, func(i, j int) bool { return points[i].x < points[j].x })
		keyX := make([]int32, 0, q)
		for i := int32(0); i < q; i++ {
			x := points[i].x
			if len(keyX) == 0 || keyX[len(keyX)-1] != x {
				keyX = append(keyX, x)
			}
			points[i].x = int32(len(keyX) - 1)
		}
		for i := int32(0); i < q; i++ {
			rects[i].x1 = lowerBound32(keyX, rects[i].x1)
			rects[i].x2 = lowerBound32(keyX, rects[i].x2)
		}
		nx = int32(len(keyX))
	} else {
		mx := INF32
		for i := int32(0); i < n; i++ {
			if tmp := points[i].x; tmp < mx {
				mx = tmp
			}
		}
		for i := int32(0); i < n; i++ {
			points[i].x -= mx
			if tmp := points[i].x + 1; tmp > nx {
				nx = tmp
			}
		}
		for i := int32(0); i < q; i++ {
			rects[i].x1 = clamp32(rects[i].x1-mx, 0, nx)
			rects[i].x2 = clamp32(rects[i].x2-mx, 0, nx)
		}
	}

	sort.Slice(points, func(i, j int) bool { return points[i].y < points[j].y })
	sort.Slice(rects, func(i, j int) bool { return rects[i].y < rects[j].y })

	bit := newFenwickTree(ps.e, ps.op, ps.inv)
	bit.Build(nx, func(i int32) E { return ps.e() })
	res := make([]E, q)
	for i := range res {
		res[i] = ps.e()
	}
	j := int32(0)
	for _, point := range points {
		q, x, y := point.id, point.x, point.y
		for j < n && rects[j].y <= y {
			rect := rects[j]
			bit.Update(rect.x1, rect.w)
			bit.Update(rect.x2, ps.inv(rect.w))
			j++
		}
		res[q] = bit.QueryPrefix(x + 1)
	}
	return res
}

type fenwickTree[E any] struct {
	n     int32
	total E
	data  []E
	e     func() E
	op    func(e1, e2 E) E
	inv   func(e E) E
}

func newFenwickTree[E any](e func() E, op func(e1, e2 E) E, inv func(e E) E) *fenwickTree[E] {
	return &fenwickTree[E]{e: e, op: op, inv: inv}
}

func (fw *fenwickTree[E]) Build(n int32, f func(i int32) E) {
	data := make([]E, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	for i := int32(1); i <= n; i++ {
		if j := i + (i & -i); j <= n {
			data[j-1] = fw.op(data[i-1], data[j-1])
		}
	}
	fw.n = n
	fw.data = data
	fw.total = fw.QueryPrefix(n)
}

func (fw *fenwickTree[E]) QueryAll() E { return fw.total }

// [0, end)
func (fw *fenwickTree[E]) QueryPrefix(end int32) E {
	if end > fw.n {
		end = fw.n
	}
	res := fw.e()
	for ; end > 0; end &= end - 1 {
		res = fw.op(res, fw.data[end-1])
	}
	return res
}

// [start, end)
func (fw *fenwickTree[E]) QueryRange(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > fw.n {
		end = fw.n
	}
	if start > end {
		return fw.e()
	}
	if start == 0 {
		return fw.QueryPrefix(end)
	}
	pos, neg := fw.e(), fw.e()
	for end > start {
		pos = fw.op(pos, fw.data[end-1])
		end &= end - 1
	}
	for start > end {
		neg = fw.op(neg, fw.data[start-1])
		start &= start - 1
	}
	return fw.op(pos, fw.inv(neg))
}

// 要求op满足交换律(commute).
func (fw *fenwickTree[E]) Update(i int32, x E) {
	fw.total = fw.op(fw.total, x)
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = fw.op(fw.data[i-1], x)
	}
}

func (fw *fenwickTree[E]) GetAll() []E {
	res := make([]E, fw.n)
	for i := int32(0); i < fw.n; i++ {
		res[i] = fw.QueryRange(i, i+1)
	}
	return res
}

func lowerBound32(nums []int32, x int32) int32 {
	left, right := int32(0), int32(len(nums)-1)
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func clamp32(x, l, r int32) int32 {
	if x < l {
		return l
	}
	if x > r {
		return r
	}
	return x
}
